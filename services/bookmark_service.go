package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/ilaziness/vexo/internal/database"
	"github.com/ilaziness/vexo/internal/secret"
	"github.com/ilaziness/vexo/internal/utils"
	"github.com/wailsapp/wails/v3/pkg/application"
	"go.uber.org/zap"
)

const (
	EventInputPassword      = "eventInputPassword"
	EventInputPasswordClose = "eventInputPasswordClose"
	EventBookmarkUpdate     = "eventBookmarkUpdate"
	EventConnectBookmark    = "eventConnectBookmark"
	BookmarkUpdateMsg       = "bookmark update"

	// PasswordMask 密码占位符，用于前端展示
	PasswordMask = "********"
)

// 用户密码，用于加密/解密书签密码
var globalUserPassword string

func init() {
	application.RegisterEvent[string](EventInputPassword)
	application.RegisterEvent[string](EventInputPasswordClose)
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

// BookmarkListItem 书签列表项（扁平结构）
type BookmarkListItem struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	GroupName string `json:"group_name"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	User      string `json:"user"`
}

// BookmarkService 书签服务结构
type BookmarkService struct {
	db           *database.Database
	sshService   *SSHService
	passwordChan chan struct{}
}

// NewBookmarkService 创建新的书签服务实例
func NewBookmarkService(db *database.Database, ssh *SSHService) *BookmarkService {
	return &BookmarkService{
		db:         db,
		sshService: ssh,
	}
}

// SetUserPassword 设置用户密码用于加密/解密私钥密码
func (bs *BookmarkService) SetUserPassword(password string) {
	Logger.Debug("set user password", zap.String("pwd", password))
	globalUserPassword = password

	// 如果有等待密码输入的channel，向channel发送完成信号
	if bs.passwordChan != nil {
		bs.passwordChan <- struct{}{}
		bs.passwordChan = nil // 清空channel引用
	}
	Logger.Debug("user password set successfully", zap.String("pwd", globalUserPassword))
}

// waitForPassword 触发事件并等待用户输入密码，超时60秒
func (bs *BookmarkService) waitForPassword(reason string) error {
	Logger.Debug("waiting for user password input", zap.String("reason", reason), zap.String("pwd", globalUserPassword))
	if globalUserPassword != "" {
		return nil
	}
	bs.passwordChan = make(chan struct{}, 1)
	app.Event.Emit(EventInputPassword, reason)

	// 等待密码输入完成，超时60秒
	select {
	case <-bs.passwordChan:
		if globalUserPassword == "" {
			return errors.New("未输入密码")
		}
		app.Event.Emit(EventInputPasswordClose, "")
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

// ConnectBookmarkByID 连接书签，通过书签ID获取连接信息并连接
func (bs *BookmarkService) ConnectBookmarkByID(bookmarkID string) (string, error) {
	bookmark, err := bs.getDecryptedBookmarkByID(bookmarkID)
	if err != nil {
		return "", err
	}

	return bs.sshService.Connect(
		bookmark.Host,
		bookmark.Port,
		bookmark.User,
		bookmark.Password,
		bookmark.PrivateKey,
		bookmark.PrivateKeyPassword,
	)
}

// encryptField 加密单个字段
func (bs *BookmarkService) encryptField(value, fieldName string) (string, error) {
	if globalUserPassword == "" {
		if err := bs.waitForPassword("需要密码来加密" + fieldName); err != nil {
			return "", err
		}
	}
	encrypted, err := secret.Encrypt(globalUserPassword, value)
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
	// 空值表示清空密码
	if newValue == "" {
		return "", nil
	}
	// 占位符表示未修改，保持原值
	if newValue == PasswordMask {
		return existingValue, nil
	}
	// 新值，需要加密
	if globalUserPassword == "" {
		if err := bs.waitForPassword("需要密码来加密" + fieldName); err != nil {
			return "", err
		}
	}
	encrypted, err := secret.Encrypt(globalUserPassword, newValue)
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
	if globalUserPassword == "" {
		if err := bs.waitForPassword("需要密码来解密" + fieldName); err != nil {
			return "", err
		}
	}
	if globalUserPassword == "" {
		Logger.Debug("user password is still empty after waiting, cannot decrypt", zap.String("field", fieldName), zap.String("pwd", globalUserPassword))
		return "", errors.New("未输入密码")
	}
	decrypted, err := secret.Decrypt(globalUserPassword, encryptedValue)
	if err != nil {
		globalUserPassword = ""
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

// ListBookmarks 列出所有书签（分组结构）
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

	// 构建分组ID到名称的映射
	groupIDToName := make(map[int]string)
	groupMap := make(map[int]*BookmarkGroup)
	for _, g := range dbGroups {
		groupIDToName[g.ID] = g.Name
		groupMap[g.ID] = &BookmarkGroup{Name: g.Name, Bookmarks: []SSHBookmark{}}
	}

	// 组装书签到对应分组（密码使用占位符，不解密）
	for _, b := range dbBookmarks {
		bookmark := SSHBookmark{
			ID:                 b.ID,
			GroupName:          groupIDToName[b.GroupID],
			Title:              b.Title,
			Host:               b.Host,
			Port:               b.Port,
			User:               b.User,
			Password:           bs.maskPassword(b.Password),
			PrivateKey:         b.PrivateKey,
			PrivateKeyPassword: bs.maskPassword(b.PrivateKeyPassword),
		}
		if group, ok := groupMap[b.GroupID]; ok {
			group.Bookmarks = append(group.Bookmarks, bookmark)
		}
	}

	// 按数据库查询顺序（id 升序，即创建顺序）转换为切片
	groups := make([]*BookmarkGroup, 0, len(dbGroups))
	for _, g := range dbGroups {
		if group, ok := groupMap[g.ID]; ok {
			groups = append(groups, group)
		}
	}

	return groups, nil
}

// GetAllBookmarks 获取全部书签列表（扁平结构）
func (bs *BookmarkService) GetAllBookmarks() ([]*BookmarkListItem, error) {
	// 从数据库加载分组
	dbGroups, err := bs.db.BookmarkRepo.GetAllGroups()
	if err != nil {
		return nil, err
	}

	// 构建分组ID到名称的映射
	groupIDToName := make(map[int]string)
	for _, g := range dbGroups {
		groupIDToName[g.ID] = g.Name
	}

	// 从数据库加载书签
	dbBookmarks, err := bs.db.BookmarkRepo.GetAllBookmarks()
	if err != nil {
		return nil, err
	}

	// 转换为列表项
	items := make([]*BookmarkListItem, 0, len(dbBookmarks))
	for _, b := range dbBookmarks {
		items = append(items, &BookmarkListItem{
			ID:        b.ID,
			Title:     b.Title,
			GroupName: groupIDToName[b.GroupID],
			Host:      b.Host,
			Port:      b.Port,
			User:      b.User,
		})
	}

	return items, nil
}

// maskPassword 密码掩码处理：如果有密码返回占位符，空则返回空
func (bs *BookmarkService) maskPassword(password string) string {
	if password == "" {
		return ""
	}
	return PasswordMask
}

// GetBookmarkByID 通过书签 ID 查找返回连接信息（密码使用占位符，不解密）
func (bs *BookmarkService) GetBookmarkByID(bookmarkID string) (*SSHBookmark, error) {
	bookmark, err := bs.getBookmarkByID(bookmarkID)
	if err != nil {
		return nil, err
	}
	bookmark.Password = bs.maskPassword(bookmark.Password)
	bookmark.PrivateKeyPassword = bs.maskPassword(bookmark.PrivateKeyPassword)
	return bookmark, nil
}

// getDecryptedBookmarkByID 内部方法：获取解密后的书签（用于连接和测试）
func (bs *BookmarkService) getDecryptedBookmarkByID(bookmarkID string) (*SSHBookmark, error) {
	bookmark, err := bs.getBookmarkByID(bookmarkID)
	if err != nil {
		return nil, err
	}
	decrypted, err := bs.decryptBookmark(*bookmark)
	if err != nil {
		return nil, err
	}
	return &decrypted, nil
}

// getBookmarkByID 内部方法：获取书签基础信息（未解密）
func (bs *BookmarkService) getBookmarkByID(bookmarkID string) (*SSHBookmark, error) {
	dbBookmark, err := bs.db.BookmarkRepo.GetBookmarkByID(bookmarkID)
	if err != nil {
		return nil, fmt.Errorf("未找到 ID 为 '%s' 的书签", bookmarkID)
	}

	group, _ := bs.db.BookmarkRepo.GetGroupByID(dbBookmark.GroupID)
	groupName := ""
	if group != nil {
		groupName = group.Name
	}

	return &SSHBookmark{
		ID:                 dbBookmark.ID,
		GroupName:          groupName,
		Title:              dbBookmark.Title,
		Host:               dbBookmark.Host,
		Port:               dbBookmark.Port,
		User:               dbBookmark.User,
		Password:           dbBookmark.Password,
		PrivateKey:         dbBookmark.PrivateKey,
		PrivateKeyPassword: dbBookmark.PrivateKeyPassword,
	}, nil
}

// SaveBookmark 保存 SSH 连接信息书签，根据 ID 判断是新增还是更新
func (bs *BookmarkService) SaveBookmark(bookmark SSHBookmark) error {
	Logger.Debug("savebookmark", zap.Any("bk", bookmark))
	if bookmark.ID != "" {
		existing, err := bs.db.BookmarkRepo.GetBookmarkByID(bookmark.ID)
		if err == nil && existing != nil {
			return bs.updateBookmark(bookmark, existing)
		}
	}
	return bs.insertBookmark(bookmark)
}

// updateBookmark 更新已有书签
func (bs *BookmarkService) updateBookmark(bookmark SSHBookmark, existing *database.BookmarkDB) error {
	existingBookmark := &SSHBookmark{
		ID:                 existing.ID,
		Password:           existing.Password,
		PrivateKeyPassword: existing.PrivateKeyPassword,
	}

	// 处理加密逻辑
	processed, err := bs.processBookmarkForSave(bookmark, existingBookmark)
	if err != nil {
		return err
	}

	// 获取分组ID
	groupID := existing.GroupID
	if bookmark.GroupName != "" {
		if group, err := bs.db.BookmarkRepo.GetGroupByName(bookmark.GroupName); err == nil {
			groupID = group.ID
		}
	}

	dbBookmark := &database.BookmarkDB{
		ID:                 processed.ID,
		GroupID:            groupID,
		Title:              processed.Title,
		Host:               processed.Host,
		Port:               processed.Port,
		User:               processed.User,
		Password:           processed.Password,
		PrivateKey:         processed.PrivateKey,
		PrivateKeyPassword: processed.PrivateKeyPassword,
		UpdatedAt:          time.Now(),
	}

	if err := bs.db.BookmarkRepo.UpdateBookmark(dbBookmark); err != nil {
		return err
	}

	app.Event.Emit(EventBookmarkUpdate, BookmarkUpdateMsg)
	Logger.Debug("bookmark updated")
	return nil
}

// insertBookmark 插入新书签
func (bs *BookmarkService) insertBookmark(bookmark SSHBookmark) error {
	group, err := bs.db.BookmarkRepo.GetGroupByName(bookmark.GroupName)
	if err != nil {
		group, _ = bs.db.BookmarkRepo.GetGroupByName("默认书签")
	}

	processed, err := bs.encryptBookmark(bookmark)
	if err != nil {
		return err
	}

	dbBookmark := &database.BookmarkDB{
		ID:                 utils.GenerateRandomID(),
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

	if err := bs.db.BookmarkRepo.InsertBookmark(dbBookmark); err != nil {
		return err
	}

	app.Event.Emit(EventBookmarkUpdate, BookmarkUpdateMsg)
	Logger.Debug("bookmark inserted")
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

// UpdateGroup 更新分组名称（同步更新该分组下所有书签）
func (bs *BookmarkService) UpdateGroup(oldName, newName string) error {
	// 不能修改默认分组名称
	if oldName == "默认书签" {
		return fmt.Errorf("不能修改默认分组名称")
	}

	// 检查新名称是否已存在
	_, err := bs.db.BookmarkRepo.GetGroupByName(newName)
	if err == nil {
		return fmt.Errorf("分组 '%s' 已存在", newName)
	}

	// 更新分组名称
	err = bs.db.BookmarkRepo.UpdateGroupName(oldName, newName)
	if err != nil {
		return err
	}

	app.Event.Emit(EventBookmarkUpdate, BookmarkUpdateMsg)
	Logger.Debug("group name updated", zap.String("old", oldName), zap.String("new", newName))
	return nil
}

// DeleteGroup 删除分组（分组下必须没有书签）
func (bs *BookmarkService) DeleteGroup(groupName string) error {
	// 不能删除默认分组
	if groupName == "默认书签" {
		return fmt.Errorf("不能删除默认分组")
	}

	// 检查分组下是否有书签
	count, err := bs.db.BookmarkRepo.GetGroupBookmarkCount(groupName)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("分组 '%s' 下还有 %d 个书签，无法删除", groupName, count)
	}

	err = bs.db.BookmarkRepo.DeleteGroup(groupName)
	if err != nil {
		return err
	}

	app.Event.Emit(EventBookmarkUpdate, BookmarkUpdateMsg)
	Logger.Debug("group deleted", zap.String("name", groupName))
	return nil
}

// TestConnection 测试连接
func (bs *BookmarkService) TestConnection(bookmark SSHBookmark) error {
	testData := bookmark

	// 有ID时检查是否有修改
	if bookmark.ID != "" {
		existing, err := bs.getBookmarkByID(bookmark.ID)
		if err != nil {
			return err
		}

		// 关键字段无变化（密码为""或占位符表示未修改），使用数据库解密数据
		if bookmark.Host == existing.Host &&
			bookmark.Port == existing.Port &&
			bookmark.User == existing.User &&
			bookmark.PrivateKey == existing.PrivateKey &&
			(bookmark.Password == "" || bookmark.Password == PasswordMask) &&
			(bookmark.PrivateKeyPassword == "" || bookmark.PrivateKeyPassword == PasswordMask) {
			decrypted, err := bs.getDecryptedBookmarkByID(bookmark.ID)
			if err != nil {
				return err
			}
			testData = *decrypted
		}
	}

	return bs.sshService.TestConnectInfo(
		testData.Host,
		testData.Port,
		testData.User,
		testData.Password,
		testData.PrivateKey,
		testData.PrivateKeyPassword,
	)
}

// SaveAndConnect 保存书签并连接（先测试，成功后保存，然后连接）
func (bs *BookmarkService) SaveAndConnect(bookmark SSHBookmark) error {
	if err := bs.TestConnection(bookmark); err != nil {
		return fmt.Errorf("连接测试失败: %w", err)
	}

	// 保存书签
	if err := bs.SaveBookmark(bookmark); err != nil {
		return fmt.Errorf("保存书签失败: %w", err)
	}

	// 触发连接事件
	bs.ConnectBookmark(bookmark.ID)
	return nil
}

// GetBookmarkForConnect 获取用于连接的书签信息
func (bs *BookmarkService) GetBookmarkForConnect(bookmarkID string) (*SSHBookmark, error) {
	bookmark, err := bs.getBookmarkByID(bookmarkID)
	if err != nil {
		return nil, err
	}
	return bookmark, nil
}
