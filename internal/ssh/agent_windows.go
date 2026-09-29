//go:build windows

package ssh

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	openSSHAgentPipe = `\\.\pipe\openssh-ssh-agent`
	agentCopyDataID  = 0x804e50ba
	agentMaxMsgLen   = 256 * 1024
	wmCopyData       = 0x004A
)

type copyDataStruct struct {
	dwData uintptr
	cbData uint32
	lpData uintptr
}

type winPipe struct {
	h windows.Handle
}

func (p *winPipe) Read(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}
	var n uint32
	err := windows.ReadFile(p.h, b, &n, nil)
	if n == 0 && err == nil {
		return 0, io.EOF
	}
	return int(n), err
}

func (p *winPipe) Write(b []byte) (int, error) {
	total := 0
	for total < len(b) {
		var n uint32
		err := windows.WriteFile(p.h, b[total:], &n, nil)
		total += int(n)
		if err != nil {
			return total, err
		}
		if n == 0 {
			return total, io.ErrShortWrite
		}
	}
	return total, nil
}

func (p *winPipe) Close() error {
	return windows.CloseHandle(p.h)
}

type pageantConn struct {
	buf   []byte
	resp  []byte
	off   int
	query func(req []byte) ([]byte, error)
}

func (p *pageantConn) Write(b []byte) (int, error) {
	p.buf = append(p.buf, b...)
	return len(b), nil
}

func (p *pageantConn) Read(b []byte) (int, error) {
	if len(p.resp) == 0 || p.off >= len(p.resp) {
		req := append([]byte(nil), p.buf...)
		p.buf = nil
		resp, err := p.query(req)
		if err != nil {
			return 0, err
		}
		if len(resp) == 0 {
			return 0, io.EOF
		}
		p.resp = resp
		p.off = 0
	}
	n := copy(b, p.resp[p.off:])
	p.off += n
	if p.off >= len(p.resp) {
		p.resp = nil
		p.off = 0
	}
	return n, nil
}

func dialAgentEndpoints() ([]agentEndpoint, error) {
	var eps []agentEndpoint
	var errs []error
	seen := map[string]struct{}{}
	tryPipe := func(path string) {
		path = normalizePipe(path)
		if !isPipePath(path) {
			return
		}
		key := strings.ToLower(path)
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		conn, err := dialPipe(path)
		if err != nil {
			errs = append(errs, err)
			return
		}
		eps = append(eps, agentEndpoint{rw: conn, closer: conn})
	}
	if sock := os.Getenv("SSH_AUTH_SOCK"); sock != "" {
		tryPipe(sock)
	}
	tryPipe(openSSHAgentPipe)
	if _, err := findPageant(); err != nil {
		errs = append(errs, err)
	} else {
		eps = append(eps, agentEndpoint{rw: &pageantConn{query: queryPageant}})
	}
	if len(eps) == 0 {
		return nil, errors.Join(errs...)
	}
	return eps, nil
}

