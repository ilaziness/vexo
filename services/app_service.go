package services

import (
	"github.com/ilaziness/vexo/internal/buildinfo"
	"github.com/ilaziness/vexo/internal/termws"
	"github.com/ilaziness/vexo/internal/updater"
	"github.com/wailsapp/wails/v3/pkg/application"
)

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

type AppService struct {
	app        *application.App
	windows    *Windows
	termWS     *termws.Server
	mainWindow *application.WebviewWindow
}

func NewAppService(app *application.App, windows *Windows, termWS *termws.Server) *AppService {
	return &AppService{app: app, windows: windows, termWS: termWS, mainWindow: windows.Main}
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
	return cs.app.Dialog.OpenFile().SetTitle("选择目录").
		CanChooseDirectories(true).
		CanChooseFiles(false).
		CanCreateDirectories(true).
		PromptForSingleSelection()
}

func (cs *AppService) SelectFile() (string, error) {
	return cs.app.Dialog.OpenFile().SetTitle("选择文件").
		CanChooseDirectories(false).
		CanChooseFiles(true).
		PromptForSingleSelection()
}

func (cs *AppService) GetAppInfo() AppInfo {
	return AppInfo{
		Version:   buildinfo.Version,
		HomeURL:   "https://github.com/ilaziness/vexo",
		RunMode:   buildinfo.Mode,
		GitInfo:   buildinfo.GitInfo,
		BuildTime: buildinfo.BuildTime,
	}
}

func (cs *AppService) CheckUpdate() (hasNew bool, newVersion NewVersion) {
	ok, rel, err := updater.CheckUpdate("ilaziness/vexo", buildinfo.Version)
	if err != nil {
		return false, NewVersion{}
	}
	if ok {
		return true, NewVersion{Version: rel.Tag, Notes: rel.Body, URL: rel.HTMLURL}
	}
	return false, NewVersion{}
}

func (cs *AppService) GetWSAddr() string {
	if cs.termWS == nil {
		return ""
	}
	return cs.termWS.Addr()
}

func (cs *AppService) NewMainWindow() {
	cs.windows.NewMainWindow()
}
