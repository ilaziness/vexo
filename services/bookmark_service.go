package services

import (
	"fmt"

	"github.com/ilaziness/vexo/internal/bookmark"
	"github.com/wailsapp/wails/v3/pkg/application"
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
	hops, err := bs.core.ResolveHopsForBookmark(bookmarkID)
	if err != nil {
		return "", err
	}
	return bs.ssh.connectHops(hops)
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
	ep, jumpID, err := bs.core.PrepareTest(b)
	if err != nil {
		return err
	}
	hops, err := bs.core.ResolveHops(ep, jumpID)
	if err != nil {
		return err
	}
	return bs.ssh.testHops(hops)
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
