package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/sftp"
	"go.uber.org/zap"
)

// FileInfo represents file information for SFTP operations
type FileInfo struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	Mode    os.FileMode `json:"mode"`
	ModTime time.Time `json:"modTime"`
	IsDir   bool      `json:"isDir"`
}

// Convert os.FileInfo to our custom FileInfo
func convertFileInfo(osInfo os.FileInfo) FileInfo {
	return FileInfo{
		Name:    osInfo.Name(),
		Size:    osInfo.Size(),
		Mode:    osInfo.Mode(),
		ModTime: osInfo.ModTime(),
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

// Connect connects to SFTP using an SSH session ID
func (sft *SftpService) Connect(sshSessionID string) error {
	sshS, ok := GetSSHSession(sshSessionID)
	if !ok {
		return fmt.Errorf("ssh session not found")
	}
	ftpClient, err := sftp.NewClient(sshS.client)
	if err != nil {
		return err
	}
	sft.ftpClient = ftpClient
	return nil
}

// ListFiles lists files and directories in the specified path
func (sft *SftpService) ListFiles(path string) ([]FileInfo, error) {
	if sft.ftpClient == nil {
		return nil, fmt.Errorf("sftp client not connected")
	}

	entries, err := sft.ftpClient.ReadDir(path)
	if err != nil {
		return nil, err
	}

	return convertFileInfos(entries), nil
}

// GetFileInfo returns information about a file or directory
func (sft *SftpService) GetFileInfo(path string) (FileInfo, error) {
	if sft.ftpClient == nil {
		return FileInfo{}, fmt.Errorf("sftp client not connected")
	}

	info, err := sft.ftpClient.Stat(path)
	if err != nil {
		return FileInfo{}, err
	}

	return convertFileInfo(info), nil
}

// UploadFile uploads a local file to the specified remote path
func (sft *SftpService) UploadFile(localPath, remotePath string) error {
	if sft.ftpClient == nil {
		return fmt.Errorf("sftp client not connected")
	}

	// Open the local file for reading
	localFile, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer localFile.Close()

	// Create the remote file for writing
	remoteFile, err := sft.ftpClient.Create(remotePath)
	if err != nil {
		return err
	}
	defer remoteFile.Close()

	// Copy the local file content to the remote file
	_, err = io.Copy(remoteFile, localFile)
	if err != nil {
		return err
	}

	return nil
}

// DownloadFile downloads a remote file to the specified local path
func (sft *SftpService) DownloadFile(remotePath, localPath string) error {
	if sft.ftpClient == nil {
		return fmt.Errorf("sftp client not connected")
	}

	// Open the remote file for reading
	remoteFile, err := sft.ftpClient.Open(remotePath)
	if err != nil {
		return err
	}
	defer remoteFile.Close()

	// Create the local file for writing
	localFile, err := os.Create(localPath)
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

// DeleteFile deletes a file or directory at the specified path
func (sft *SftpService) DeleteFile(path string) error {
	if sft.ftpClient == nil {
		return fmt.Errorf("sftp client not connected")
	}

	// Check if it's a directory
	info, err := sft.ftpClient.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		// It's a directory, recursively delete contents
		return sft.deleteDirectory(path)
	} else {
		// It's a file
		return sft.ftpClient.Remove(path)
	}
}

// deleteDirectory recursively deletes a directory and its contents
func (sft *SftpService) deleteDirectory(path string) error {
	entries, err := sft.ftpClient.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			err = sft.deleteDirectory(entryPath)
			if err != nil {
				return err
			}
		} else {
			err = sft.ftpClient.Remove(entryPath)
			if err != nil {
				return err
			}
		}
	}

	// Remove the empty directory
	return sft.ftpClient.RemoveDirectory(path)
}

// RenameFile renames a file or directory at the specified path
func (sft *SftpService) RenameFile(oldPath, newPath string) error {
	if sft.ftpClient == nil {
		return fmt.Errorf("sftp client not connected")
	}

	return sft.ftpClient.Rename(oldPath, newPath)
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