func normalizePipe(path string) string {
	path = strings.TrimSpace(path)
	lower := strings.ToLower(path)
	const alt = `//./pipe/`
	if strings.HasPrefix(lower, alt) {
		return `\\.\pipe\` + path[len(alt):]
	}
	return path
}

func isPipePath(path string) bool {
	return strings.HasPrefix(strings.ToLower(path), `\\.\pipe\`)
}

func dialPipe(path string) (io.ReadWriteCloser, error) {
	name, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}
	h, err := windows.CreateFile(name, windows.GENERIC_READ|windows.GENERIC_WRITE, 0, nil, windows.OPEN_EXISTING, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("dial %s: %w", path, err)
	}
	return &winPipe{h: h}, nil
}

func findPageant() (uintptr, error) {
	user32 := windows.NewLazySystemDLL("user32.dll")
	findWindow := user32.NewProc("FindWindowW")
	class, err := windows.UTF16PtrFromString("Pageant")
	if err != nil {
		return 0, err
	}
	hwnd, _, _ := findWindow.Call(uintptr(unsafe.Pointer(class)), uintptr(unsafe.Pointer(class)))
	runtime.KeepAlive(class)
	if hwnd == 0 {
		return 0, fmt.Errorf("Pageant is not running")
	}
	return hwnd, nil
}

func queryPageant(req []byte) ([]byte, error) {
	if len(req) == 0 || len(req) > agentMaxMsgLen {
		return nil, fmt.Errorf("pageant request length invalid")
	}
	hwnd, err := findPageant()
	if err != nil {
		return nil, err
	}
	sa, err := currentUserMappingSA()
	if err != nil {
		return nil, err
	}
	mapping, mapName, err := createPageantMapping(sa)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(mapping)
	runtime.KeepAlive(sa)
	addr, err := windows.MapViewOfFile(mapping, windows.FILE_MAP_READ|windows.FILE_MAP_WRITE, 0, 0, 0)
	if err != nil {
		return nil, err
	}
	defer windows.UnmapViewOfFile(addr)

	view := unsafe.Slice((*byte)(unsafe.Pointer(addr)), agentMaxMsgLen)
	copy(view, req)

	nameC := append([]byte(mapName), 0)
	cds := copyDataStruct{
		dwData: agentCopyDataID,
		cbData: uint32(len(nameC)),
		lpData: uintptr(unsafe.Pointer(&nameC[0])),
	}
	user32 := windows.NewLazySystemDLL("user32.dll")
	sendMessage := user32.NewProc("SendMessageW")
	r, _, callErr := sendMessage.Call(hwnd, wmCopyData, 0, uintptr(unsafe.Pointer(&cds)))
	runtime.KeepAlive(nameC)
	runtime.KeepAlive(&cds)
	if r == 0 {
		if callErr != nil && !errors.Is(callErr, windows.ERROR_SUCCESS) {
			return nil, fmt.Errorf("pageant: %w", callErr)
		}
		return nil, fmt.Errorf("pageant rejected the request")
	}
	n := binary.BigEndian.Uint32(view[:4])
	total := int(n) + 4
	if total < 4 || total > len(view) {
		return nil, fmt.Errorf("pageant response length invalid")
	}
	out := make([]byte, total)
	copy(out, view[:total])
	return out, nil
}

func randomPageantMapName() (string, error) {
	var rnd [4]byte
	if _, err := rand.Read(rnd[:]); err != nil {
		return "", err
	}
	return fmt.Sprintf("PageantRequest%08x%08x", uint32(os.Getpid()), binary.BigEndian.Uint32(rnd[:])), nil
}

// createPageantMapping 建立仅当前用户可打开的共享内存。名字冲突时关掉句柄再换一个名字。
func createPageantMapping(sa *windows.SecurityAttributes) (windows.Handle, string, error) {
	var lastErr error
	for attempt := 0; attempt < 4; attempt++ {
		mapName, err := randomPageantMapName()
		if err != nil {
			return 0, "", err
		}
		nameUTF16, err := windows.UTF16PtrFromString(mapName)
		if err != nil {
			return 0, "", err
		}
		mapping, err := windows.CreateFileMapping(windows.InvalidHandle, sa, windows.PAGE_READWRITE, 0, uint32(agentMaxMsgLen), nameUTF16)
		runtime.KeepAlive(nameUTF16)
		if err == nil {
			return mapping, mapName, nil
		}
		if mapping != 0 {
			_ = windows.CloseHandle(mapping)
		}
		lastErr = err
		if !errors.Is(err, windows.ERROR_ALREADY_EXISTS) {
			return 0, "", err
		}
	}
	return 0, "", lastErr
}

func currentUserMappingSA() (*windows.SecurityAttributes, error) {
	user, err := windows.GetCurrentProcessToken().GetTokenUser()
	if err != nil {
		return nil, err
	}
	sd, err := windows.SecurityDescriptorFromString("D:P(A;;GA;;;" + user.User.Sid.String() + ")")
	if err != nil {
		return nil, err
	}
	return &windows.SecurityAttributes{
		Length:             uint32(unsafe.Sizeof(windows.SecurityAttributes{})),
		SecurityDescriptor: sd,
	}, nil
}
