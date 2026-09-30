package bookmark

import (
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/ilaziness/vexo/internal/database"
	"github.com/ilaziness/vexo/internal/secret"
)

func TestPrepareTestSwitchToStoredKeyKeepsPassword(t *testing.T) {
	svc, id := testBookmarkService(t, "/keys/id_ed25519", "key-pass")
	ep, _, err := svc.PrepareTest(Bookmark{
		ID: id, Host: "h", Port: 22, User: "root",
		Password: PasswordMask, PrivateKey: "", PrivateKeyPassword: "",
		SshKeyID: "key-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if ep.Password != "login-secret" {
		t.Fatalf("password = %q", ep.Password)
	}
	if ep.Key != "" || ep.KeyPassword != "" {
		t.Fatalf("file key leaked: path=%q pass=%q", ep.Key, ep.KeyPassword)
	}
	if ep.KeyPEM != "PEM" {
		t.Fatalf("key pem = %q", ep.KeyPEM)
	}
}

func TestPrepareTestClearedPasswordStaysEmpty(t *testing.T) {
	svc, id := testBookmarkService(t, "/keys/id_ed25519", "key-pass")
	ep, _, err := svc.PrepareTest(Bookmark{
		ID: id, Host: "h", Port: 22, User: "root",
		Password: "", PrivateKey: "", PrivateKeyPassword: "",
		SshKeyID: "key-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if ep.Password != "" {
		t.Fatalf("password = %q", ep.Password)
	}
	if ep.KeyPEM != "PEM" {
		t.Fatalf("key pem = %q", ep.KeyPEM)
	}
}

func TestPrepareTestUnchangedFileKeepsKeyPassword(t *testing.T) {
	svc, id := testBookmarkService(t, "/keys/id_ed25519", "key-pass")
	ep, jumpID, err := svc.PrepareTest(Bookmark{
		ID: id, Host: "h", Port: 22, User: "root",
		Password: PasswordMask, PrivateKey: "/keys/id_ed25519", PrivateKeyPassword: PasswordMask,
		ProxyJumpID: "jump-new",
	})
	if err != nil {
		t.Fatal(err)
	}
	if jumpID != "jump-new" {
		t.Fatalf("jump = %q", jumpID)
	}
	if ep.Password != "login-secret" || ep.KeyPassword != "key-pass" || ep.Key != "/keys/id_ed25519" {
		t.Fatalf("endpoint = %+v", ep)
	}
	if ep.KeyPEM != "" {
		t.Fatal("unexpected stored key")
	}
}

func TestPrepareTestKeepsFormCertificate(t *testing.T) {
	svc, id := testBookmarkService(t, "/keys/id_ed25519", "key-pass")
	ep, _, err := svc.PrepareTest(Bookmark{
		ID: id, Host: "h", Port: 22, User: "root",
		Password: PasswordMask, PrivateKey: "/keys/id_ed25519", PrivateKeyPassword: PasswordMask,
		Certificate: "/keys/user-cert.pub",
	})
	if err != nil {
		t.Fatal(err)
	}
	if ep.Certificate != "/keys/user-cert.pub" {
		t.Fatalf("certificate = %q", ep.Certificate)
	}
	if ep.Password != "login-secret" || ep.KeyPassword != "key-pass" {
		t.Fatalf("endpoint = %+v", ep)
	}
}

func testBookmarkService(t *testing.T, keyPath, keyPass string) (*Service, string) {
	t.Helper()
	db := database.NewDatabase(t.TempDir(), zap.NewNop())
	if err := db.Initialize(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := db.BookmarkRepo.InsertGroup(&database.BookmarkGroupDB{Name: "默认书签"}); err != nil {
		t.Fatal(err)
	}
	group, err := db.BookmarkRepo.GetGroupByName("默认书签")
	if err != nil {
		t.Fatal(err)
	}
	login, err := secret.Encrypt("master", "login-secret")
	if err != nil {
		t.Fatal(err)
	}
	keyPassword, err := secret.Encrypt("master", keyPass)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	const id = "bm-1"
	if err := db.BookmarkRepo.InsertBookmark(&database.BookmarkDB{
		ID: id, GroupID: group.ID, Title: "t", Host: "h", Port: 22, User: "root",
		Password: login, PrivateKey: keyPath, PrivateKeyPassword: keyPassword,
		ProxyJumpID: "jump-old", CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatal(err)
	}
	svc := New(zap.NewNop(), db, func(string) (string, error) { return "master", nil }, nil, nil, func(keyID string) (string, error) {
		if keyID != "key-1" {
			t.Fatalf("unexpected key id %s", keyID)
		}
		return "PEM", nil
	})
	return svc, id
}
