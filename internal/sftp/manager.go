package sftp

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/ilaziness/vexo/internal/ssh"
	"github.com/ilaziness/vexo/internal/transfer"
	"github.com/pkg/sftp"
	"go.uber.org/zap"
)

const (
	ErrTrackerRequired = "tracker is required"
)

// FileInfo represents file information for SFTP operations
type FileInfo struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	Mode    string `json:"mode"`
	ModTime string `json:"modTime"`
	IsDir   bool   `json:"isDir"`
}

// Convert os.FileInfo to our custom FileInfo
func convertFileInfo(osInfo os.FileInfo) FileInfo {
	return FileInfo{
		Name:    osInfo.Name(),
		Size:    osInfo.Size(),
		Mode:    osInfo.Mode().String(),
		ModTime: osInfo.ModTime().Format(time.RFC3339),
		IsDir:   osInfo.IsDir(),
	}
}

// Convert []os.FileInfo to []FileInfo
func convertFileInfos(osInfos []os.FileInfo) []FileInfo {
	result := make([]FileInfo, len(osInfos))
	for i, osInfo := range osInfos {
		result[i] = convertFileInfo(osInfo)
	}
	return result
}

type Manager struct {
	logger    *zap.Logger
	ssh       *ssh.Manager
	transfers *transfer.Registry
	clients   sync.Map
}

func NewManager(logger *zap.Logger, sshMgr *ssh.Manager, transfers *transfer.Registry) *Manager {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &Manager{logger: logger, ssh: sshMgr, transfers: transfers}
}

// joinRemotePath joins remote path elements with forward slash
// SFTP servers (typically Unix/Linux) always use forward slash as path separator
func joinRemotePath(elem ...string) string {
	return path.Join(elem...)
}

// getSftpClient 通过 sessionID 获取对应的 SFTP 客户端
func (sft *Manager) getSftpClient(sessionID string) (*sftp.Client, error) {
	connVal, ok := sft.clients.Load(sessionID)
	if !ok {
		return nil, fmt.Errorf("SSH session with ID %s not found", sessionID)
	}
	return connVal.(*sftp.Client), nil
}

func (sft *Manager) Connect(sessionID string) error {
	if _, err := sft.getSftpClient(sessionID); err == nil {
		return nil
	}
	sshClient, err := sft.ssh.GetClient(sessionID)
	if err != nil {
		return err
	}
	ftpClient, err := sftp.NewClient(
		sshClient,
		sftp.MaxPacket(32768),
		sftp.UseConcurrentWrites(true),
		sftp.UseConcurrentReads(true),
	)
	if err != nil {
		return err
	}
	sft.clients.Store(sessionID, ftpClient)
	return nil
}

func JoinRemotePath(elem ...string) string {
	return joinRemotePath(elem...)
}

// ListFiles lists files and directories in the specified path
func (sft *Manager) ListFiles(sessionID string, path string, showHidden bool) ([]FileInfo, error) {
	sft.logger.Debug("ListFiles", zap.String("sessionID", sessionID), zap.String("path", path), zap.Bool("showHidden", showHidden))
	ftpClient, err := sft.getSftpClient(sessionID)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()
	entries, err := ftpClient.ReadDirContext(ctx, path)
	if err != nil {
		return nil, err
	}

	// Filter out hidden files if showHidden is false
	if !showHidden {
		var filteredEntries []os.FileInfo

		for _, entry := range entries {
			// Check if the file/directory name starts with a dot (hidden file in Unix systems)
			if len(entry.Name()) > 0 && entry.Name()[0] != '.' {
				filteredEntries = append(filteredEntries, entry)
			}
		}

		entries = filteredEntries
	}

	return convertFileInfos(entries), nil
}

// GetFileInfo returns information about a file or directory
func (sft *Manager) GetFileInfo(sessionID string, path string) (FileInfo, error) {
	sft.logger.Debug("GetFileInfo", zap.String("sessionID", sessionID), zap.String("path", path))
	ftpClient, err := sft.getSftpClient(sessionID)
	if err != nil {
		return FileInfo{}, err
	}

	info, err := ftpClient.Stat(path)
	if err != nil {
		return FileInfo{}, err
	}

	return convertFileInfo(info), nil
}

