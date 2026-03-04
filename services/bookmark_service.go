package services

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ilaziness/vexo/internal/secret"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

const (
	EventInputPasswort      = "eventInputPassword"
	EventInputPasswortClose = "eventInputPasswordClose"
	EventBookmarkUpdate     = "eventBookmarkUpdate"
	EventConnectBookmark    = "eventConnectBookmark"
)

func init() {
	application.RegisterEvent[string](EventInputPasswort)
	application.RegisterEvent[string](EventInputPasswortClose)
	application.RegisterEvent[string](EventBookmarkUpdate)
	application.RegisterEvent[string](EventConnectBookmark)
}

// SSHBookmark 书签连接信息结构
type SSHBookmark struct {
	GroupName          string `json:"group_name"`
	ID                 string `json:"id"`
	Title              string `json:"title"`
	Host               string `json:"host"`
	Port               int    `json:"port"`
	PrivateKey         string `json:"private_key"`
	PrivateKeyPassword string `json:"private_key_password"`
	User               string `json:"user"`
	Password           string `json:"password"`
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
	userPassword string
	passwordChan chan struct{} // 用于等待用户密码输入完成的信号channel
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

// ShowWindow show bookmark manage window
func (bs *BookmarkService) ShowWindow() {
	if AppInstance.BookmarkWindow == nil {
		AppInstance.BookmarkWindow = newBookmarkWindow()
	}
	AppInstance.BookmarkWindow.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
		AppInstance.BookmarkWindow = nil
	})
	AppInstance.BookmarkWindow.Show()
	AppInstance.BookmarkWindow.Focus()
}

func (cs *BookmarkService) CloseWindow() {
	if AppInstance.BookmarkWindow != nil {
		AppInstance.BookmarkWindow.Close()
		AppInstance.BookmarkWindow = nil
	}
}

// SetUserPassword 设置用户密码用于加密/解密私钥密码
func (bs *BookmarkService) SetUserPassword(password string) {
	bs.userPassword = password

	// 如果有等待密码输入的channel，向channel发送完成信号
	if bs.passwordChan != nil {
		select {
		case bs.passwordChan <- struct{}{}:
		default:
			// 如果channel已满或关闭，忽略
		}
		bs.passwordChan = nil // 清空channel引用
	}
}

