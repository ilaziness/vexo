package services

import (
	"fmt"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"

	"github.com/ilaziness/vexo/internal/bookmark"
	"github.com/ilaziness/vexo/internal/config"
	"github.com/ilaziness/vexo/internal/ssh"
)

const (
	EventBookmarkUpdate  = "eventBookmarkUpdate"
	EventConnectBookmark = "eventConnectBookmark"
	BookmarkUpdateMsg    = "bookmark update"
)

func init() {
	application.RegisterEvent[string](EventBookmarkUpdate)
	application.RegisterEvent[string](EventConnectBookmark)
}

type SSHBookmark = bookmark.Bookmark
type BookmarkGroup = bookmark.Group
type BookmarkListItem = bookmark.ListItem

type BookmarkService struct {
	app  *application.App
	core *bookmark.Service
	ssh  *SSHService
}

func NewBookmarkService(app *application.App, core *bookmark.Service, ssh *SSHService) *BookmarkService {
	return &BookmarkService{app: app, core: core, ssh: ssh}
}

func (bs *BookmarkService) ConnectBookmark(bookmarkID string) {
	bs.app.Event.Emit(EventConnectBookmark, bookmarkID)
}

func (bs *BookmarkService) ConnectBookmarkByID(bookmarkID string) (string, error) {
	b, err := bs.core.GetDecrypted(bookmarkID)
	if err != nil {
		return "", err
	}
	hops, err := bs.core.ResolveHopsBookmark(*b)
	if err != nil {
		return "", err
	}
	proxy, err := resolveDialProxy(b.ProxyMode, ssh.ProxyConfig{
		Type: b.ProxyType, Host: b.ProxyHost, Port: b.ProxyPort,
		User: b.ProxyUser, Password: b.ProxyPassword,
	}, bs.sshConfig())
	if err != nil {
		return "", err
	}
	return bs.ssh.connectHops(hops, proxy)
}

func (bs *BookmarkService) sshConfig() config.SSHConfig {
	if bs.ssh == nil || bs.ssh.getSSHConfig == nil {
		return config.SSHConfig{}
	}
	return bs.ssh.getSSHConfig()
}

func (bs *BookmarkService) ListBookmarks() ([]*BookmarkGroup, error) {
	return bs.core.ListGrouped()
}
func (bs *BookmarkService) GetAllBookmarks() ([]*BookmarkListItem, error) {
	return bs.core.ListItems()
}
func (bs *BookmarkService) GetBookmarkByID(id string) (*SSHBookmark, error) {
	return bs.core.GetMasked(id)
}
func (bs *BookmarkService) SaveBookmark(b SSHBookmark) (string, error) {
	if err := validateBookmarkProxy(b); err != nil {
		return "", err
	}
	return bs.core.Save(b)
}
func (bs *BookmarkService) DeleteBookmark(id string) error {
	return bs.core.Delete(id)
}
func (bs *BookmarkService) CopyBookmark(id string) (*SSHBookmark, error) {
	return bs.core.Copy(id)
}
func (bs *BookmarkService) AddGroup(name string) error {
	return bs.core.AddGroup(name)
}
func (bs *BookmarkService) UpdateGroup(oldName, newName, icon string) error {
	return bs.core.UpdateGroup(oldName, newName, icon)
}
func (bs *BookmarkService) DeleteGroup(name string) error {
	return bs.core.DeleteGroup(name)
}
func (bs *BookmarkService) TestConnection(b SSHBookmark) error {
	prep, err := bs.core.PrepareTest(b)
	if err != nil {
		return err
	}
	hops, err := bs.core.ResolveHops(prep.Endpoint, prep.JumpID)
	if err != nil {
		return err
	}
	proxy, err := resolveDialProxy(b.ProxyMode, ssh.ProxyConfig{
		Type: b.ProxyType, Host: b.ProxyHost, Port: b.ProxyPort,
		User: b.ProxyUser, Password: prep.ProxyPassword,
	}, bs.sshConfig())
	if err != nil {
		return err
	}
	return bs.ssh.testHops(hops, proxy)
}

func validateBookmarkProxy(b SSHBookmark) error {
	mode := strings.ToLower(strings.TrimSpace(b.ProxyMode))
	if mode == "" {
		mode = "inherit"
	}
	switch mode {
	case "inherit", "none":
		return nil
	case "custom":
		p := ssh.ProxyConfig{Type: b.ProxyType, Host: b.ProxyHost, Port: b.ProxyPort}
		return p.Validate("书签自定义代理")
	default:
		return fmt.Errorf("无效的代理模式: %s", b.ProxyMode)
	}
}

func (bs *BookmarkService) SaveAndConnect(b SSHBookmark) (string, error) {
	if err := bs.TestConnection(b); err != nil {
		return "", fmt.Errorf("连接测试失败: %w", err)
	}
	id, err := bs.SaveBookmark(b)
	if err != nil {
		return "", fmt.Errorf("保存书签失败: %w", err)
	}
	bs.ConnectBookmark(id)
	return id, nil
}