// uploadFile uploads a local file to the specified remote path.
// localPathFile is local file path with file full name.
// remoteDir is remote directory without file name.
func (sft *Manager) uploadFile(sessionID string, localPathFile, remoteDir string, tracker *transfer.Tracker) error {
	sft.logger.Debug("uploadFile", zap.String("sessionID", sessionID), zap.String("localPathFile", localPathFile), zap.String("remoteDir", remoteDir))
	ftpClient, err := sft.getSftpClient(sessionID)
	if err != nil {
		return err
	}

	// Open the local file for reading
	localFile, err := os.Open(localPathFile)
	if err != nil {
		return err
	}
	defer localFile.Close()
	_, err = ftpClient.Stat(remoteDir)
	if err != nil && os.IsNotExist(err) {
		if err = ftpClient.MkdirAll(remoteDir); err != nil {
			return err
		}
	}

	// Create the remote file for writing
	remoteFilePath := joinRemotePath(remoteDir, filepath.Base(localPathFile))
	sft.logger.Debug("uploadFile Create File", zap.String("file", remoteFilePath))
	remoteFile, err := ftpClient.Create(remoteFilePath)
	if err != nil {
		return err
	}
	defer remoteFile.Close()

	if tracker == nil {
		return fmt.Errorf(ErrTrackerRequired)
	}

	progressReader := &transfer.Reader{Reader: localFile, Tracker: tracker}

	// Copy the local file content to the remote file using io.Copy to leverage sftp.File.ReadFrom concurrent writes
	_, err = io.Copy(remoteFile, progressReader)
	return err
}

func (sft *Manager) UploadFile(sessionID, localPath, remotePath string) error {
	return sft.uploadFileWithTracker(sessionID, localPath, remotePath)
}

func (sft *Manager) uploadFileWithTracker(sessionID, localPath, remotePath string) error {
	info, err := os.Stat(localPath)
	if err != nil {
		return err
	}
	tracker := sft.transfers.New(sessionID, transfer.TypeUpload, localPath, joinRemotePath(remotePath, filepath.Base(localPath)), info.Size())
	err = sft.uploadFile(sessionID, localPath, remotePath, tracker)
	tracker.Stop(err)
	return err
}

func (sft *Manager) startUploadFileAsync(sessionID, localPath, remotePath string, size int64) {
	go func() {
		tracker := sft.transfers.New(sessionID, transfer.TypeUpload, localPath, joinRemotePath(remotePath, filepath.Base(localPath)), size)
		err := sft.uploadFile(sessionID, localPath, remotePath, tracker)
		tracker.Stop(err)
	}()
}

// downloadFile downloads a remote file to the specified local path
// localPath local save file path with filename
// remotePathFile remote file path (full path including filename)
func (sft *Manager) downloadFile(sessionID string, localPathFile, remotePathFile string, tracker *transfer.Tracker) error {
	sft.logger.Debug("downloadFile", zap.String("sessionID", sessionID), zap.String("localPath", localPathFile), zap.String("remotePath", remotePathFile))
	ftpClient, err := sft.getSftpClient(sessionID)
	if err != nil {
		return err
	}

	// Ensure local directory exists
	localDir := filepath.Dir(localPathFile)
	if err := os.MkdirAll(localDir, 0755); err != nil {
		return err
	}

	// Open the remote file for reading
	remoteFile, err := ftpClient.Open(remotePathFile)
	if err != nil {
		return err
	}
	defer remoteFile.Close()

	// Create the local file for writing
	localFile, err := os.Create(localPathFile)
	if err != nil {
		return err
	}
	defer localFile.Close()

	if tracker == nil {
		return fmt.Errorf(ErrTrackerRequired)
	}

	progressWriter := &transfer.Writer{Writer: localFile, Tracker: tracker}

	// Copy the remote file content to the local file using io.Copy to leverage sftp.File.WriteTo concurrent reads
	_, err = io.Copy(progressWriter, remoteFile)
	if err != nil {
		return err
	}

	return nil
}

func (sft *Manager) DownloadFile(sessionID, localPathFile, remotePathFile string) error {
	ftpClient, err := sft.getSftpClient(sessionID)
	if err != nil {
		return err
	}
	info, err := ftpClient.Stat(remotePathFile)
	if err != nil {
		return err
	}
	tracker := sft.transfers.New(sessionID, transfer.TypeDownload, localPathFile, remotePathFile, info.Size())
	err = sft.downloadFile(sessionID, localPathFile, remotePathFile, tracker)
	tracker.Stop(err)
	return err
}

func (sft *Manager) UploadDirectory(sessionID, localPath, remotePath string) error {
	return sft.uploadDirectoryWithTracker(sessionID, localPath, remotePath)
}

func (sft *Manager) uploadDirectoryWithTracker(sessionID, localPath, remotePath string) error {
	total, err := sft.calcLocalDirSize(localPath)
	if err != nil {
		return err
	}
	tracker := sft.transfers.New(sessionID, transfer.TypeUpload, localPath, joinRemotePath(remotePath, filepath.Base(localPath)), total)
	err = sft.uploadDirectory(sessionID, localPath, remotePath, tracker)
	tracker.Stop(err)
	return err
}

