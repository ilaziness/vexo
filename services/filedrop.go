package services

import (
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

const EventSftpFilesDropped = "eventSftpFilesDropped"

type SftpFilesDroppedData struct {
	SessionID  string   `json:"sessionID"`
	LocalPaths []string `json:"localPaths"`
}

func init() {
	application.RegisterEvent[SftpFilesDroppedData](EventSftpFilesDropped)
}

func RegisterFileDropHandler(window *application.WebviewWindow) {
	window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
		files := event.Context().DroppedFiles()
		if len(files) == 0 {
			return
		}
		details := event.Context().DropTargetDetails()
		if details == nil || details.ElementID == "" {
			return
		}
		const prefix = "sftp-drop-"
		if !strings.HasPrefix(details.ElementID, prefix) {
			return
		}
		sessionID := strings.TrimPrefix(details.ElementID, prefix)
		if sessionID == "" {
			return
		}
		if app == nil {
			return
		}
		app.Event.Emit(EventSftpFilesDropped, SftpFilesDroppedData{
			SessionID:  sessionID,
			LocalPaths: files,
		})
	})
}
