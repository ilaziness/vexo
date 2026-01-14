package services

import "github.com/wailsapp/wails/v3/pkg/application"

var app *application.App

func RegisterServices(a *application.App, mainWindow *application.WebviewWindow) {
	app = a
	commonService := &WindowService{
		mainWindow: mainWindow,
	}
	sshService := NewSSHService()
	sftpService := NewSftpService()
	configService := NewConfigService()
	bookmarkService := NewBookmarkService(configService)

	app.RegisterService(application.NewService(commonService))
	app.RegisterService(application.NewService(sshService))
	app.RegisterService(application.NewService(sftpService))
	app.RegisterService(application.NewService(configService))
	app.RegisterService(application.NewService(bookmarkService))

	app.OnShutdown(func() {
		Logger.Sugar().Debugln("run app OnShutdown...")
		sshService.Close()
	})
}

type WindowService struct {
	mainWindow *application.WebviewWindow
}

func (cs *WindowService) MainWindowMin() {
	cs.mainWindow.Minimise()
}
func (cs *WindowService) MainWindowMax() {
	if cs.mainWindow.IsMaximised() {
		cs.mainWindow.UnMaximise()
		return
	}
	cs.mainWindow.Maximise()
}
func (cs *WindowService) MainWindowClose() {
	cs.mainWindow.Close()
}

func (cs *WindowService) SelectDirectory() (string, error) {
	return app.Dialog.OpenFile().SetTitle("选择目录").
		CanChooseDirectories(true).
		CanChooseFiles(false).
		CanCreateDirectories(true).
		PromptForSingleSelection()
}

func (cs *WindowService) SelectFile() (string, error) {
	return app.Dialog.OpenFile().SetTitle("选择文件").
		CanChooseDirectories(false).
		CanChooseFiles(true).
		PromptForSingleSelection()
}