func (sft *Manager) startUploadDirectoryAsync(sessionID, localPath, remotePath string, total int64) {
	go func() {
		tracker := sft.transfers.New(sessionID, transfer.TypeUpload, localPath, joinRemotePath(remotePath, filepath.Base(localPath)), total)
		err := sft.uploadDirectory(sessionID, localPath, remotePath, tracker)
		tracker.Stop(err)
	}()
}

// UploadPaths uploads local files or directories to the remote path in parallel.
func (sft *Manager) UploadPaths(sessionID, remotePath string, localPaths []string) error {
	if sessionID == "" || remotePath == "" || len(localPaths) == 0 {
		return fmt.Errorf("invalid upload parameters")
	}
	if _, err := sft.getSftpClient(sessionID); err != nil {
		return err
	}
	for _, p := range localPaths {
		info, err := os.Stat(p)
		if err != nil {
			return fmt.Errorf("local path %s: %w", p, err)
		}
		if info.IsDir() {
			total, err := sft.calcLocalDirSize(p)
			if err != nil {
				return fmt.Errorf("local path %s: %w", p, err)
			}
			sft.startUploadDirectoryAsync(sessionID, p, remotePath, total)
		} else {
			sft.startUploadFileAsync(sessionID, p, remotePath, info.Size())
		}
	}
	return nil
}

// uploadDirectory recursively uploads a local directory to the specified remote path
func (sft *Manager) uploadDirectory(sessionID string, localPath, remotePath string, tracker *transfer.Tracker) error {
	sft.logger.Debug("uploadDirectory", zap.String("sessionID", sessionID), zap.String("localPath", localPath), zap.String("remotePath", remotePath))

	// Validate inputs and get clients
	ftpClient, err := sft.getSftpClient(sessionID)
	if err != nil {
		return err
	}

	if tracker == nil {
		return fmt.Errorf(ErrTrackerRequired)
	}

	// Validate local directory
	info, err := os.Stat(localPath)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory")
	}

	// Create remote directory
	remoteDir := joinRemotePath(remotePath, filepath.Base(localPath))
	if err := sft.ensureRemoteDirExists(ftpClient, remoteDir); err != nil {
		return err
	}

	// Process directory contents
	return sft.processDirectoryEntries(sessionID, localPath, remoteDir, tracker)
}

