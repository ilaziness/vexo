package services

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/ilaziness/vexo/internal/database"
	"github.com/ilaziness/vexo/internal/secret"
	"github.com/wailsapp/wails/v3/pkg/application"
	"go.uber.org/zap"
)

const (
	EventInputPasswort      = "eventInputPassword"
	EventInputPasswortClose = "eventInputPasswordClose"
	EventBookmarkUpdate     = "eventBookmarkUpdate"
	EventConnectBookmark    = "eventConnectBookmark"
	BookmarkUpdateMsg       = "bookmark update"
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
	db           *database.Database
	window       *application.WebviewWindow
	userPassword string
	passwordChan chan struct{}
}

// NewBookmarkService 创建新的书签服务实例
func NewBookmarkService(cfgs *ConfigService, db *database.Database) *BookmarkService {
	return &BookmarkService{
		cfgService: cfgs,
		db:         db,
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
	bs.passwordChan = make(chan struct{}, 1)
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

// ListBookmarks 列出所有书签
func (bs *BookmarkService) ListBookmarks() ([]*BookmarkGroup, error) {
	// 从数据库加载分组
	dbGroups, err := bs.db.BookmarkRepo.GetAllGroups()
	if err != nil {
		return nil, err
	}

	// 从数据库加载书签
	dbBookmarks, err := bs.db.BookmarkRepo.GetAllBookmarks()
	if err != nil {
		return nil, err
	}

	// 临时组装为分组结构返回
	groupMap := make(map[int]*BookmarkGroup)
	for _, g := range dbGroups {
		groupMap[g.ID] = &BookmarkGroup{Name: g.Name, Bookmarks: []SSHBookmark{}}
	}

	for _, b := range dbBookmarks {
		bookmark := SSHBookmark{
			ID:                 b.ID,
			GroupName:          "",
			Title:              b.Title,
			Host:               b.Host,
			Port:               b.Port,
			User:               b.User,
			Password:           b.Password,
			PrivateKey:         b.PrivateKey,
			PrivateKeyPassword: b.PrivateKeyPassword,
		}
		if group, ok := groupMap[b.GroupID]; ok {
			bookmark.GroupName = group.Name
			group.Bookmarks = append(group.Bookmarks, bookmark)
		}
	}

	groups := make([]*BookmarkGroup, 0, len(groupMap))
	for _, group := range groupMap {
		groups = append(groups, group)
	}

	return groups, nil
}

// SaveBookmark 保存 SSH 连接信息书签，根据 ID 判断是新增还是更新
func (bs *BookmarkService) SaveBookmark(bookmark SSHBookmark) error {
	// 查询是否已存在
	existing, err := bs.db.BookmarkRepo.GetBookmarkByID(bookmark.ID)

	if err == nil && existing != nil {
		// 更新：先解密原有数据进行对比
		decryptedExisting, _ := bs.decryptBookmark(SSHBookmark{
			ID:                 existing.ID,
			Password:           existing.Password,
			PrivateKeyPassword: existing.PrivateKeyPassword,
		})

		// 处理加密逻辑
		processed, err := bs.processBookmarkForSave(bookmark, &decryptedExisting)
		if err != nil {
			return err
		}

		// 更新数据库
		dbBookmark := &database.BookmarkDB{
			ID:                 processed.ID,
			GroupID:            existing.GroupID,
			Title:              processed.Title,
			Host:               processed.Host,
			Port:               processed.Port,
			User:               processed.User,
			Password:           processed.Password,
			PrivateKey:         processed.PrivateKey,
			PrivateKeyPassword: processed.PrivateKeyPassword,
			UpdatedAt:          time.Now(),
		}
		err = bs.db.BookmarkRepo.UpdateBookmark(dbBookmark)
		if err != nil {
			return err
		}
	} else {
		// 新增：需要获取分组 ID
		group, err := bs.db.BookmarkRepo.GetGroupByName(bookmark.GroupName)
		if err != nil {
			// 如果分组不存在，使用默认分组
			group, _ = bs.db.BookmarkRepo.GetGroupByName("默认书签")
		}

		// 处理加密逻辑
		processed, err := bs.encryptBookmark(bookmark)
		if err != nil {
			return err
		}

		// 插入数据库
		dbBookmark := &database.BookmarkDB{
			ID:                 processed.ID,
			GroupID:            group.ID,
			Title:              processed.Title,
			Host:               processed.Host,
			Port:               processed.Port,
			User:               processed.User,
			Password:           processed.Password,
			PrivateKey:         processed.PrivateKey,
			PrivateKeyPassword: processed.PrivateKeyPassword,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}
		err = bs.db.BookmarkRepo.InsertBookmark(dbBookmark)
		if err != nil {
			return err
		}
	}

	app.Event.Emit(EventBookmarkUpdate, BookmarkUpdateMsg)
	Logger.Debug("bookmark saved")
	return nil
}

// DeleteGroup 删除分组
func (bs *BookmarkService) DeleteGroup(groupName string) error {
	// 不能删除默认分组
	if groupName == "默认书签" {
		return fmt.Errorf("不能删除默认分组")
	}

	err := bs.db.BookmarkRepo.DeleteGroup(groupName)
	if err != nil {
		return err
	}

	app.Event.Emit(EventBookmarkUpdate, BookmarkUpdateMsg)
	Logger.Debug("group deleted", zap.String("name", groupName))
	return nil
}

// DeleteBookmark 删除书签
func (bs *BookmarkService) DeleteBookmark(bookmarkID string) error {
	err := bs.db.BookmarkRepo.DeleteBookmark(bookmarkID)
	if err != nil {
		return err
	}

	app.Event.Emit(EventBookmarkUpdate, BookmarkUpdateMsg)
	Logger.Debug("bookmark deleted", zap.String("id", bookmarkID))
	return nil
}

// GetBookmarkByID 通过书签 ID 查找返回连接信息
func (bs *BookmarkService) GetBookmarkByID(bookmarkID string) (*SSHBookmark, error) {
	dbBookmark, err := bs.db.BookmarkRepo.GetBookmarkByID(bookmarkID)
	if err != nil {
		return nil, fmt.Errorf("未找到 ID 为 '%s' 的书签", bookmarkID)
	}

	// 解密敏感字段
	bookmark := &SSHBookmark{
		ID:                 dbBookmark.ID,
		GroupName:          "",
		Title:              dbBookmark.Title,
		Host:               dbBookmark.Host,
		Port:               dbBookmark.Port,
		User:               dbBookmark.User,
		Password:           dbBookmark.Password,
		PrivateKey:         dbBookmark.PrivateKey,
		PrivateKeyPassword: dbBookmark.PrivateKeyPassword,
	}

	decrypted, err := bs.decryptBookmark(*bookmark)
	if err != nil {
		return nil, err
	}

	return &decrypted, nil
}

// AddGroup 新增分组
func (bs *BookmarkService) AddGroup(groupName string) error {
	// 检查是否已存在
	_, err := bs.db.BookmarkRepo.GetGroupByName(groupName)
	if err == nil {
		return fmt.Errorf("分组 '%s' 已存在", groupName)
	}

	err = bs.db.BookmarkRepo.InsertGroup(&database.BookmarkGroupDB{
		Name: groupName,
	})
	if err != nil {
		return err
	}

	app.Event.Emit(EventBookmarkUpdate, BookmarkUpdateMsg)
	Logger.Debug("group added", zap.String("name", groupName))
	return nil
}

// UpdateGroupName 更新分组名称
func (bs *BookmarkService) UpdateGroupName(oldName, newName string) error {
	// 不能修改默认分组名称
	if oldName == "默认书签" {
		return fmt.Errorf("不能修改默认分组名称")
	}

	// 检查新名称是否已存在
	_, err := bs.db.BookmarkRepo.GetGroupByName(newName)
	if err == nil {
		return fmt.Errorf("分组 '%s' 已存在", newName)
	}

	err = bs.db.BookmarkRepo.UpdateGroupName(oldName, newName)
	if err != nil {
		return err
	}

	app.Event.Emit(EventBookmarkUpdate, BookmarkUpdateMsg)
	Logger.Debug("group name updated", zap.String("old", oldName), zap.String("new", newName))
	return nil
}
