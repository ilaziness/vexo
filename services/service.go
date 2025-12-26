package services

import "github.com/wailsapp/wails/v3/pkg/application"

var app *application.App

func RegisterServices(a *application.App) {
	app = a
	sshService := NewSSHService(app)
	sftpService := NewSftpService()
	app.RegisterService(application.NewService(sshService))
	app.RegisterService(application.NewService(sftpService))

	app.OnShutdown(func() {
		Logger.Sugar().Debugln("run app OnShutdown...")
		sshService.Close()
		sftpService.Close()
	})
}
