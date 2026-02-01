package services

import (
	"github.com/ilaziness/vexo/internal/updater"
	"github.com/wailsapp/wails/v3/pkg/application"
)

var app *application.App

type AppInfo struct {
	Version string
	HomeURL string
	RunMode string
}

type NewVersion struct {
	Version string
	Notes   string
	URL     string
}

func RegisterServices(a *application.App, mainWindow *application.WebviewWindow) {
	app = a
	appService := &AppService{
		mainWindow: mainWindow,
	}
	sshService := NewSSHService()
	sftpService := NewSftpService()
	configService := NewConfigService()
	bookmarkService := NewBookmarkService(configService)

	app.RegisterService(application.NewService(appService))
	app.RegisterService(application.NewService(sshService))
	app.RegisterService(application.NewService(sftpService))
	app.RegisterService(application.NewService(configService))
	app.RegisterService(application.NewService(bookmarkService))

	app.OnShutdown(func() {
		Logger.Sugar().Debugln("run app OnShutdown...")
		sshService.Close()
	})
}

type AppService struct {
	mainWindow *application.WebviewWindow
}

func (cs *AppService) MainWindowMin() {
	cs.mainWindow.Minimise()
}
func (cs *AppService) MainWindowMax() {
	if cs.mainWindow.IsMaximised() {
		cs.mainWindow.UnMaximise()
		return
	}
	cs.mainWindow.Maximise()
}
func (cs *AppService) MainWindowClose() {
	cs.mainWindow.Close()
}

func (cs *AppService) SelectDirectory() (string, error) {
	return app.Dialog.OpenFile().SetTitle("选择目录").
		CanChooseDirectories(true).
		CanChooseFiles(false).
		CanCreateDirectories(true).
		PromptForSingleSelection()
}

func (cs *AppService) SelectFile() (string, error) {
	return app.Dialog.OpenFile().SetTitle("选择文件").
		CanChooseDirectories(false).
		CanChooseFiles(true).
		PromptForSingleSelection()
}

func (cs *AppService) GetAppInfo() AppInfo {
	return AppInfo{
		Version: Version,
		HomeURL: "https://github.com/ilaziness/vexo",
		RunMode: Mode,
	}
}

func (cs *AppService) CheckUpdate() (hasNew bool, newVersion NewVersion) {
	ok, rel, err := updater.CheckUpdate("ilaziness/vexo", Version)
	if err != nil {
		return false, NewVersion{}
	}
	if ok {
		return true, NewVersion{
			Version: rel.Tag,
			Notes:   rel.Body,
			URL:     rel.HTMLURL,
		}
	}
	return false, NewVersion{}
}
