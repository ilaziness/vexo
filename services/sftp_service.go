package services

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/pkg/sftp"
	"github.com/wailsapp/wails/v3/pkg/application"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"

	"github.com/ilaziness/vexo/internal/utils"
)

const (
	EventProgress = "eventProgress"

	TransferTypeUpload   = "upload"
	TransferTypeDownload = "download"
)

func init() {
	application.RegisterEvent[ProgressData](EventProgress)
}

var sftpClient = new(sync.Map)

type ProgressData struct {
	ID           string // 本地传输ID
	SessionID    string
	TransferType string // 传输类型
	LocalFile    string // 本地文件/目录，包含完整路径
	RemoteFile   string // 远程文件/目录，包含完整路径
	TotalSize    int64
	Rate         int
}

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

type SftpService struct {
	ftpClient *sftp.Client
}

func NewSftpService() *SftpService {
	return &SftpService{}
}

// joinRemotePath joins remote path elements with forward slash
// SFTP servers (typically Unix/Linux) always use forward slash as path separator
func joinRemotePath(elem ...string) string {
	return path.Join(elem...)
}

// getSftpClient 通过 sessionID 获取对应的 SFTP 客户端
func (sft *SftpService) getSftpClient(sessionID string) (*sftp.Client, error) {
	// 从全局 sshSession 获取对应的连接
	connVal, ok := sftpClient.Load(sessionID)
	if !ok {
		return nil, fmt.Errorf("SSH session with ID %s not found", sessionID)
	}
	conn := connVal.(*SftpService)
	return conn.ftpClient, nil
}

// Connect connects to SFTP using an SSH session ID
func (sft *SftpService) Connect(sshClient *ssh.Client) error {
	ftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return err
	}
	sft.ftpClient = ftpClient
	return nil
}

