package bookmark

import (
	"errors"
	"fmt"
	"time"

	"github.com/ilaziness/vexo/internal/database"
	"github.com/ilaziness/vexo/internal/secret"
	"github.com/ilaziness/vexo/internal/ssh"
	"github.com/ilaziness/vexo/internal/utils"
	"go.uber.org/zap"
)

const PasswordMask = "********"

type PasswordFn func(reason string) (string, error)

type Bookmark struct {
	GroupName          string `json:"group_name"`
	ID                 string `json:"id"`
	Title              string `json:"title"`
	Host               string `json:"host"`
	Port               int    `json:"port"`
	PrivateKey         string `json:"private_key"`
	PrivateKeyPassword string `json:"private_key_password"`
	ProxyJumpID        string `json:"proxy_jump_id"`
	User               string `json:"user"`
	Password           string `json:"password"`
}

func (b Bookmark) Endpoint() ssh.Endpoint {
	return ssh.Endpoint{
		Host: b.Host, Port: b.Port, User: b.User,
		Password: b.Password, Key: b.PrivateKey, KeyPassword: b.PrivateKeyPassword,
	}
}

type Group struct {
	Name      string     `json:"name"`
	Bookmarks []Bookmark `json:"bookmarks"`
}

type ListItem struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	GroupName string `json:"group_name"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	User      string `json:"user"`
}

type Service struct {
	logger     *zap.Logger
	db         *database.Database
	passwordFn PasswordFn
	onUpdate   func()
	onBadPass  func()
}

func New(logger *zap.Logger, db *database.Database, passwordFn PasswordFn, onUpdate, onBadPass func()) *Service {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &Service{logger: logger, db: db, passwordFn: passwordFn, onUpdate: onUpdate, onBadPass: onBadPass}
}

func (s *Service) emit() {
	if s.onUpdate != nil {
		s.onUpdate()
	}
}

func mask(p string) string {
	if p == "" {
		return ""
	}
	return PasswordMask
}

func (s *Service) password(reason string) (string, error) {
	if s.passwordFn == nil {
		return "", errors.New("password not entered")
	}
	return s.passwordFn(reason)
}

func (s *Service) encryptField(value, fieldName string) (string, error) {
	password, err := s.password("需要密码来加密书签相关密码")
	if err != nil {
		return "", err
	}
	encrypted, err := secret.Encrypt(password, value)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt %s: %w", fieldName, err)
	}
	return encrypted, nil
}

func (s *Service) encryptBookmark(b Bookmark) (Bookmark, error) {
	if b.PrivateKeyPassword != "" {
		v, err := s.encryptField(b.PrivateKeyPassword, "private key password")
		if err != nil {
			return b, err
		}
		b.PrivateKeyPassword = v
	}
	if b.Password != "" {
		v, err := s.encryptField(b.Password, "login password")
		if err != nil {
			return b, err
		}
		b.Password = v
	}
	return b, nil
}

func (s *Service) encryptFieldIfNeeded(newValue, existing, fieldName string) (string, error) {
	if newValue == "" {
		return "", nil
	}
	if newValue == PasswordMask {
		return existing, nil
	}
	return s.encryptField(newValue, fieldName)
}

func (s *Service) decryptField(encrypted, fieldName string) (string, error) {
	password, err := s.password("需要密码来加密书签相关密码")
	if err != nil {
		return "", err
	}
	if password == "" {
		return "", errors.New("password not entered")
	}
	decrypted, err := secret.Decrypt(password, encrypted)
	if err != nil {
		if s.onBadPass != nil {
			s.onBadPass()
		}
		return "", fmt.Errorf("failed to decrypt %s: %w", fieldName, err)
	}
	return decrypted, nil
}

func (s *Service) decryptBookmark(b Bookmark) (Bookmark, error) {
	if b.PrivateKeyPassword != "" {
		v, err := s.decryptField(b.PrivateKeyPassword, "private key password")
		if err != nil {
			return b, err
		}
		b.PrivateKeyPassword = v
	}
	if b.Password != "" {
		v, err := s.decryptField(b.Password, "login password")
		if err != nil {
			return b, err
		}
		b.Password = v
	}
	return b, nil
}

func (s *Service) Get(id string) (*Bookmark, error) {
	dbBookmark, err := s.db.BookmarkRepo.GetBookmarkByID(id)
	if err != nil {
		return nil, fmt.Errorf("未找到 ID 为 '%s' 的书签", id)
	}
	groupName := ""
	if group, _ := s.db.BookmarkRepo.GetGroupByID(dbBookmark.GroupID); group != nil {
		groupName = group.Name
	}
	return fromDB(dbBookmark, groupName), nil
}

func fromDB(b *database.BookmarkDB, groupName string) *Bookmark {
	return &Bookmark{
		ID: b.ID, GroupName: groupName, Title: b.Title, Host: b.Host, Port: b.Port,
		User: b.User, Password: b.Password, PrivateKey: b.PrivateKey,
		PrivateKeyPassword: b.PrivateKeyPassword, ProxyJumpID: b.ProxyJumpID,
	}
}

func (s *Service) GetMasked(id string) (*Bookmark, error) {
	b, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	b.Password = mask(b.Password)
	b.PrivateKeyPassword = mask(b.PrivateKeyPassword)
	return b, nil
}

func (s *Service) GetDecrypted(id string) (*Bookmark, error) {
	b, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	d, err := s.decryptBookmark(*b)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (s *Service) ResolveHops(target ssh.Endpoint, jumpID string) ([]ssh.Endpoint, error) {
	hops := []ssh.Endpoint{target}
	seen := map[string]bool{target.Addr(): true}
	id := jumpID
	for id != "" {
		if len(hops) > 6 {
			return nil, fmt.Errorf("跳板机层数超过最大限制 (5)")
		}
		b, err := s.GetDecrypted(id)
		if err != nil {
			return nil, fmt.Errorf("failed to load proxy jump bookmark: %v", err)
		}
		ep := b.Endpoint()
		if seen[ep.Addr()] {
			return nil, fmt.Errorf("检测到跳板机循环引用: %s", ep.Addr())
		}
		if ep.Host == hops[len(hops)-1].Host && ep.Port == hops[len(hops)-1].Port {
			return nil, fmt.Errorf("跳板机不能与目标主机相同")
		}
		seen[ep.Addr()] = true
		hops = append([]ssh.Endpoint{ep}, hops...)
		id = b.ProxyJumpID
	}
	return hops, nil
}

func (s *Service) ResolveHopsForBookmark(id string) ([]ssh.Endpoint, error) {
	b, err := s.GetDecrypted(id)
	if err != nil {
		return nil, err
	}
	return s.ResolveHops(b.Endpoint(), b.ProxyJumpID)
}

func (s *Service) ListGrouped() ([]*Group, error) {
	dbGroups, err := s.db.BookmarkRepo.GetAllGroups()
	if err != nil {
		return nil, err
	}
	dbBookmarks, err := s.db.BookmarkRepo.GetAllBookmarks()
	if err != nil {
		return nil, err
	}
	groupIDToName := make(map[int]string)
	groupMap := make(map[int]*Group)
	for _, g := range dbGroups {
		groupIDToName[g.ID] = g.Name
		groupMap[g.ID] = &Group{Name: g.Name, Bookmarks: []Bookmark{}}
	}
	for _, b := range dbBookmarks {
		item := Bookmark{
			ID: b.ID, GroupName: groupIDToName[b.GroupID], Title: b.Title, Host: b.Host, Port: b.Port,
			User: b.User, Password: mask(b.Password), PrivateKey: b.PrivateKey,
			PrivateKeyPassword: mask(b.PrivateKeyPassword), ProxyJumpID: b.ProxyJumpID,
		}
		if g, ok := groupMap[b.GroupID]; ok {
			g.Bookmarks = append(g.Bookmarks, item)
		}
	}
	groups := make([]*Group, 0, len(dbGroups))
	for _, g := range dbGroups {
		if group, ok := groupMap[g.ID]; ok {
			groups = append(groups, group)
		}
	}
	return groups, nil
}

func (s *Service) ListItems() ([]*ListItem, error) {
	dbGroups, err := s.db.BookmarkRepo.GetAllGroups()
	if err != nil {
		return nil, err
	}
	names := map[int]string{}
	for _, g := range dbGroups {
		names[g.ID] = g.Name
	}
	dbBookmarks, err := s.db.BookmarkRepo.GetAllBookmarks()
	if err != nil {
		return nil, err
	}
	items := make([]*ListItem, 0, len(dbBookmarks))
	for _, b := range dbBookmarks {
		items = append(items, &ListItem{ID: b.ID, Title: b.Title, GroupName: names[b.GroupID], Host: b.Host, Port: b.Port, User: b.User})
	}
	return items, nil
}

func (s *Service) Save(b Bookmark) (string, error) {
	if b.ID != "" {
		existing, err := s.db.BookmarkRepo.GetBookmarkByID(b.ID)
		if err == nil && existing != nil {
			return b.ID, s.update(b, existing)
		}
	}
	return s.insert(b)
}

func (s *Service) update(b Bookmark, existing *database.BookmarkDB) error {
	processed, err := s.encryptBookmarkForSave(b, &Bookmark{Password: existing.Password, PrivateKeyPassword: existing.PrivateKeyPassword})
	if err != nil {
		return err
	}
	groupID := existing.GroupID
	if b.GroupName != "" {
		if g, err := s.db.BookmarkRepo.GetGroupByName(b.GroupName); err == nil {
			groupID = g.ID
		}
	}
	if existing.Title != b.Title || existing.GroupID != groupID {
		dup, err := s.db.BookmarkRepo.GetBookmarkByTitleAndGroup(b.Title, groupID)
		if err == nil && dup.ID != b.ID {
			return fmt.Errorf("分组 '%s' 中已存在名称为 '%s' 的书签", b.GroupName, b.Title)
		}
	}
	dbBookmark := &database.BookmarkDB{
		ID: processed.ID, GroupID: groupID, Title: processed.Title, Host: processed.Host, Port: processed.Port,
		User: processed.User, Password: processed.Password, PrivateKey: processed.PrivateKey,
		PrivateKeyPassword: processed.PrivateKeyPassword, ProxyJumpID: processed.ProxyJumpID, UpdatedAt: time.Now(),
	}
	if err := s.db.BookmarkRepo.UpdateBookmark(dbBookmark); err != nil {
		return err
	}
	s.emit()
	return nil
}

func (s *Service) encryptBookmarkForSave(b Bookmark, existing *Bookmark) (Bookmark, error) {
	if existing == nil {
		return s.encryptBookmark(b)
	}
	var err error
	b.PrivateKeyPassword, err = s.encryptFieldIfNeeded(b.PrivateKeyPassword, existing.PrivateKeyPassword, "private key password")
	if err != nil {
		return b, err
	}
	b.Password, err = s.encryptFieldIfNeeded(b.Password, existing.Password, "login password")
	return b, err
}

func (s *Service) insert(b Bookmark) (string, error) {
	group, err := s.db.BookmarkRepo.GetGroupByName(b.GroupName)
	if err != nil {
		group, err = s.db.BookmarkRepo.GetGroupByName("默认书签")
		if err != nil || group == nil {
			return "", fmt.Errorf("未找到分组 '%s'", b.GroupName)
		}
	}
	_, err = s.db.BookmarkRepo.GetBookmarkByTitleAndGroup(b.Title, group.ID)
	if err == nil {
		return "", fmt.Errorf("分组 '%s' 中已存在名称为 '%s' 的书签", b.GroupName, b.Title)
	}
	processed, err := s.encryptBookmark(b)
	if err != nil {
		return "", err
	}
	id := utils.GenerateRandomID()
	now := time.Now()
	if err := s.db.BookmarkRepo.InsertBookmark(&database.BookmarkDB{
		ID: id, GroupID: group.ID, Title: processed.Title, Host: processed.Host, Port: processed.Port,
		User: processed.User, Password: processed.Password, PrivateKey: processed.PrivateKey,
		PrivateKeyPassword: processed.PrivateKeyPassword, ProxyJumpID: processed.ProxyJumpID,
		CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		return "", err
	}
	s.emit()
	return id, nil
}

func (s *Service) Delete(id string) error {
	if err := s.db.BookmarkRepo.DeleteBookmark(id); err != nil {
		return err
	}
	s.emit()
	return nil
}

func (s *Service) AddGroup(name string) error {
	if _, err := s.db.BookmarkRepo.GetGroupByName(name); err == nil {
		return fmt.Errorf("分组 '%s' 已存在", name)
	}
	if err := s.db.BookmarkRepo.InsertGroup(&database.BookmarkGroupDB{Name: name}); err != nil {
		return err
	}
	s.emit()
	return nil
}

func (s *Service) UpdateGroup(oldName, newName string) error {
	if oldName == "默认书签" {
		return fmt.Errorf("不能修改默认分组名称")
	}
	if _, err := s.db.BookmarkRepo.GetGroupByName(newName); err == nil {
		return fmt.Errorf("分组 '%s' 已存在", newName)
	}
	if err := s.db.BookmarkRepo.UpdateGroupName(oldName, newName); err != nil {
		return err
	}
	s.emit()
	return nil
}

func (s *Service) DeleteGroup(name string) error {
	if name == "默认书签" {
		return fmt.Errorf("不能删除默认分组")
	}
	count, err := s.db.BookmarkRepo.GetGroupBookmarkCount(name)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("分组 '%s' 下还有 %d 个书签，无法删除", name, count)
	}
	if err := s.db.BookmarkRepo.DeleteGroup(name); err != nil {
		return err
	}
	s.emit()
	return nil
}

func (s *Service) PrepareTest(b Bookmark) (ssh.Endpoint, string, error) {
	if b.ID != "" {
		existing, err := s.Get(b.ID)
		if err != nil {
			return ssh.Endpoint{}, "", err
		}
		if b.Host == existing.Host && b.Port == existing.Port && b.User == existing.User &&
			b.PrivateKey == existing.PrivateKey &&
			(b.Password == "" || b.Password == PasswordMask) &&
			(b.PrivateKeyPassword == "" || b.PrivateKeyPassword == PasswordMask) {
			decrypted, err := s.GetDecrypted(b.ID)
			if err != nil {
				return ssh.Endpoint{}, "", err
			}
			return decrypted.Endpoint(), decrypted.ProxyJumpID, nil
		}
	}
	return b.Endpoint(), b.ProxyJumpID, nil
}
