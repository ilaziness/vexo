package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

const (
	EventInputPasswort = "eventInput_Password"
)

func init() {
	application.RegisterEvent[string](EventInputPasswort)
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
	userPassword []byte
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

// SetUserPassword 设置用户密码用于加密/解密私钥密码
func (bs *BookmarkService) SetUserPassword(password string) {
	// 使用SHA-256哈希密码作为AES密钥，提供更好的安全性
	hash := sha256.Sum256([]byte(password))
	bs.userPassword = hash[:]

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
		// 密码已经设置完成
		return nil
	case <-time.After(60 * time.Second):
		// 超时，关闭channel
		close(bs.passwordChan)
		bs.passwordChan = nil
		return fmt.Errorf("密码输入超时")
	}
}

// encrypt 使用用户密码加密数据
func (bs *BookmarkService) encrypt(plaintext string) (string, error) {
	if len(bs.userPassword) == 0 {
		// 如果没有设置密码，等待用户输入密码
		if err := bs.waitForPassword("需要密码来加密私钥密码"); err != nil {
			return "", err
		}
	}

	block, err := aes.NewCipher(bs.userPassword)
	if err != nil {
		return "", err
	}

	plaintextBytes := []byte(plaintext)
	// CTR模式使用nonce而不是IV，nonce长度等于block size
	ciphertext := make([]byte, aes.BlockSize+len(plaintextBytes))
	nonce := ciphertext[:aes.BlockSize]
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	stream := cipher.NewCTR(block, nonce)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintextBytes)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decrypt 使用用户密码解密数据
func (bs *BookmarkService) decrypt(encryptedText string) (string, error) {
	if len(bs.userPassword) == 0 {
		// 如果没有设置密码，等待用户输入密码
		if err := bs.waitForPassword("需要密码来解密私钥密码"); err != nil {
			return "", err
		}
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(bs.userPassword)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCTR(block, nonce)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
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
		return
	}
	bs.window = app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "书签",
		URL:              "/bookmark",
		Width:            1200,
		Height:           800,
		BackgroundColour: application.NewRGB(27, 38, 54),
	})
	bs.window.OnWindowEvent(events.Common.WindowClosing, func(_ *application.WindowEvent) {
		bs.window = nil
	})
	bs.window.Show()
}

// CloseWindow 关闭窗口
func (bs *BookmarkService) CloseWindow() {
	if bs.window != nil {
		bs.window.Close()
		bs.window = nil
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

// SaveBookmark 保存SSH连接信息书签，根据ID判断是新增还是更新
func (bs *BookmarkService) SaveBookmark(bookmark SSHBookmark) error {
	bookmarkGroups := bs.loadBookmarks()

	// 检查ID是否已存在
	bookmarkExists := false
	var existingGroupIndex int
	var existingBookmarkIndex int

	for i, group := range bookmarkGroups {
		for j, existingBookmark := range group.Bookmarks {
			if existingBookmark.ID == bookmark.ID {
				bookmarkExists = true
				existingGroupIndex = i
				existingBookmarkIndex = j
				break
			}
		}
		if bookmarkExists {
			break
		}
	}

	if bookmarkExists {
		// 更新现有书签
		existingBookmark := bookmarkGroups[existingGroupIndex].Bookmarks[existingBookmarkIndex]

		// 检查目标分组是否存在
		targetGroupExists := false
		var targetGroupIndex int
		for i, g := range bookmarkGroups {
			if g.Name == bookmark.GroupName {
				targetGroupExists = true
				targetGroupIndex = i
				break
			}
		}

		if !targetGroupExists {
			// 如果目标分组不存在，使用原分组
			bookmark.GroupName = existingBookmark.GroupName
			targetGroupIndex = existingGroupIndex
		}

		// 处理私钥密码更新逻辑
		if bookmark.PrivateKeyPassword == "" {
			// 如果前端没有传递私钥密码，保持原有的加密密码
			bookmark.PrivateKeyPassword = existingBookmark.PrivateKeyPassword
		} else {
			// 如果前端传递了私钥密码，加密保存
			encryptedPassword, err := bs.encrypt(bookmark.PrivateKeyPassword)
			if err != nil {
				return fmt.Errorf("加密私钥密码失败: %v", err)
			}
			bookmark.PrivateKeyPassword = encryptedPassword
		}

		// 处理登录密码更新逻辑
		if bookmark.Password == "" {
			// 如果前端没有传递登录密码，保持原有的加密密码
			bookmark.Password = existingBookmark.Password
		} else {
			// 如果前端传递了登录密码，加密保存
			encryptedPassword, err := bs.encrypt(bookmark.Password)
			if err != nil {
				return fmt.Errorf("加密登录密码失败: %v", err)
			}
			bookmark.Password = encryptedPassword
		}

		// 如果分组发生变化，需要移动书签
		if targetGroupIndex != existingGroupIndex {
			// 从原分组移除
			bookmarks := bookmarkGroups[existingGroupIndex].Bookmarks
			bookmarkGroups[existingGroupIndex].Bookmarks = append(bookmarks[:existingBookmarkIndex], bookmarks[existingBookmarkIndex+1:]...)

			// 添加到新分组
			bookmarkGroups[targetGroupIndex].Bookmarks = append(bookmarkGroups[targetGroupIndex].Bookmarks, bookmark)
		} else {
			// 更新书签信息
			bookmarkGroups[existingGroupIndex].Bookmarks[existingBookmarkIndex] = bookmark
		}
	} else {
		// 新增书签
		// 如果有私钥密码，加密保存
		if bookmark.PrivateKeyPassword != "" {
			encryptedPassword, err := bs.encrypt(bookmark.PrivateKeyPassword)
			if err != nil {
				return fmt.Errorf("加密私钥密码失败: %v", err)
			}
			bookmark.PrivateKeyPassword = encryptedPassword
		}

		// 如果有登录密码，加密保存
		if bookmark.Password != "" {
			encryptedPassword, err := bs.encrypt(bookmark.Password)
			if err != nil {
				return fmt.Errorf("加密登录密码失败: %v", err)
			}
			bookmark.Password = encryptedPassword
		}

		// 检查分组是否存在，如果不存在则使用默认分组
		groupExists := false
		for i, group := range bookmarkGroups {
			if group.Name == bookmark.GroupName {
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
					bookmarkGroups[i].Bookmarks = append(group.Bookmarks, bookmark)
					break
				}
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
				// 如果有加密的私钥密码，解密返回
				if bookmark.PrivateKeyPassword != "" {
					decryptedPassword, err := bs.decrypt(bookmark.PrivateKeyPassword)
					if err != nil {
						return nil, fmt.Errorf("解密私钥密码失败: %v", err)
					}
					bookmark.PrivateKeyPassword = decryptedPassword
				}

				// 如果有加密的登录密码，解密返回
				if bookmark.Password != "" {
					decryptedPassword, err := bs.decrypt(bookmark.Password)
					if err != nil {
						return nil, fmt.Errorf("解密登录密码失败: %v", err)
					}
					bookmark.Password = decryptedPassword
				}

				return &bookmark, nil
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
