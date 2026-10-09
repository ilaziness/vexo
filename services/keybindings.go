package services

import (
	"runtime"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

const EventShortcut = "eventShortcut"

const (
	ShortcutTerminalFind      = "terminal.find"
	ShortcutTerminalCopy      = "terminal.copy"
	ShortcutTerminalPaste     = "terminal.paste"
	ShortcutTerminalAddToChat = "terminal.addToChat"
	ShortcutTerminalToggleLog = "terminal.toggleLog"
	ShortcutTerminalClear     = "terminal.clear"
	ShortcutPanelAI           = "panel.ai"
	ShortcutPanelBookmarks    = "panel.bookmarks"
	ShortcutPanelCommand      = "panel.command"
	ShortcutPanelTools        = "panel.tools"
	ShortcutPanelSettings     = "panel.settings"
)

const (
	CategoryTerminal = "终端"
	CategoryPanel    = "面板"
)

type ShortcutEvent struct {
	Action string `json:"action"`
}

// KeyBindingInfo 内置快捷键说明，供设置页与前端提示展示。
type KeyBindingInfo struct {
	Action   string `json:"action"`
	Keys     string `json:"keys"`
	Label    string `json:"label"`
	Category string `json:"category"`
}

type bindingDef struct {
	action   string
	logical  string // CmdOrCtrl[+Shift]+Key 或 CmdOrCtrl+,
	label    string
	category string
	toggle   func(*Windows) // 非空则 Go 侧 Toggle；空则向聚焦窗口派发 action
}

func init() {
	application.RegisterEvent[ShortcutEvent](EventShortcut)
}

func builtinBindings() []bindingDef {
	return []bindingDef{
		{action: ShortcutTerminalFind, logical: "CmdOrCtrl+F", label: "查找", category: CategoryTerminal},
		{action: ShortcutTerminalCopy, logical: "CmdOrCtrl+Shift+C", label: "复制", category: CategoryTerminal},
		{action: ShortcutTerminalPaste, logical: "CmdOrCtrl+Shift+V", label: "粘贴", category: CategoryTerminal},
		{action: ShortcutTerminalAddToChat, logical: "CmdOrCtrl+Shift+L", label: "添加到聊天", category: CategoryTerminal},
		{action: ShortcutTerminalToggleLog, logical: "CmdOrCtrl+Shift+R", label: "开始/停止记录", category: CategoryTerminal},
		{action: ShortcutTerminalClear, logical: "CmdOrCtrl+Shift+K", label: "清屏", category: CategoryTerminal},
		{action: ShortcutPanelAI, logical: "CmdOrCtrl+Shift+A", label: "打开/关闭 AI 助手", category: CategoryPanel},
		{action: ShortcutPanelBookmarks, logical: "CmdOrCtrl+Shift+B", label: "打开/关闭书签管理", category: CategoryPanel},
		{action: ShortcutPanelCommand, logical: "CmdOrCtrl+Shift+P", label: "打开/关闭命令面板", category: CategoryPanel, toggle: (*Windows).ToggleCommand},
		{action: ShortcutPanelTools, logical: "CmdOrCtrl+Shift+T", label: "打开/关闭工具面板", category: CategoryPanel, toggle: (*Windows).ToggleTool},
		{action: ShortcutPanelSettings, logical: "CmdOrCtrl+,", label: "打开/关闭设置", category: CategoryPanel, toggle: (*Windows).ToggleSetting},
	}
}

// ListKeyBindings 返回当前平台的内置快捷键列表。
func ListKeyBindings() []KeyBindingInfo {
	defs := builtinBindings()
	out := make([]KeyBindingInfo, 0, len(defs))
	for _, d := range defs {
		out = append(out, KeyBindingInfo{
			Action:   d.action,
			Keys:     displayKeys(d.logical),
			Label:    d.label,
			Category: d.category,
		})
	}
	return out
}

// RegisterKeyBindings 注册应用内快捷键。
// 使用 DispatchWailsEvent 只投递给当前聚焦窗口（EmitEvent 会广播到所有窗口）。
// KeyBinding.Add 必须使用平台归一化键名，不能直接传 CmdOrCtrl。
func RegisterKeyBindings(app *application.App, windows *Windows) {
	for _, d := range builtinBindings() {
		def := d
		if def.toggle != nil {
			toggle := def.toggle
			// 必须离开 AcceleratorKeyPressed 回调后再建窗：同步 NewWithOptions
			// 会在 WebView2 COM 回调栈上嵌套初始化，导致卡住并超时退出。
			addKeyBinding(app, def.logical, func(application.Window) {
				go func() {
					application.InvokeAsync(func() {
						toggle(windows)
					})
				}()
			})
			continue
		}
		action := def.action
		addKeyBinding(app, def.logical, func(window application.Window) {
			window.DispatchWailsEvent(&application.CustomEvent{
				Name: EventShortcut,
				Data: ShortcutEvent{Action: action},
			})
		})
	}
}

// addKeyBinding 注册归一化后的加速键；Windows 上将 "," 同步注册为 OEM_COMMA
// （系统上报的键名，见 VirtualKeyCodes[0xBC]）。
func addKeyBinding(app *application.App, logical string, cb func(window application.Window)) {
	accel := normalizeAccel(logical)
	app.KeyBinding.Add(accel, cb)
	for _, alias := range accelAliases(accel) {
		app.KeyBinding.Add(alias, cb)
	}
}

func accelAliases(accel string) []string {
	if runtime.GOOS != "windows" || !strings.HasSuffix(accel, "+,") {
		return nil
	}
	return []string{strings.TrimSuffix(accel, ",") + "OEM_COMMA"}
}

// normalizeAccel 将逻辑加速键转为 KeyBinding.Add 可用的平台字符串。
// 修饰键顺序与 Wails accelerator.String() 一致：字母序 + 主键。
func normalizeAccel(logical string) string {
	parts := strings.Split(logical, "+")
	if len(parts) == 0 {
		return logical
	}
	key := parts[len(parts)-1]
	modSet := map[string]struct{}{}
	for _, p := range parts[:len(parts)-1] {
		switch strings.ToLower(p) {
		case "cmdorctrl", "cmd", "command":
			if runtime.GOOS == "darwin" {
				modSet["Cmd"] = struct{}{}
			} else {
				modSet["Ctrl"] = struct{}{}
			}
		case "ctrl", "control":
			modSet["Ctrl"] = struct{}{}
		case "shift":
			modSet["Shift"] = struct{}{}
		case "alt", "option":
			if runtime.GOOS == "darwin" {
				modSet["Option"] = struct{}{}
			} else {
				modSet["Alt"] = struct{}{}
			}
		}
	}
	order := []string{"Alt", "Cmd", "Ctrl", "Option", "Shift", "Win"}
	var mods []string
	for _, m := range order {
		if _, ok := modSet[m]; ok {
			mods = append(mods, m)
		}
	}
	// 与 Wails accelerator.String() 一致：主键大写（字母 / 命名键）。
	key = strings.ToUpper(key)
	if len(mods) == 0 {
		return key
	}
	return strings.Join(append(mods, key), "+")
}

func displayKeys(logical string) string {
	isMac := runtime.GOOS == "darwin"
	parts := strings.Split(logical, "+")
	if len(parts) == 0 {
		return logical
	}
	key := parts[len(parts)-1]
	var out []string
	for _, p := range parts[:len(parts)-1] {
		switch strings.ToLower(p) {
		case "cmdorctrl", "cmd", "command":
			if isMac {
				out = append(out, "⌘")
			} else {
				out = append(out, "Ctrl")
			}
		case "ctrl", "control":
			if isMac {
				out = append(out, "⌃")
			} else {
				out = append(out, "Ctrl")
			}
		case "shift":
			if isMac {
				out = append(out, "⇧")
			} else {
				out = append(out, "Shift")
			}
		case "alt", "option":
			if isMac {
				out = append(out, "⌥")
			} else {
				out = append(out, "Alt")
			}
		default:
			out = append(out, p)
		}
	}
	if isMac {
		return strings.Join(append(out, strings.ToUpper(key)), "")
	}
	return strings.Join(append(out, strings.ToUpper(key)), "+")
}