// ListFiles lists files and directories in the specified path
func (sft *SftpService) ListFiles(sessionID string, path string, showHidden bool) ([]FileInfo, error) {
	Logger.Debug("ListFiles", zap.String("sessionID", sessionID), zap.String("path", path), zap.Bool("showHidden", showHidden))
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
func (sft *SftpService) GetFileInfo(sessionID string, path string) (FileInfo, error) {
	Logger.Debug("GetFileInfo", zap.String("sessionID", sessionID), zap.String("path", path))
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
func (sft *SftpService) uploadFile(sessionID string, localPathFile, remoteDir string, tracker *transferTracker) error {
	Logger.Debug("uploadFile", zap.String("sessionID", sessionID), zap.String("localPathFile", localPathFile), zap.String("remoteDir", remoteDir))
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
	Logger.Debug("uploadFile Create File", zap.String("file", remoteFilePath))
	remoteFile, err := ftpClient.Create(remoteFilePath)
	if err != nil {
		return err
	}
	defer remoteFile.Close()

	if tracker == nil {
		return fmt.Errorf("tracker is required")
	}

	progressWriter := &progressWriter{writer: remoteFile, tracker: tracker}

	// Copy the local file content to the remote file
	_, err = io.Copy(progressWriter, localFile)
	return err
}

// UploadFileDialog opens a file dialog to select a file to upload
func (sft *SftpService) UploadFileDialog(sessionID string, remotePath string) error {
	Logger.Debug("UploadFileDialog", zap.String("sessionID", sessionID), zap.String("remotePath", remotePath))
	// 选择文件
	localPath, err := app.Dialog.OpenFile().SetTitle("选择文件").PromptForSingleSelection()
	if err != nil {
		return err
	}
	Logger.Info("UploadFileDialog", zap.String("localPath", localPath), zap.String("remotePath", remotePath))
	// Calculate file size and create tracker
	info, err := os.Stat(localPath)
	if err != nil {
		return err
	}
	tracker := newTransferTracker(sessionID, TransferTypeUpload, localPath, joinRemotePath(remotePath, filepath.Base(localPath)), filepath.Base(localPath), info.Size())
	tracker.startProgress()
	return sft.uploadFile(sessionID, localPath, remotePath, tracker)
}

// downloadFile downloads a remote file to the specified local path
// localPath local save file directory
// remotePathFile remote file
func (sft *SftpService) downloadFile(sessionID string, localPath, remotePathFile string, tracker *transferTracker) error {
	Logger.Debug("downloadFile", zap.String("sessionID", sessionID), zap.String("localPath", localPath), zap.String("remotePath", remotePathFile))
	ftpClient, err := sft.getSftpClient(sessionID)
	if err != nil {
		return err
	}

	// Open the remote file for reading
	remoteFile, err := ftpClient.Open(remotePathFile)
	if err != nil {
		return err
	}
	defer remoteFile.Close()

	// Create the local file for writing
	localFilePath := filepath.Join(localPath, filepath.Base(remotePathFile))
	localFile, err := os.Create(localFilePath)
	if err != nil {
		return err
	}
	defer localFile.Close()

	if tracker == nil {
		return fmt.Errorf("tracker is required")
	}

	progressWriter := &progressWriter{writer: localFile, tracker: tracker}

	// Copy the remote file content to the local file
	_, err = io.Copy(progressWriter, remoteFile)
	if err != nil {
		return err
	}

	return nil
}

// DownloadFileDialog opens a save dialog to select where to download a remote file
func (sft *SftpService) DownloadFileDialog(sessionID string, remotePathFile string) error {
	Logger.Debug("DownloadFileDialog", zap.String("sessionID", sessionID), zap.String("remotePath", remotePathFile))
	// 选择保存位置
	localPath, err := app.Dialog.SaveFile().
		SetMessage("保存文件").
		SetFilename(filepath.Base(remotePathFile)).
		CanCreateDirectories(true).
		PromptForSingleSelection()
	if err != nil {
		return err
	}
	// Get remote file info and create tracker
	ftpClient, err := sft.getSftpClient(sessionID)
	if err != nil {
		return err
	}
	info, err := ftpClient.Stat(remotePathFile)
	if err != nil {
		return err
	}
	tracker := newTransferTracker(sessionID, TransferTypeDownload, filepath.Join(localPath, filepath.Base(remotePathFile)), remotePathFile, filepath.Base(remotePathFile), info.Size())
	tracker.startProgress()
	return sft.downloadFile(sessionID, localPath, remotePathFile, tracker)
}

func (sft *SftpService) UploadDirectoryDialog(sessionID, remotePath string) error {
	localPath, err := app.Dialog.OpenFile().SetTitle("选择目录").
		CanChooseDirectories(true).
		CanChooseFiles(false).
		PromptForSingleSelection()
	if err != nil {
		return err
	}
	// Calculate directory size and create tracker
	total, err := sft.calcLocalDirSize(localPath)
	if err != nil {
		return err
	}
	tracker := newTransferTracker(sessionID, TransferTypeUpload, localPath, joinRemotePath(remotePath, filepath.Base(localPath)), filepath.Base(localPath), total)
	tracker.startProgress()
	return sft.uploadDirectory(sessionID, localPath, remotePath, tracker)
}

// uploadDirectory recursively uploads a local directory to the specified remote path
func (sft *SftpService) uploadDirectory(sessionID string, localPath, remotePath string, tracker *transferTracker) error {
	Logger.Debug("uploadDirectory", zap.String("sessionID", sessionID), zap.String("localPath", localPath), zap.String("remotePath", remotePath))
	var err error
	ftpClient, err := sft.getSftpClient(sessionID)
	if err != nil {
		return err
	}

	// Stat the local directory
	info, err := os.Stat(localPath)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("path is not a directory")
	}

	remoteDir := joinRemotePath(remotePath, filepath.Base(localPath))
	// Create the remote directory if it doesn't exist
	info, err = ftpClient.Stat(remoteDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = ftpClient.MkdirAll(remoteDir)
			if err != nil {
				return err
			}
		}
	}

	if tracker == nil {
		return fmt.Errorf("tracker is required")
	}

	// Read the directory contents
	entries, err := os.ReadDir(localPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		localEntryPath := filepath.Join(localPath, entry.Name())

		if entry.IsDir() {
			// Recursively upload subdirectory
			remoteEntryPath := joinRemotePath(remoteDir, entry.Name())
			err = sft.uploadDirectory(sessionID, localEntryPath, remoteEntryPath, tracker)
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

func (sft *SftpService) DownloadDirectoryDialog(sessionID, remotePath string) error {
	localPath, err := app.Dialog.OpenFile().SetTitle("选择目录").
		CanChooseDirectories(true).
		CanChooseFiles(false).
		PromptForSingleSelection()
	if err != nil {
		return err
	}
	// Calculate remote directory size and create tracker
	total, err := sft.calcRemoteDirSize(sessionID, remotePath)
	if err != nil {
		return err
	}
	tracker := newTransferTracker(sessionID, TransferTypeDownload, localPath, remotePath, filepath.Base(remotePath), total)
	tracker.startProgress()
	return sft.downloadDirectory(sessionID, localPath, remotePath, tracker)
}

// downloadDirectory recursively downloads a remote directory to the specified local path
func (sft *SftpService) downloadDirectory(sessionID string, localPath, remotePath string, tracker *transferTracker) error {
	Logger.Debug("downloadDirectory", zap.String("sessionID", sessionID), zap.String("localPath", localPath), zap.String("remotePath", remotePath))
	ftpClient, err := sft.getSftpClient(sessionID)
	if err != nil {
		return err
	}

	// Stat the remote directory
	info, err := ftpClient.Stat(remotePath)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("path is not a directory")
	}

	localDir := filepath.Join(localPath, filepath.Base(remotePath))
	// Create the local directory if it doesn't exist
	if err := os.MkdirAll(localDir, 0755); err != nil {
		return err
	}

	if tracker == nil {
		return fmt.Errorf("tracker is required")
	}

	// Read the directory contents
	entries, err := ftpClient.ReadDir(remotePath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		remoteEntryPath := joinRemotePath(remotePath, entry.Name())
		localEntryPath := filepath.Join(localDir, entry.Name())

		if entry.IsDir() {
			// Recursively download subdirectory
			err = sft.downloadDirectory(sessionID, localEntryPath, remoteEntryPath, tracker)
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
func (sft *SftpService) DeleteFile(sessionID string, path string) error {
	Logger.Debug("DeleteFile", zap.String("sessionID", sessionID), zap.String("path", path))
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
func (sft *SftpService) RenameFile(sessionID string, oldPath, newPath string) error {
	Logger.Debug("RenameFile", zap.String("sessionID", sessionID), zap.String("oldPath", oldPath), zap.String("newPath", newPath))
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
func (sft *SftpService) GetWd(sessionID string) (string, error) {
	Logger.Debug("GetWd", zap.String("sessionID", sessionID))
	ftpClient, err := sft.getSftpClient(sessionID)
	if err != nil {
		return "", err
	}
	return ftpClient.Getwd()
}

// Close closes the SFTP connection
func (sft *SftpService) Close() {
	var err error
	if sft.ftpClient != nil {
		err = sft.ftpClient.Close()
	}
	if err != nil {
		Logger.Warn("Error closing SFTP connection:", zap.Error(err))
	}
}

// CreateFile creates a new file at the specified path
func (sft *SftpService) CreateFile(sessionID string, path string) error {
	Logger.Debug("CreateFile", zap.String("sessionID", sessionID), zap.String("path", path))
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
func (sft *SftpService) CreateDirectory(sessionID string, path string) error {
	Logger.Debug("CreateDirectory", zap.String("sessionID", sessionID), zap.String("path", path))
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

func (sft *SftpService) calcLocalDirSize(path string) (int64, error) {
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

func (sft *SftpService) calcRemoteDirSize(sessionID, remotePath string) (int64, error) {
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

type transferTracker struct {
	sessionID    string
	id           string
	transferType string
	localFile    string
	remoteFile   string
	file         string
	total        int64
	transferred  int64
	mutex        sync.Mutex
}

func newTransferTracker(sessionID, transferType, localFile, remoteFile, fileName string, total int64) *transferTracker {
	return &transferTracker{
		sessionID:    sessionID,
		id:           utils.GenerateRandomID(),
		transferType: transferType,
		localFile:    localFile,
		remoteFile:   remoteFile,
		file:         fileName,
		total:        total,
	}
}

func (t *transferTracker) update(n int64) {
	t.mutex.Lock()
	t.transferred += n
	t.mutex.Unlock()
}

func (t *transferTracker) getRate() int {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.total == 0 {
		return 100
	}
	rate := int(t.transferred * 100 / t.total)
	rate = min(rate, 100)
	return rate
}

func (t *transferTracker) startProgress() {
	// Emit initial progress with rate 0
	app.Event.Emit(EventProgress, ProgressData{
		ID:           t.id,
		SessionID:    t.sessionID,
		TransferType: t.transferType,
		LocalFile:    t.localFile,
		RemoteFile:   t.remoteFile,
		TotalSize:    t.total,
		Rate:         0,
	})
	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			rate := t.getRate()
			app.Event.Emit(EventProgress, ProgressData{
				ID:           t.id,
				SessionID:    t.sessionID,
				TransferType: t.transferType,
				LocalFile:    t.localFile,
				RemoteFile:   t.remoteFile,
				TotalSize:    t.total,
				Rate:         rate,
			})
			if rate >= 100 {
				ticker.Stop()
				return
			}
		}
	}()
}

type progressWriter struct {
	writer  io.Writer
	tracker *transferTracker
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n, err := pw.writer.Write(p)
	if n > 0 {
		pw.tracker.update(int64(n))
	}
	return n, err
}
