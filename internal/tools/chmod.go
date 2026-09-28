package tools

import (
	"fmt"
	"strconv"
	"strings"
)

// ChmodResult 权限计算结果
type ChmodResult struct {
	Success  bool   `json:"success"`
	Octal    string `json:"octal,omitempty"`
	Symbolic string `json:"symbolic,omitempty"`
	Error    string `json:"error,omitempty"`
}

// ConvertChmod 在八进制与 rwx 符号表示之间转换
// direction: toSymbolic | toOctal
func (ts *Service) ConvertChmod(input string, direction string) ChmodResult {
	input = strings.TrimSpace(input)
	if input == "" {
		return ChmodResult{Success: false, Error: "输入内容不能为空"}
	}

	switch strings.ToLower(direction) {
	case "tosymbolic", "to-symbolic", "symbolic":
		return octalToSymbolic(input)
	case "tooctal", "to-octal", "octal":
		return symbolicToOctal(input)
	default:
		return ChmodResult{Success: false, Error: "不支持的方向: " + direction}
	}
}

func formatChmodOctal(mode uint16) string {
	if mode <= 0o777 {
		return fmt.Sprintf("%03o", mode)
	}
	return fmt.Sprintf("%04o", mode)
}

func octalToSymbolic(input string) ChmodResult {
	n, err := strconv.ParseUint(input, 8, 16)
	if err != nil {
		return ChmodResult{Success: false, Error: "无效的八进制权限: " + err.Error()}
	}
	if n > 0o7777 {
		return ChmodResult{Success: false, Error: "权限值超出范围 (最大 7777)"}
	}

	mode := uint16(n)
	perms := make([]byte, 9)
	bits := []uint16{0o400, 0o200, 0o100, 0o040, 0o020, 0o010, 0o004, 0o002, 0o001}
	chars := []byte{'r', 'w', 'x', 'r', 'w', 'x', 'r', 'w', 'x'}
	for i := range bits {
		if mode&bits[i] != 0 {
			perms[i] = chars[i]
		} else {
			perms[i] = '-'
		}
	}
	if mode&0o4000 != 0 {
		if perms[2] == 'x' {
			perms[2] = 's'
		} else {
			perms[2] = 'S'
		}
	}
	if mode&0o2000 != 0 {
		if perms[5] == 'x' {
			perms[5] = 's'
		} else {
			perms[5] = 'S'
		}
	}
	if mode&0o1000 != 0 {
		if perms[8] == 'x' {
			perms[8] = 't'
		} else {
			perms[8] = 'T'
		}
	}

	return ChmodResult{Success: true, Octal: formatChmodOctal(mode), Symbolic: string(perms)}
}

func symbolicToOctal(input string) ChmodResult {
	if len(input) == 10 {
		// 兼容 ls -l：首字符为文件类型 (d/- /l/…)
		input = input[1:]
	}
	if len(input) != 9 {
		return ChmodResult{Success: false, Error: "符号权限须为 9 位，例如 rwxr-xr-x"}
	}

	var mode uint16
	setBit := func(ch byte, expect byte, bit uint16) error {
		switch ch {
		case expect:
			mode |= bit
			return nil
		case '-':
			return nil
		default:
			return fmt.Errorf("无效字符: %c（期望 %c/-）", ch, expect)
		}
	}
	setExec := func(ch byte, bit uint16, specialBit uint16, lower byte, upper byte) error {
		switch ch {
		case 'x':
			mode |= bit
			return nil
		case '-':
			return nil
		case lower:
			mode |= specialBit | bit
			return nil
		case upper:
			mode |= specialBit
			return nil
		default:
			return fmt.Errorf("无效字符: %c（期望 x/-/%c/%c）", ch, lower, upper)
		}
	}

	if err := setBit(input[0], 'r', 0o400); err != nil {
		return ChmodResult{Success: false, Error: err.Error()}
	}
	if err := setBit(input[1], 'w', 0o200); err != nil {
		return ChmodResult{Success: false, Error: err.Error()}
	}
	if err := setExec(input[2], 0o100, 0o4000, 's', 'S'); err != nil {
		return ChmodResult{Success: false, Error: err.Error()}
	}
	if err := setBit(input[3], 'r', 0o040); err != nil {
		return ChmodResult{Success: false, Error: err.Error()}
	}
	if err := setBit(input[4], 'w', 0o020); err != nil {
		return ChmodResult{Success: false, Error: err.Error()}
	}
	if err := setExec(input[5], 0o010, 0o2000, 's', 'S'); err != nil {
		return ChmodResult{Success: false, Error: err.Error()}
	}
	if err := setBit(input[6], 'r', 0o004); err != nil {
		return ChmodResult{Success: false, Error: err.Error()}
	}
	if err := setBit(input[7], 'w', 0o002); err != nil {
		return ChmodResult{Success: false, Error: err.Error()}
	}
	if err := setExec(input[8], 0o001, 0o1000, 't', 'T'); err != nil {
		return ChmodResult{Success: false, Error: err.Error()}
	}

	return ChmodResult{Success: true, Octal: formatChmodOctal(mode), Symbolic: input}
}
