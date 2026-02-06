package services

import (
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

var AppInstance *App

type App struct {
	app            *application.App
	MainWindow     *application.WebviewWindow
	SettingWindow  *application.WebviewWindow
	BookmarkWindow *application.WebviewWindow
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
	AppInstance.SettingWindow = newSettingWindow()
	AppInstance.BookmarkWindow = newBookmarkWindow()
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
	window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
		event.Cancel() // 取消关闭事件
		window.Hide()
	})
	return window
}

func newBookmarkWindow() *application.WebviewWindow {
	window := AppInstance.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:                      "书签管理",
		URL:                        "/#/bookmark",
		Width:                      1200,
		Height:                     800,
		MinWidth:                   1200,
		MinHeight:                  800,
		DefaultContextMenuDisabled: true,
		Hidden:                     true,
		Frameless:                  true,
	})
	window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
		event.Cancel()
		window.Hide()
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