// ensureRemoteDirExists creates the remote directory if it doesn't exist
func (sft *Manager) ensureRemoteDirExists(ftpClient *sftp.Client, remoteDir string) error {
	_, err := ftpClient.Stat(remoteDir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if os.IsNotExist(err) {
		return ftpClient.MkdirAll(remoteDir)
	}
	return nil
}

// processDirectoryEntries processes all entries in a directory
func (sft *Manager) processDirectoryEntries(sessionID, localPath, remoteDir string, tracker *transfer.Tracker) error {
	entries, err := os.ReadDir(localPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		localEntryPath := filepath.Join(localPath, entry.Name())

		if entry.IsDir() {
			// Recursively upload subdirectory
			err = sft.uploadDirectory(sessionID, localEntryPath, remoteDir, tracker)
			if err != nil {
				return err
			}
		} else {
			// Upload file
			err = sft.uploadFile(sessionID, localEntryPath, remoteDir, tracker)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (sft *Manager) DownloadDirectory(sessionID, localPath, remotePath string) error {
	total, err := sft.calcRemoteDirSize(sessionID, remotePath)
	if err != nil {
		return err
	}
	tracker := sft.transfers.New(sessionID, transfer.TypeDownload, localPath, remotePath, total)
	err = sft.downloadDirectory(sessionID, localPath, remotePath, tracker)
	tracker.Stop(err)
	return err
}

// downloadDirectory recursively downloads a remote directory to the specified local path
func (sft *Manager) downloadDirectory(sessionID string, localPath, remotePath string, tracker *transfer.Tracker) error {
	sft.logger.Debug("downloadDirectory", zap.String("sessionID", sessionID), zap.String("localPath", localPath), zap.String("remotePath", remotePath))

	// Validate inputs and get clients
	ftpClient, err := sft.getSftpClient(sessionID)
	if err != nil {
		return err
	}

	if tracker == nil {
		return fmt.Errorf(ErrTrackerRequired)
	}

	// Validate remote directory
	info, err := ftpClient.Stat(remotePath)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory")
	}

	// Create local directory
	localDir := filepath.Join(localPath, filepath.Base(remotePath))
	if err := sft.ensureLocalDirExists(localDir); err != nil {
		return err
	}

	// Process directory contents
	return sft.processRemoteDirectoryEntries(sessionID, localDir, remotePath, tracker)
}

// ensureLocalDirExists creates the local directory if it doesn't exist
func (sft *Manager) ensureLocalDirExists(localDir string) error {
	return os.MkdirAll(localDir, 0755)
}

// processRemoteDirectoryEntries processes all entries in a remote directory
func (sft *Manager) processRemoteDirectoryEntries(sessionID, localDir, remotePath string, tracker *transfer.Tracker) error {
	ftpClient, err := sft.getSftpClient(sessionID)
	if err != nil {
		return err
	}

	entries, err := ftpClient.ReadDir(remotePath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		remoteEntryPath := joinRemotePath(remotePath, entry.Name())
		localEntryPath := filepath.Join(localDir, entry.Name())

		if entry.IsDir() {
			// Recursively download subdirectory
			err = sft.downloadDirectory(sessionID, localDir, remoteEntryPath, tracker)
			if err != nil {
				return err
			}
		} else {
			// Download file
			err = sft.downloadFile(sessionID, localEntryPath, remoteEntryPath, tracker)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// DeleteFile deletes a file or directory at the specified path
func (sft *Manager) DeleteFile(sessionID string, path string) error {
	sft.logger.Debug("DeleteFile", zap.String("sessionID", sessionID), zap.String("path", path))
	ftpClient, err := sft.getSftpClient(sessionID)
	if err != nil {
		return err
	}

	// Check if it's a directory
	info, err := ftpClient.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		// It's a directory, recursively delete contents
		return ftpClient.RemoveAll(path)
	} else {
		// It's a file
		return ftpClient.Remove(path)
	}
}

// RenameFile renames a file or directory at the specified path
func (sft *Manager) RenameFile(sessionID string, oldPath, newPath string) error {
	sft.logger.Debug("RenameFile", zap.String("sessionID", sessionID), zap.String("oldPath", oldPath), zap.String("newPath", newPath))
	ftpClient, err := sft.getSftpClient(sessionID)
	if err != nil {
		return err
	}
	if _, err := ftpClient.Stat(newPath); err == nil {
		return fmt.Errorf("%s exists", newPath)
	}

	return ftpClient.Rename(oldPath, newPath)
}

// GetWd returns the current working directory
func (sft *Manager) GetWd(sessionID string) (string, error) {
	sft.logger.Debug("GetWd", zap.String("sessionID", sessionID))
	ftpClient, err := sft.getSftpClient(sessionID)
	if err != nil {
		return "", err
	}
	return ftpClient.Getwd()
}

func (sft *Manager) CloseSession(sessionID string) {
	if val, ok := sft.clients.LoadAndDelete(sessionID); ok {
		if err := val.(*sftp.Client).Close(); err != nil {
			sft.logger.Warn("Error closing SFTP connection:", zap.Error(err))
		}
	}
}

func (sft *Manager) CloseAll() {
	sft.clients.Range(func(key, value any) bool {
		if err := value.(*sftp.Client).Close(); err != nil {
			sft.logger.Warn("Error closing SFTP connection:", zap.Error(err))
		}
		sft.clients.Delete(key)
		return true
	})
}

// CreateFile creates a new file at the specified path
func (sft *Manager) CreateFile(sessionID string, path string) error {
	sft.logger.Debug("CreateFile", zap.String("sessionID", sessionID), zap.String("path", path))
	ftpClient, err := sft.getSftpClient(sessionID)
	if err != nil {
		return err
	}
	// file exists return error
	if _, err := ftpClient.Stat(path); err == nil {
		return fmt.Errorf("file already exists")
	}
	_, err = ftpClient.Create(path)
	return err
}

// CreateDirectory creates a new directory at the specified path
func (sft *Manager) CreateDirectory(sessionID string, path string) error {
	sft.logger.Debug("CreateDirectory", zap.String("sessionID", sessionID), zap.String("path", path))
	ftpClient, err := sft.getSftpClient(sessionID)
	if err != nil {
		return err
	}
	// directory exists return error
	if _, err := ftpClient.Stat(path); err == nil {
		return fmt.Errorf("directory already exists")
	}
	return ftpClient.MkdirAll(path)
}

// CancelTransfer cancels an ongoing transfer by its ID
func (sft *Manager) CancelTransfer(transferID string) error {
	sft.logger.Debug("CancelTransfer", zap.String("transferID", transferID))
	return sft.transfers.Cancel(transferID)
}

func (sft *Manager) calcLocalDirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

func (sft *Manager) calcRemoteDirSize(sessionID, remotePath string) (int64, error) {
	ftpClient, err := sft.getSftpClient(sessionID)
	if err != nil {
		return 0, err
	}
	var size int64
	var walk func(string) error
	walk = func(path string) error {
		entries, err := ftpClient.ReadDir(path)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if entry.IsDir() {
				if err := walk(joinRemotePath(path, entry.Name())); err != nil {
					return err
				}
			} else {
				size += entry.Size()
			}
		}
		return nil
	}
	err = walk(remotePath)
	return size, err
}
