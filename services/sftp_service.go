package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/pkg/sftp"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

var sftpClient = new(sync.Map)

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

	entries, err := ftpClient.ReadDir(path)
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

// UploadFile uploads a local file to the specified remote path.
// localPathFile is local file path with file full name.
// remoteDir is remote directory without file name.
func (sft *SftpService) UploadFile(sessionID string, localPathFile, remoteDir string) error {
	Logger.Debug("UploadFile", zap.String("sessionID", sessionID), zap.String("localPathFile", localPathFile), zap.String("remoteDir", remoteDir))
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

	// Create the remote file for writing
	remoteFile, err := ftpClient.Create(filepath.Join(remoteDir, filepath.Base(localPathFile)))
	if err != nil {
		return err
	}
	defer remoteFile.Close()

	// Copy the local file content to the remote file
	_, err = io.Copy(remoteFile, localFile)
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
	if localPath == "" {
		return fmt.Errorf("no file selected")
	}
	Logger.Info("UploadFileDialog", zap.String("localPath", localPath), zap.String("remotePath", remotePath))
	return sft.UploadFile(sessionID, localPath, remotePath)
}

// DownloadFile downloads a remote file to the specified local path
func (sft *SftpService) DownloadFile(sessionID string, localPath, remotePath string) error {
	Logger.Debug("DownloadFile", zap.String("sessionID", sessionID), zap.String("localPath", localPath), zap.String("remotePath", remotePath))
	ftpClient, err := sft.getSftpClient(sessionID)
	if err != nil {
		return err
	}

	// Open the remote file for reading
	remoteFile, err := ftpClient.Open(remotePath)
	if err != nil {
		return err
	}
	defer remoteFile.Close()

	// Create the local file for writing
	localFile, err := os.Create(filepath.Join(localPath, filepath.Base(remotePath)))
	if err != nil {
		return err
	}
	defer localFile.Close()

	// Copy the remote file content to the local file
	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		return err
	}

	return nil
}

// DownloadFileDialog opens a save dialog to select where to download a remote file
func (sft *SftpService) DownloadFileDialog(sessionID string, remotePath string) error {
	Logger.Debug("DownloadFileDialog", zap.String("sessionID", sessionID), zap.String("remotePath", remotePath))
	// 选择保存位置
	localPath, err := app.Dialog.SaveFile().
		SetMessage("保存文件").
		SetFilename(filepath.Base(remotePath)).
		PromptForSingleSelection()
	if err != nil {
		return err
	}
	if localPath == "" {
		return fmt.Errorf("no file selected")
	}
	Logger.Info("下载文件", zap.String("localPath", localPath), zap.String("remotePath", remotePath))
	return sft.DownloadFile(sessionID, localPath, remotePath)
}

// UploadDirectory recursively uploads a local directory to the specified remote path
// if localPath empty string, it will open a file dialog to select a directory
func (sft *SftpService) UploadDirectory(sessionID string, localPath, remotePath string) error {
	Logger.Debug("UploadDirectory", zap.String("sessionID", sessionID), zap.String("localPath", localPath), zap.String("remotePath", remotePath))
	var err error
	ftpClient, err := sft.getSftpClient(sessionID)
	if err != nil {
		return err
	}

	// if localPath is empty string
	if localPath == "" {
		localPath, err = app.Dialog.OpenFile().SetTitle("选择目录").PromptForSingleSelection()
		if err != nil {
			return err
		}
		if localPath == "" {
			return fmt.Errorf("no directory selected")
		}
	}
	Logger.Info("上传目录", zap.String("localPath", localPath), zap.String("remotePath", remotePath))

	// Stat the local directory
	info, err := os.Stat(localPath)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("path is not a directory")
	}

	// Create the remote directory if it doesn't exist
	info, err = ftpClient.Stat(remotePath)
	if err != nil {
		if os.IsNotExist(err) {
			err = ftpClient.MkdirAll(remotePath)
			if err != nil {
				return err
			}
		}
	}

	// Read the directory contents
	entries, err := os.ReadDir(localPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		localEntryPath := filepath.Join(localPath, entry.Name())
		remoteEntryPath := remotePath

		if entry.IsDir() {
			// Recursively upload subdirectory
			err = sft.UploadDirectory(sessionID, localEntryPath, remoteEntryPath)
			if err != nil {
				return err
			}
		} else {
			// Upload file
			err = sft.UploadFile(sessionID, localEntryPath, remoteEntryPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DownloadDirectory recursively downloads a remote directory to the specified local path
// if localPath empty string, it will open a file dialog to select a directory
func (sft *SftpService) DownloadDirectory(sessionID string, localPath, remotePath string) error {
	Logger.Debug("DownloadDirectory", zap.String("sessionID", sessionID), zap.String("localPath", localPath), zap.String("remotePath", remotePath))
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

	if localPath == "" {
		localPath, err = app.Dialog.SaveFile().
			SetMessage("保存目录").
			PromptForSingleSelection()
		if err != nil {
			return err
		}
	}
	if localPath == "" {
		return fmt.Errorf("no directory selected")
	}
	Logger.Info("下载目录", zap.String("localPath", localPath), zap.String("remotePath", remotePath))

	// Read the directory contents
	entries, err := ftpClient.ReadDir(remotePath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		remoteEntryPath := filepath.Join(remotePath, entry.Name())
		localEntryPath := filepath.Join(localPath, entry.Name())

		if entry.IsDir() {
			// Recursively download subdirectory
			err = sft.DownloadDirectory(sessionID, remoteEntryPath, localEntryPath)
			if err != nil {
				return err
			}
		} else {
			// Download file
			err = sft.DownloadFile(sessionID, remoteEntryPath, localEntryPath)
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
