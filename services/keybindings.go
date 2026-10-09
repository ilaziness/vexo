package services

import (
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
)

const EventShortcut = "eventShortcut"

const ShortcutTerminalFind = "terminal.find"

type ShortcutEvent struct {
	Action string `json:"action"`
}

func init() {
	application.RegisterEvent[ShortcutEvent](EventShortcut)
}

// RegisterKeyBindings 注册应用内快捷键。
// 使用 DispatchWailsEvent 只投递给当前聚焦窗口（EmitEvent 会广播到所有窗口）。
func RegisterKeyBindings(app *application.App) {
	emitFind := func(window application.Window) {
		window.DispatchWailsEvent(&application.CustomEvent{
			Name: EventShortcut,
			Data: ShortcutEvent{Action: ShortcutTerminalFind},
		})
	}
	if runtime.GOOS == "darwin" {
		app.KeyBinding.Add("Cmd+F", emitFind)
	} else {
		app.KeyBinding.Add("Ctrl+F", emitFind)
	}
}
