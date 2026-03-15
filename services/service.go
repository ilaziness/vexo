package services

import (
	"github.com/ilaziness/vexo/internal/database"
	"github.com/ilaziness/vexo/internal/updater"
	"github.com/wailsapp/wails/v3/pkg/application"
	"go.uber.org/zap"
)

var app *application.App
var ConfigSvc *ConfigService

type AppInfo struct {
	Version   string
	HomeURL   string
	RunMode   string
	GitInfo   string
	BuildTime string
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

	configService := NewConfigService()

	database.Logger = Logger
	db := database.NewDatabase(configService.Config.General.UserDataDir)
	skipCreateTables := !configService.CheckAndUpdateDBVersion()
	if err := db.Initialize(skipCreateTables); err != nil {
		Logger.Error("init db failed", zap.Error(err))
		return
	}

	sshService := NewSSHService()
	sftpService := NewSftpService()
	bookmarkService := NewBookmarkService(configService, db)
	sshTunnelService := NewSSHTunnelService(sshService)
	commandService := NewCommandService(sshService, db)

	ConfigSvc = configService

	app.RegisterService(application.NewService(appService))
	app.RegisterService(application.NewService(sshService))
	app.RegisterService(application.NewService(sftpService))
	app.RegisterService(application.NewService(configService))
	app.RegisterService(application.NewService(bookmarkService))
	app.RegisterService(application.NewService(sshTunnelService))
	app.RegisterService(application.NewService(commandService))

	wsService := NewWebSocketService(app, sshService)
	wsService.Start()

	app.OnShutdown(func() {
		Logger.Sugar().Debugln("run app OnShutdown...")
		wsService.Stop()
		sshService.Close()
		sshTunnelService.StopAll()

		if err := db.Close(); err != nil {
			Logger.Error("close db failed", zap.Error(err))
		}
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

// SelectDirectory 打开目录选择对话框，返回选择的目录路径
func (cs *AppService) SelectDirectory() (string, error) {
	return app.Dialog.OpenFile().SetTitle("选择目录").
		CanChooseDirectories(true).
		CanChooseFiles(false).
		CanCreateDirectories(true).
		PromptForSingleSelection()
}

// SelectFile 打开文件选择对话框，返回选择的文件路径
func (cs *AppService) SelectFile() (string, error) {
	return app.Dialog.OpenFile().SetTitle("选择文件").
		CanChooseDirectories(false).
		CanChooseFiles(true).
		PromptForSingleSelection()
}

func (cs *AppService) GetAppInfo() AppInfo {
	return AppInfo{
		Version:   Version,
		HomeURL:   "https://github.com/ilaziness/vexo",
		RunMode:   Mode,
		GitInfo:   GitInfo,
		BuildTime: BuildTime,
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

// GetWSAddr 获取 WebSocket 服务器地址
func (cs *AppService) GetWSAddr() string {
	return wsAddr
}