// waitForPassword 触发事件并等待用户输入密码，超时60秒
func (bs *BookmarkService) waitForPassword(reason string) error {
	// 创建新的channel用于等待密码输入完成信号
	bs.passwordChan = make(chan struct{}, 1)

	// 触发密码输入事件
	app.Event.Emit(EventInputPasswort, reason)

	// 等待密码输入完成，超时60秒
	select {
	case <-bs.passwordChan:
		if bs.userPassword == "" {
			return errors.New("未输入密码")
		}
		app.Event.Emit(EventInputPasswortClose, "")
		return nil
	case <-time.After(60 * time.Second):
		// 超时，关闭channel
		close(bs.passwordChan)
		bs.passwordChan = nil
		return errors.New("密码输入超时")
	}
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

// ConnectBookmark 连接书签，通过书签ID触发连接事件
func (bs *BookmarkService) ConnectBookmark(bookmarkID string) {
	app.Event.Emit(EventConnectBookmark, bookmarkID)
}

// encryptField 加密单个字段
func (bs *BookmarkService) encryptField(value, fieldName string) (string, error) {
	if bs.userPassword == "" {
		if err := bs.waitForPassword("需要密码来加密" + fieldName); err != nil {
			return "", err
		}
	}
	encrypted, err := secret.Encrypt(bs.userPassword, value)
	if err != nil {
		return "", fmt.Errorf("加密%s失败: %v", fieldName, err)
	}
	return encrypted, nil
}

// encryptBookmark 对书签中的敏感字段进行加密
func (bs *BookmarkService) encryptBookmark(bookmark SSHBookmark) (SSHBookmark, error) {
	if bookmark.PrivateKeyPassword != "" {
		encrypted, err := bs.encryptField(bookmark.PrivateKeyPassword, "私钥密码")
		if err != nil {
			return bookmark, err
		}
		bookmark.PrivateKeyPassword = encrypted
	}

	if bookmark.Password != "" {
		encrypted, err := bs.encryptField(bookmark.Password, "登录密码")
		if err != nil {
			return bookmark, err
		}
		bookmark.Password = encrypted
	}

	return bookmark, nil
}

// encryptFieldIfNeeded 在需要时加密字段值
func (bs *BookmarkService) encryptFieldIfNeeded(newValue, existingValue, fieldName string) (string, error) {
	if newValue == "" || newValue == existingValue {
		return existingValue, nil
	}
	if bs.userPassword == "" {
		if err := bs.waitForPassword("需要密码来加密" + fieldName); err != nil {
			return "", err
		}
	}
	encrypted, err := secret.Encrypt(bs.userPassword, newValue)
	if err != nil {
		return "", fmt.Errorf("加密%s失败: %v", fieldName, err)
	}
	return encrypted, nil
}

// processBookmarkForSave 处理书签保存逻辑，包括条件加密
func (bs *BookmarkService) processBookmarkForSave(bookmark SSHBookmark, existingBookmark *SSHBookmark) (SSHBookmark, error) {
	if existingBookmark == nil {
		return bs.encryptBookmark(bookmark)
	}

	var err error
	bookmark.PrivateKeyPassword, err = bs.encryptFieldIfNeeded(
		bookmark.PrivateKeyPassword, existingBookmark.PrivateKeyPassword, "私钥密码")
	if err != nil {
		return bookmark, err
	}

	bookmark.Password, err = bs.encryptFieldIfNeeded(
		bookmark.Password, existingBookmark.Password, "登录密码")
	if err != nil {
		return bookmark, err
	}

	return bookmark, nil
}

// decryptField 解密单个字段，如果失败则清空用户密码
func (bs *BookmarkService) decryptField(encryptedValue, fieldName string) (string, error) {
	if bs.userPassword == "" {
		if err := bs.waitForPassword("需要密码来解密" + fieldName); err != nil {
			return "", err
		}
	}
	decrypted, err := secret.Decrypt(bs.userPassword, encryptedValue)
	if err != nil {
		bs.userPassword = ""
		return "", fmt.Errorf("解密%s失败: %v", fieldName, err)
	}
	return decrypted, nil
}

// decryptBookmark 对书签中的敏感字段进行解密
func (bs *BookmarkService) decryptBookmark(bookmark SSHBookmark) (SSHBookmark, error) {
	if bookmark.PrivateKeyPassword != "" {
		decrypted, err := bs.decryptField(bookmark.PrivateKeyPassword, "私钥密码")
		if err != nil {
			return bookmark, err
		}
		bookmark.PrivateKeyPassword = decrypted
	}

	if bookmark.Password != "" {
		decrypted, err := bs.decryptField(bookmark.Password, "登录密码")
		if err != nil {
			return bookmark, err
		}
		bookmark.Password = decrypted
	}

	return bookmark, nil
}

// DecryptPassword 解密密码字符串，如果不是base64则直接返回，是则解密后返回，失败则返回原始字符串
func (bs *BookmarkService) DecryptPassword(password string) string {
	if password == "" {
		return password
	}
	// 检查是否为base64格式
	if _, err := base64.StdEncoding.DecodeString(password); err != nil {
		// 不是base64，直接返回原字符串
		return password
	}

	// 是base64，尝试解密
	if bs.userPassword == "" {
		if err := bs.waitForPassword("需要密码来解密密码"); err != nil {
			return password // 解密失败，返回原始字符串
		}
	}

	decrypted, err := secret.Decrypt(bs.userPassword, password)
	if err != nil {
		// 解密失败，清空用户密码，让用户重新输入
		bs.userPassword = ""
		return password // 解密失败，返回原始字符串
	}

	return decrypted
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

	app.Event.Emit(EventBookmarkUpdate, "bookmark update")

	return os.WriteFile(bs.bookmarkFile, data, 0644)
}

// ListBookmarks 列出所有书签
func (bs *BookmarkService) ListBookmarks() ([]*BookmarkGroup, error) {
	return bs.loadBookmarks(), nil
}

// findBookmarkByID 查找书签位置，返回是否找到及分组索引和书签索引
func (bs *BookmarkService) findBookmarkByID(groups []*BookmarkGroup, id string) (found bool, groupIdx, bookmarkIdx int) {
	for i, group := range groups {
		for j, b := range group.Bookmarks {
			if b.ID == id {
				return true, i, j
			}
		}
	}
	return false, 0, 0
}

// findGroupIndexByName 根据分组名称查找分组索引
func (bs *BookmarkService) findGroupIndexByName(groups []*BookmarkGroup, name string) (int, bool) {
	for i, g := range groups {
		if g.Name == name {
			return i, true
		}
	}
	return 0, false
}

// updateExistingBookmark 更新现有书签
func (bs *BookmarkService) updateExistingBookmark(groups []*BookmarkGroup, bookmark SSHBookmark, groupIdx, bookmarkIdx int) ([]*BookmarkGroup, error) {
	existing := groups[groupIdx].Bookmarks[bookmarkIdx]

	targetIdx, found := bs.findGroupIndexByName(groups, bookmark.GroupName)
	if !found {
		bookmark.GroupName = existing.GroupName
		targetIdx = groupIdx
	}

	processed, err := bs.processBookmarkForSave(bookmark, &existing)
	if err != nil {
		return groups, err
	}

	if targetIdx != groupIdx {
		bookmarks := groups[groupIdx].Bookmarks
		groups[groupIdx].Bookmarks = append(bookmarks[:bookmarkIdx], bookmarks[bookmarkIdx+1:]...)
		groups[targetIdx].Bookmarks = append(groups[targetIdx].Bookmarks, processed)
	} else {
		groups[groupIdx].Bookmarks[bookmarkIdx] = processed
	}

	return groups, nil
}

// addNewBookmark 添加新书签
func (bs *BookmarkService) addNewBookmark(groups []*BookmarkGroup, bookmark SSHBookmark) ([]*BookmarkGroup, error) {
	processed, err := bs.processBookmarkForSave(bookmark, nil)
	if err != nil {
		return groups, err
	}

	if idx, found := bs.findGroupIndexByName(groups, bookmark.GroupName); found {
		groups[idx].Bookmarks = append(groups[idx].Bookmarks, processed)
		return groups, nil
	}

	defaultName := "默认书签"
	if idx, found := bs.findGroupIndexByName(groups, defaultName); found {
		groups[idx].Bookmarks = append(groups[idx].Bookmarks, processed)
	}

	return groups, nil
}

// SaveBookmark 保存SSH连接信息书签，根据ID判断是新增还是更新
func (bs *BookmarkService) SaveBookmark(bookmark SSHBookmark) error {
	groups := bs.loadBookmarks()

	found, groupIdx, bookmarkIdx := bs.findBookmarkByID(groups, bookmark.ID)

	var err error
	if found {
		groups, err = bs.updateExistingBookmark(groups, bookmark, groupIdx, bookmarkIdx)
	} else {
		groups, err = bs.addNewBookmark(groups, bookmark)
	}

	if err != nil {
		return err
	}

	bs.bookmarks = groups
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
				// 解密敏感字段
				decryptedBookmark, err := bs.decryptBookmark(bookmark)
				if err != nil {
					return nil, err
				}
				return &decryptedBookmark, nil
			}
		}
	}

	return nil, fmt.Errorf("未找到ID为 '%s' 的书签", bookmarkID)
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
