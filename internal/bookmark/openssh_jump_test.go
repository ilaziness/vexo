package bookmark

import (
	"os"
	"path/filepath"
	"testing"
)

func TestImportMultiHopDoesNotClobberOwnJump(t *testing.T) {
	svc, _ := testBookmarkService(t, "", "")
	if err := svc.ensureGroup(ImportedGroupName); err != nil {
		t.Fatal(err)
	}

	dir := t.TempDir()
	cfg := filepath.Join(dir, "config")
	content := `
Host a
  HostName 10.0.0.1
  User u

Host b
  HostName 10.0.0.2
  User u
  ProxyJump a

Host c
  HostName 10.0.0.3
  User u

Host target
  HostName 10.0.0.9
  User u
  ProxyJump c,b
`
	if err := os.WriteFile(cfg, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	res, err := svc.ImportOpenSSH(cfg, ImportedGroupName)
	if err != nil {
		t.Fatal(err)
	}
	if res.Created < 4 {
		t.Fatalf("created=%d warnings=%v", res.Created, res.Warnings)
	}

	items, err := svc.ListItems()
	if err != nil {
		t.Fatal(err)
	}
	idOf := map[string]string{}
	for _, it := range items {
		if it.GroupName == ImportedGroupName {
			idOf[it.Title] = it.ID
		}
	}
	b, err := svc.Get(idOf["b"])
	if err != nil {
		t.Fatal(err)
	}
	if b.ProxyJumpID != idOf["a"] {
		t.Fatalf("b.proxy_jump_id=%q want a=%q warnings=%v", b.ProxyJumpID, idOf["a"], res.Warnings)
	}
	target, err := svc.Get(idOf["target"])
	if err != nil {
		t.Fatal(err)
	}
	if target.ProxyJumpID != idOf["b"] {
		t.Fatalf("target.proxy_jump_id=%q want b=%q", target.ProxyJumpID, idOf["b"])
	}
}
