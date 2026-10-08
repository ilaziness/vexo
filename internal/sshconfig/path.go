package sshconfig

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// ExpandPath expands leading ~/ and ~user/ to an absolute path.
// On failure it returns the original path and a non-nil error.
func ExpandPath(p string) (string, error) {
	p = unquoteSSH(strings.TrimSpace(p))
	if p == "" {
		return "", nil
	}
	if !strings.HasPrefix(p, "~") {
		return p, nil
	}
	var rest string
	var home string
	if p == "~" || strings.HasPrefix(p, "~/") || strings.HasPrefix(p, `~\`) {
		rest = strings.TrimPrefix(strings.TrimPrefix(p[1:], "/"), `\`)
		h, err := os.UserHomeDir()
		if err != nil {
			return p, err
		}
		home = h
	} else {
		// ~otheruser/...
		body := p[1:]
		slash := strings.IndexAny(body, `/\`)
		var uname string
		if slash < 0 {
			uname = body
			rest = ""
		} else {
			uname = body[:slash]
			rest = body[slash+1:]
		}
		u, err := user.Lookup(uname)
		if err != nil {
			return p, err
		}
		home = u.HomeDir
	}
	if rest == "" {
		return home, nil
	}
	return filepath.Join(home, filepath.FromSlash(rest)), nil
}
