package services

import (
	"path/filepath"

	"github.com/ilaziness/vexo/internal/bookmark"
	"github.com/ilaziness/vexo/internal/command"
	"github.com/ilaziness/vexo/internal/config"
	"github.com/ilaziness/vexo/internal/database"
	"github.com/ilaziness/vexo/internal/secret"
	"github.com/ilaziness/vexo/internal/sftp"
	"github.com/ilaziness/vexo/internal/ssh"
	"github.com/ilaziness/vexo/internal/termws"
	"github.com/ilaziness/vexo/internal/transfer"
	"github.com/ilaziness/vexo/internal/tunnel"
	"github.com/wailsapp/wails/v3/pkg/application"
	"go.uber.org/zap"
)

func RegisterServices(a *application.App, mainWindow *application.WebviewWindow, log *zap.Logger, commandsJSON []byte) error {
	if log == nil {
		log = zap.NewNop()
	}
	windows := NewWindows(a, mainWindow)

	cfgStore, err := config.Load(log)
	if err != nil {
		return err
	}
	vault := secret.NewVault()
	configService := NewConfigService(a, windows, cfgStore, vault)

	db := database.NewDatabase(cfgStore.UserDataDir(), log)
	if err := db.Initialize(); err != nil {
		return err
	}

	prompter := &hostKeyPrompter{app: a}
	sshMgr := ssh.NewManager(log, filepath.Join(cfgStore.UserDataDir(), "known_hosts"), prompter)
	transfers := transfer.NewRegistry(func(p transfer.ProgressData) {
		a.Event.Emit(EventProgress, p)
	})
	sftpMgr := sftp.NewManager(log, sshMgr, transfers)
	tunnelMgr := tunnel.NewManager(log, sshMgr)

	sshService := NewSSHService(a, sshMgr, sftpMgr, tunnelMgr)
	termWS := termws.NewServer(log, sshMgr, func(id string) { _ = sshService.CloseByID(id) })
	sshMgr.SetOnClose(func(id string) { _ = sshService.CloseByID(id) })

	bookmarks := bookmark.New(log, db, configService.getPasswordWithPrompt, func() {
		a.Event.Emit(EventBookmarkUpdate, BookmarkUpdateMsg)
	}, vault.Clear)
	sshService.bind(termWS, bookmarks)
	bookmarkService := NewBookmarkService(a, bookmarks, sshService)

	commandSvc := command.New(log, db, commandsJSON, sshService.SendToSession)
	commandService := NewCommandService(windows, commandSvc)
	syncService := NewSyncService(configService, db, log)
	toolService := NewToolService(windows)
	aiService := NewAIService(a, log, configService, sshService, db)
	appService := NewAppService(a, windows, termWS)
	tunnelService := NewSSHTunnelService(tunnelMgr)

	a.RegisterService(application.NewService(appService))
	a.RegisterService(application.NewService(sshService))
	a.RegisterService(application.NewService(NewSftpService(a, sftpMgr)))
	a.RegisterService(application.NewService(configService))
	a.RegisterService(application.NewService(bookmarkService))
	a.RegisterService(application.NewService(tunnelService))
	a.RegisterService(application.NewService(commandService))
	a.RegisterService(application.NewService(syncService))
	a.RegisterService(application.NewService(toolService))
	a.RegisterService(application.NewService(aiService))

	if err := termWS.Start(); err != nil {
		return err
	}

	a.OnShutdown(func() {
		log.Debug("run app OnShutdown...")
		termWS.Stop()
		sshService.Close()
		if err := db.Close(); err != nil {
			log.Error("close db failed", zap.Error(err))
		}
	})
	return nil
}
