package services

import (
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

var AppInstance *App

type App struct {
	app           *application.App
	MainWindow    *application.WebviewWindow
	SettingWindow *application.WebviewWindow
	CommandWindow *application.WebviewWindow // 新增
	ToolWindow    *application.WebviewWindow // 工具窗口
}

func NewApp(a *application.App, mainWindow *application.WebviewWindow) {
	AppInstance = &App{
		app:        a,
		MainWindow: mainWindow,
	}
	mainWindow.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
		for _, window := range app.Window.GetAll() {
			window.Close()
		}
	})
}

func newSettingWindow() *application.WebviewWindow {
	window := AppInstance.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:                      "设置",
		URL:                        "/#/setting",
		Width:                      1200,
		Height:                     800,
		MinWidth:                   1200,
		MinHeight:                  800,
		DefaultContextMenuDisabled: true,
		Hidden:                     true,
		Frameless:                  true,
	})
	return window
}

func newCommandWindow() *application.WebviewWindow {
	window := AppInstance.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:                      "命令面板",
		URL:                        "/#/command",
		Width:                      1600,
		Height:                     900,
		MinWidth:                   1600,
		MinHeight:                  900,
		DefaultContextMenuDisabled: true,
		Hidden:                     true,
		Frameless:                  true,
	})
	return window
}

func newToolWindow() *application.WebviewWindow {
	window := AppInstance.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:                      "运维工具",
		URL:                        "/#/tools",
		Width:                      1200,
		Height:                     800,
		MinWidth:                   800,
		MinHeight:                  600,
		DefaultContextMenuDisabled: true,
		Hidden:                     true,
		Frameless:                  true,
	})
	return window
}

// SelectDirectory 打开目录选择对话框，返回选择的目录路径
func (cs *App) SelectDirectory() (string, error) {
	return app.Dialog.OpenFile().SetTitle("选择目录").
		CanChooseDirectories(true).
		CanChooseFiles(false).
		CanCreateDirectories(true).
		PromptForSingleSelection()
}

// SelectFile 打开文件选择对话框，返回选择的文件路径
func (cs *App) SelectFile() (string, error) {
	return app.Dialog.OpenFile().SetTitle("选择文件").
		CanChooseDirectories(false).
		CanChooseFiles(true).
		PromptForSingleSelection()
}
