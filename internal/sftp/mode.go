package sftp

import "os"

// fileModeToUnix maps Go FileMode permission/special bits to Unix mode (0o7777).
func fileModeToUnix(m os.FileMode) uint32 {
	unix := uint32(m.Perm())
	if m&os.ModeSetuid != 0 {
		unix |= 0o4000
	}
	if m&os.ModeSetgid != 0 {
		unix |= 0o2000
	}
	if m&os.ModeSticky != 0 {
		unix |= 0o1000
	}
	return unix
}

// unixToFileMode maps Unix mode (0o7777) to Go FileMode for pkg/sftp Chmod.
func unixToFileMode(unix uint32) os.FileMode {
	m := os.FileMode(unix & 0o777)
	if unix&0o4000 != 0 {
		m |= os.ModeSetuid
	}
	if unix&0o2000 != 0 {
		m |= os.ModeSetgid
	}
	if unix&0o1000 != 0 {
		m |= os.ModeSticky
	}
	return m
}
