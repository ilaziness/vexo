package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// SSHBookmark 书签连接信息结构
type SSHBookmark struct {
	GroupName  string `json:"group_name"`
	ID         string `json:"id"`
	Title      string `json:"title"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	PrivateKey string `json:"private_key"`
	User       string `json:"user"`
	Password   string `json:"password"`
}

// BookmarkGroup 书签分组结构
type BookmarkGroup struct {
	Name      string        `json:"name"`
	Bookmarks []SSHBookmark `json:"bookmarks"`
}

// BookmarkService 书签服务结构
type BookmarkService struct {
	cfgService   *ConfigService
	bookmarks    []*BookmarkGroup
	bookmarkFile string
	window       *application.WebviewWindow
}

// NewBookmarkService 创建新的书签服务实例
func NewBookmarkService(cfgs *ConfigService) *BookmarkService {
	bookmarkFile := filepath.Join(cfgs.Config.General.UserDataDir, "bookmarks.json")
	bs := &BookmarkService{
		cfgService:   cfgs,
		bookmarkFile: bookmarkFile,
	}
	// 初始化时加载书签数据到内存
	bs.loadBookmarksToMemory()
	return bs
}

// initializeDefaultBookmarks 初始化默认书签数据
func (bs *BookmarkService) initializeDefaultBookmarks() error {
	defaultGroup := &BookmarkGroup{
		Name:      "默认书签",
		Bookmarks: []SSHBookmark{},
	}

	bookmarkGroups := []*BookmarkGroup{defaultGroup}

	// 更新内存中的数据
	bs.bookmarks = bookmarkGroups

	return bs.saveBookmarks()
}

// ShowWindow show bookmark manage window
func (bs *BookmarkService) ShowWindow() {
	if bs.window != nil {
		bs.window.Focus()
		return
	}
	bs.window = app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "书签",
		URL:              "/bookmark",
		MinWidth:         1200,
		MinHeight:        700,
		BackgroundColour: application.NewRGB(27, 38, 54),
	})
	bs.window.Show()
}

// CloseWindow 关闭窗口
func (bs *BookmarkService) CloseWindow() {
	if bs.window != nil {
		bs.window.Close()
	}
}

// loadBookmarksToMemory 从文件加载书签数据到内存
func (bs *BookmarkService) loadBookmarksToMemory() {
	data, err := os.ReadFile(bs.bookmarkFile)
	if err != nil {
		if os.IsNotExist(err) {
			// 如果文件不存在，创建默认书签
			bs.initializeDefaultBookmarks()
			// 再次尝试读取
			data, err = os.ReadFile(bs.bookmarkFile)
			if err != nil {
				// 如果再次失败，创建默认数据
				bs.bookmarks = []*BookmarkGroup{
					{
						Name:      "默认书签",
						Bookmarks: []SSHBookmark{},
					},
				}
				return
			}
		} else {
			// 其他错误，创建默认数据
			bs.bookmarks = []*BookmarkGroup{
				{
					Name:      "默认书签",
					Bookmarks: []SSHBookmark{},
				},
			}
			return
		}
	}

	var groups []*BookmarkGroup
	if err := json.Unmarshal(data, &groups); err != nil {
		// 解析错误，创建默认数据
		bs.bookmarks = []*BookmarkGroup{
			{
				Name:      "默认书签",
				Bookmarks: []SSHBookmark{},
			},
		}
		return
	}

	// 确保至少有一个默认分组
	if len(groups) == 0 {
		groups = append(groups, &BookmarkGroup{
			Name:      "默认书签",
			Bookmarks: []SSHBookmark{},
		})
	}

	bs.bookmarks = groups
}

// loadBookmarks 从内存加载书签数据
func (bs *BookmarkService) loadBookmarks() []*BookmarkGroup {
	return bs.bookmarks
}

// saveBookmarks 保存书签数据到文件
func (bs *BookmarkService) saveBookmarks() error {
	// 确保目录存在
	bookmarkDir := filepath.Dir(bs.bookmarkFile)
	if err := os.MkdirAll(bookmarkDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(bs.bookmarks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(bs.bookmarkFile, data, 0644)
}

// ListBookmarks 列出所有书签
func (bs *BookmarkService) ListBookmarks() ([]*BookmarkGroup, error) {
	return bs.loadBookmarks(), nil
}

// AddGroup 新增分组
func (bs *BookmarkService) AddGroup(groupName string) error {
	bookmarkGroups := bs.loadBookmarks()

	// 检查分组是否已存在
	for _, group := range bookmarkGroups {
		if group.Name == groupName {
			return fmt.Errorf("分组 '%s' 已存在", groupName)
		}
	}

	newGroup := &BookmarkGroup{
		Name:      groupName,
		Bookmarks: []SSHBookmark{},
	}
	bookmarkGroups = append(bookmarkGroups, newGroup)

	// 更新内存中的数据
	bs.bookmarks = bookmarkGroups

	return bs.saveBookmarks()
}

// AddBookmark 新增SSH连接信息书签
func (bs *BookmarkService) AddBookmark(bookmark SSHBookmark) error {
	bookmarkGroups := bs.loadBookmarks()

	// 检查分组是否存在，如果不存在则使用默认分组
	groupExists := false
	for i, group := range bookmarkGroups {
		if group.Name == bookmark.GroupName {
			// 检查ID是否已存在
			for _, existingBookmark := range group.Bookmarks {
				if existingBookmark.ID == bookmark.ID {
					return fmt.Errorf("书签ID '%s' 已存在", bookmark.ID)
				}
			}

			bookmarkGroups[i].Bookmarks = append(group.Bookmarks, bookmark)
			groupExists = true
			break
		}
	}

	if !groupExists {
		// 如果分组不存在，添加到默认分组
		defaultGroupName := "默认书签"
		for i, group := range bookmarkGroups {
			if group.Name == defaultGroupName {
				// 检查ID是否已存在
				for _, existingBookmark := range group.Bookmarks {
					if existingBookmark.ID == bookmark.ID {
						return fmt.Errorf("书签ID '%s' 已存在", bookmark.ID)
					}
				}

				bookmarkGroups[i].Bookmarks = append(group.Bookmarks, bookmark)
				break
			}
		}
	}

	// 更新内存中的数据
	bs.bookmarks = bookmarkGroups

	return bs.saveBookmarks()
}

// DeleteGroup 删除分组
func (bs *BookmarkService) DeleteGroup(groupName string) error {
	bookmarkGroups := bs.loadBookmarks()

	// 不能删除默认分组
	if groupName == "默认书签" {
		return fmt.Errorf("不能删除默认分组")
	}

	// 查找并移除分组
	newGroups := []*BookmarkGroup{}
	removed := false

	for _, group := range bookmarkGroups {
		if group.Name != groupName {
			newGroups = append(newGroups, group)
		} else {
			removed = true
		}
	}

	if !removed {
		return fmt.Errorf("分组 '%s' 不存在", groupName)
	}

	// 更新内存中的数据
	bs.bookmarks = newGroups

	return bs.saveBookmarks()
}

// DeleteBookmark 删除书签
func (bs *BookmarkService) DeleteBookmark(bookmarkID string) error {
	bookmarkGroups := bs.loadBookmarks()

	bookmarkFound := false

	for i, group := range bookmarkGroups {
		for j, bookmark := range group.Bookmarks {
			if bookmark.ID == bookmarkID {
				// 从切片中移除书签
				bookmarks := bookmarkGroups[i].Bookmarks
				bookmarkGroups[i].Bookmarks = append(bookmarks[:j], bookmarks[j+1:]...)
				bookmarkFound = true
				break
			}
		}
		if bookmarkFound {
			break
		}
	}

	if !bookmarkFound {
		return fmt.Errorf("书签ID '%s' 不存在", bookmarkID)
	}

	// 更新内存中的数据
	bs.bookmarks = bookmarkGroups

	return bs.saveBookmarks()
}

// GetBookmarkByID 通过书签ID查找返回连接信息
func (bs *BookmarkService) GetBookmarkByID(bookmarkID string) (*SSHBookmark, error) {
	bookmarkGroups := bs.loadBookmarks()

	for _, group := range bookmarkGroups {
		for _, bookmark := range group.Bookmarks {
			if bookmark.ID == bookmarkID {
				return &bookmark, nil
			}
		}
	}

	return nil, fmt.Errorf("未找到ID为 '%s' 的书签", bookmarkID)
}

// UpdateBookmark 更新SSH连接信息书签
func (bs *BookmarkService) UpdateBookmark(bookmark SSHBookmark) error {
	bookmarkGroups := bs.loadBookmarks()

	bookmarkFound := false
	for i, group := range bookmarkGroups {
		for j, existingBookmark := range group.Bookmarks {
			if existingBookmark.ID == bookmark.ID {
				// 检查目标分组是否存在
				targetGroupExists := false
				for _, g := range bookmarkGroups {
					if g.Name == bookmark.GroupName {
						targetGroupExists = true
						break
					}
				}

				if !targetGroupExists {
					// 如果目标分组不存在，使用原分组
					bookmark.GroupName = existingBookmark.GroupName
				}

				// 更新书签信息
				bookmarkGroups[i].Bookmarks[j] = bookmark
				bookmarkFound = true
				break
			}
		}
		if bookmarkFound {
			break
		}
	}

	if !bookmarkFound {
		return fmt.Errorf("书签ID '%s' 不存在", bookmark.ID)
	}

	// 更新内存中的数据
	bs.bookmarks = bookmarkGroups

	return bs.saveBookmarks()
}

// UpdateGroupName 更新分组名称
func (bs *BookmarkService) UpdateGroupName(oldName, newName string) error {
	bookmarkGroups := bs.loadBookmarks()

	// 不能修改默认分组名称
	if oldName == "默认书签" {
		return fmt.Errorf("不能修改默认分组名称")
	}

	// 检查新分组名是否已存在
	for _, group := range bookmarkGroups {
		if group.Name == newName {
			return fmt.Errorf("分组 '%s' 已存在", newName)
		}
	}

	groupFound := false
	for i, group := range bookmarkGroups {
		if group.Name == oldName {
			bookmarkGroups[i].Name = newName
			// 同时更新该组内所有书签的分组名
			for j := range group.Bookmarks {
				bookmarkGroups[i].Bookmarks[j].GroupName = newName
			}
			groupFound = true
			break
		}
	}

	if !groupFound {
		return fmt.Errorf("分组 '%s' 不存在", oldName)
	}

	// 更新内存中的数据
	bs.bookmarks = bookmarkGroups

	return bs.saveBookmarks()
}
