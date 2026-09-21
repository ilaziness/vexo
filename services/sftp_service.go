package services

import (
	"path/filepath"

	"github.com/ilaziness/vexo/internal/sftp"
	"github.com/ilaziness/vexo/internal/transfer"
	"github.com/wailsapp/wails/v3/pkg/application"
)

const (
	EventProgress        = "eventProgress"
	TransferTypeUpload   = transfer.TypeUpload
	TransferTypeDownload = transfer.TypeDownload
)

func init() {
	application.RegisterEvent[ProgressData](EventProgress)
}

type ProgressData = transfer.ProgressData
type FileInfo = sftp.FileInfo

type SftpService struct {
	app *application.App
	mgr *sftp.Manager
}

func NewSftpService(app *application.App, mgr *sftp.Manager) *SftpService {
	return &SftpService{app: app, mgr: mgr}
}

func (s *SftpService) ListFiles(sessionID, path string, showHidden bool) ([]FileInfo, error) {
	return s.mgr.ListFiles(sessionID, path, showHidden)
}
func (s *SftpService) GetFileInfo(sessionID, path string) (FileInfo, error) {
	return s.mgr.GetFileInfo(sessionID, path)
}
func (s *SftpService) UploadFileDialog(sessionID, remotePath string) error {
	localPath, err := s.app.Dialog.OpenFile().SetTitle("选择文件").PromptForSingleSelection()
	if err != nil || localPath == "" {
		return err
	}
	return s.mgr.UploadFile(sessionID, localPath, remotePath)
}
func (s *SftpService) DownloadFileDialog(sessionID, remotePathFile string) error {
	localPathFile, err := s.app.Dialog.SaveFile().
		SetMessage("保存文件").SetFilename(filepath.Base(remotePathFile)).
		CanCreateDirectories(true).PromptForSingleSelection()
	if err != nil {
		return err
	}
	return s.mgr.DownloadFile(sessionID, localPathFile, remotePathFile)
}
func (s *SftpService) UploadDirectoryDialog(sessionID, remotePath string) error {
	localPath, err := s.app.Dialog.OpenFile().SetTitle("选择目录").
		CanChooseDirectories(true).CanChooseFiles(false).PromptForSingleSelection()
	if err != nil || localPath == "" {
		return err
	}
	return s.mgr.UploadDirectory(sessionID, localPath, remotePath)
}
func (s *SftpService) DownloadDirectoryDialog(sessionID, remotePath string) error {
	localPath, err := s.app.Dialog.OpenFile().SetTitle("选择目录").
		CanChooseDirectories(true).CanChooseFiles(false).PromptForSingleSelection()
	if err != nil {
		return err
	}
	return s.mgr.DownloadDirectory(sessionID, localPath, remotePath)
}
func (s *SftpService) UploadPaths(sessionID, remotePath string, localPaths []string) error {
	return s.mgr.UploadPaths(sessionID, remotePath, localPaths)
}
func (s *SftpService) DeleteFile(sessionID, path string) error {
	return s.mgr.DeleteFile(sessionID, path)
}
func (s *SftpService) RenameFile(sessionID, oldPath, newPath string) error {
	return s.mgr.RenameFile(sessionID, oldPath, newPath)
}
func (s *SftpService) GetWd(sessionID string) (string, error) {
	return s.mgr.GetWd(sessionID)
}
func (s *SftpService) CreateFile(sessionID, path string) error {
	return s.mgr.CreateFile(sessionID, path)
}
func (s *SftpService) CreateDirectory(sessionID, path string) error {
	return s.mgr.CreateDirectory(sessionID, path)
}
func (s *SftpService) CancelTransfer(transferID string) error {
	return s.mgr.CancelTransfer(transferID)
}
