package services

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"html"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"

	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
)

// Tool 工具定义
type Tool struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Category    string `json:"category"`
}

// PortCheckResult 端口检测结果
type PortCheckResult struct {
	Success      bool   `json:"success"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	ResponseTime int64  `json:"responseTime"` // 毫秒
	Error        string `json:"error,omitempty"`
}

// EncodeRequest 编码请求
type EncodeRequest struct {
	ToolType string `json:"toolType"` // base64, url, html
	Input    string `json:"input"`
}

// EncodeResponse 编码响应
type EncodeResponse struct {
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}

// Match 正则匹配结果
type Match struct {
	Text   string   `json:"text"`
	Index  int      `json:"index"`
	Groups []string `json:"groups"`
}

// RegexMatchResult 正则匹配结果
type RegexMatchResult struct {
	Matches []Match `json:"matches"`
	Count   int     `json:"count"`
	Error   string  `json:"error,omitempty"`
}

// ToolService 工具服务
type ToolService struct{}

// NewToolService 创建工具服务实例
func NewToolService() *ToolService {
	return &ToolService{}
}

// ShowWindow 显示工具窗口
func (ts *ToolService) ShowWindow() {
	if AppInstance.ToolWindow == nil {
		AppInstance.ToolWindow = newToolWindow()
	}
	AppInstance.ToolWindow.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
		AppInstance.ToolWindow = nil
	})
	AppInstance.ToolWindow.Show()
	AppInstance.ToolWindow.Focus()
}

// CloseWindow 关闭工具窗口
func (ts *ToolService) CloseWindow() {
	if AppInstance.ToolWindow != nil {
		AppInstance.ToolWindow.Close()
		AppInstance.ToolWindow = nil
	}
}

// GetTools 获取可用工具列表
func (ts *ToolService) GetTools() []Tool {
	return []Tool{
		{
			ID:          "port-check",
			Name:        "端口检测",
			Description: "TCP端口连通性测试工具",
			Icon:        "NetworkCheck",
			Category:    "网络",
		},
		{
			ID:          "encoder",
			Name:        "编码解码",
			Description: "Base64、URL、HTML编解码工具",
			Icon:        "Code",
			Category:    "编码",
		},
		{
			ID:          "regex",
			Name:        "正则测试",
			Description: "正则表达式匹配测试工具",
			Icon:        "RegularExpression",
			Category:    "开发",
		},
		{
			ID:          "hash",
			Name:        "哈希计算",
			Description: "文件和字符串哈希值计算（MD5、SHA1、SHA256等）",
			Icon:        "EditNote",
			Category:    "安全",
		},
		{
			ID:          "timestamp",
			Name:        "时间戳转换",
			Description: "Unix时间戳和日期时间相互转换工具",
			Icon:        "Code",
			Category:    "开发",
		},
	}
}

// CheckPort 检测端口连通性
func (ts *ToolService) CheckPort(host string, port int) PortCheckResult {
	if host == "" {
		return PortCheckResult{
			Success: false,
			Host:    host,
			Port:    port,
			Error:   "主机地址不能为空",
		}
	}
	if port < 1 || port > 65535 {
		return PortCheckResult{
			Success: false,
			Host:    host,
			Port:    port,
			Error:   "端口号必须在1-65535之间",
		}
	}

	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	start := time.Now()
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	elapsed := time.Since(start).Milliseconds()

	if err != nil {
		return PortCheckResult{
			Success:      false,
			Host:         host,
			Port:         port,
			ResponseTime: elapsed,
			Error:        err.Error(),
		}
	}
	defer conn.Close()

	return PortCheckResult{
		Success:      true,
		Host:         host,
		Port:         port,
		ResponseTime: elapsed,
	}
}

// Encode 编码
func (ts *ToolService) Encode(toolType string, input string) EncodeResponse {
	var result string

	switch toolType {
	case "base64":
		result = base64.StdEncoding.EncodeToString([]byte(input))
	case "url":
		result = url.QueryEscape(input)
	case "html":
		result = html.EscapeString(input)
	default:
		return EncodeResponse{
			Error: "不支持的编码类型: " + toolType,
		}
	}

	return EncodeResponse{Result: result}
}

// Decode 解码
func (ts *ToolService) Decode(toolType string, input string) EncodeResponse {
	var result string
	var err error

	switch toolType {
	case "base64":
		var decoded []byte
		decoded, err = base64.StdEncoding.DecodeString(input)
		if err != nil {
			return EncodeResponse{Error: "Base64解码失败: " + err.Error()}
		}
		result = string(decoded)
	case "url":
		result, err = url.QueryUnescape(input)
		if err != nil {
			return EncodeResponse{Error: "URL解码失败: " + err.Error()}
		}
	case "html":
		result = html.UnescapeString(input)
	default:
		return EncodeResponse{
			Error: "不支持的解码类型: " + toolType,
		}
	}

	return EncodeResponse{Result: result}
}

// HashResult 哈希计算结果
type HashResult struct {
	Success   bool   `json:"success"`
	Input     string `json:"input"`
	Algorithm string `json:"algorithm"` // 哈希算法
	Result    string `json:"result"`    // 哈希值
	Error     string `json:"error,omitempty"`
}

// TimestampResult 时间戳转换结果
type TimestampResult struct {
	Success   bool   `json:"success"`
	Timestamp int64  `json:"timestamp,omitempty"` // Unix时间戳（秒）
	Datetime  string `json:"datetime,omitempty"`  // 日期时间字符串
	Error     string `json:"error,omitempty"`
}

// RegexMatch 正则匹配
func (ts *ToolService) RegexMatch(pattern string, text string, flags string) RegexMatchResult {
	if pattern == "" {
		return RegexMatchResult{Matches: []Match{}, Count: 0}
	}

	// 构建正则标志
	flagMap := map[rune]bool{
		'i': false, // 忽略大小写
		'm': false, // 多行模式
		'g': false, // 全局匹配
	}
	for _, f := range flags {
		if _, ok := flagMap[f]; ok {
			flagMap[f] = true
		}
	}

	// 构建 Go 正则标志
	goFlags := ""
	if flagMap['i'] {
		goFlags += "(?i)"
	}
	if flagMap['m'] {
		goFlags += "(?m)"
	}

	// 编译正则表达式
	fullPattern := goFlags + pattern
	re, err := regexp.Compile(fullPattern)
	if err != nil {
		return RegexMatchResult{
			Error: "正则表达式编译失败: " + err.Error(),
		}
	}

	// 查找所有匹配
	matches := re.FindAllStringSubmatchIndex(text, -1)
	if matches == nil {
		return RegexMatchResult{Matches: []Match{}, Count: 0}
	}

	result := make([]Match, 0, len(matches))
	for _, match := range matches {
		if len(match) >= 2 {
			start, end := match[0], match[1]
			matchedText := text[start:end]

			// 提取捕获组
			groups := make([]string, 0)
			for i := 2; i < len(match); i += 2 {
				if i+1 < len(match) {
					groupStart, groupEnd := match[i], match[i+1]
					if groupStart >= 0 && groupEnd >= 0 {
						groups = append(groups, text[groupStart:groupEnd])
					} else {
						groups = append(groups, "")
					}
				}
			}

			result = append(result, Match{
				Text:   matchedText,
				Index:  start,
				Groups: groups,
			})
		}
	}

	return RegexMatchResult{
		Matches: result,
		Count:   len(result),
	}
}

// CalculateHash 计算哈希值
func (ts *ToolService) CalculateHash(input string, algorithm string) HashResult {
	if input == "" {
		return HashResult{
			Success: false,
			Input:   input,
			Error:   "输入内容不能为空",
		}
	}
	if algorithm == "" {
		algorithm = "md5"
	}

	var result string

	switch strings.ToLower(algorithm) {
	case "md5":
		hash := md5.Sum([]byte(input))
		result = hex.EncodeToString(hash[:])
	case "sha1":
		hash := sha1.Sum([]byte(input))
		result = hex.EncodeToString(hash[:])
	case "sha256":
		hash := sha256.Sum256([]byte(input))
		result = hex.EncodeToString(hash[:])
	case "sha512":
		hash := sha512.Sum512([]byte(input))
		result = hex.EncodeToString(hash[:])
	default:
		return HashResult{
			Success: false,
			Input:   input,
			Error:   "不支持的哈希算法: " + algorithm,
		}
	}

	return HashResult{
		Success:   true,
		Input:     input,
		Algorithm: strings.ToLower(algorithm),
		Result:    result,
	}
}

// ConvertTimestamp 转换时间戳
func (ts *ToolService) ConvertTimestamp(input string, toTimestamp bool) TimestampResult {
	if input == "" {
		return TimestampResult{
			Success: false,
			Error:   "输入内容不能为空",
		}
	}

	if toTimestamp {
		// 将日期时间转换为时间戳
		// 支持多种日期时间格式
		layouts := []string{
			"2006-01-02 15:04:05",
			"2006-01-02 15:04",
			"2006-01-02",
			time.RFC3339,
			time.RFC1123,
			time.RFC822,
			time.ANSIC,
			time.UnixDate,
			time.RubyDate,
		}

		var t time.Time
		var err error

		for _, layout := range layouts {
			t, err = time.Parse(layout, input)
			if err == nil {
				break
			}
		}

		if err != nil {
			return TimestampResult{
				Success: false,
				Error:   "无法解析日期时间格式: " + err.Error(),
			}
		}

		return TimestampResult{
			Success:   true,
			Timestamp: t.Unix(),
		}
	} else {
		// 将时间戳转换为日期时间
		timestamp, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return TimestampResult{
				Success: false,
				Error:   "无效的时间戳: " + err.Error(),
			}
		}

		t := time.Unix(timestamp, 0)
		return TimestampResult{
			Success:  true,
			Datetime: t.Format("2006-01-02 15:04:05"),
		}
	}
}
