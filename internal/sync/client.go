package sync

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Client 同步客户端
type Client struct {
	config     *SyncConfig
	httpClient *http.Client
	baseURL    string
}

// VersionInfo 版本信息
type VersionInfo struct {
	VersionNumber int       `json:"version_number"`
	FileSize      int64     `json:"file_size"`
	CreatedAt     time.Time `json:"created_at"`
}

// NewClient 创建同步客户端
func NewClient(config *SyncConfig) *Client {
	return &Client{
		config:     config,
		httpClient: &http.Client{Timeout: 5 * time.Minute},
		baseURL:    normalizeURL(config.ServerURL),
	}
}

// normalizeURL 规范化 URL，确保有 http:// 或 https:// 前缀
func normalizeURL(url string) string {
	if url == "" {
		return ""
	}
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	}
	return "http://" + url
}

// setAuthHeaders 设置认证头
func (c *Client) setAuthHeaders(req *http.Request) {
	req.Header.Set("X-Sync-ID", c.config.SyncID)
	req.Header.Set("X-User-Key", c.config.UserKey)
}

// Upload 上传数据
func (c *Client) Upload(reader io.Reader) (*VersionInfo, error) {
	url := c.baseURL + "/upload"

	req, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setAuthHeaders(req)
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to upload: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("upload failed: %s - %s", resp.Status, string(body))
	}

	var result VersionInfo
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// Download 下载数据
func (c *Client) Download(version int) (io.ReadCloser, error) {
	url := c.baseURL + "/download"
	if version > 0 {
		url += "?version=" + strconv.Itoa(version)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setAuthHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("download failed: %s - %s", resp.Status, string(body))
	}

	return resp.Body, nil
}

// ListVersions 获取版本列表（支持分页，limit=0表示不分页）
func (c *Client) ListVersions(limit, offset int) ([]VersionInfo, error) {
	url := c.baseURL + "/versions"
	if limit > 0 {
		url += fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setAuthHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to list versions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("list versions failed: %s - %s", resp.Status, string(body))
	}

	var result struct {
		Versions []VersionInfo `json:"versions"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Versions, nil
}

// DeleteVersion 删除指定版本
func (c *Client) DeleteVersion(version int) error {
	url := c.baseURL + "/versions?version=" + strconv.Itoa(version)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	c.setAuthHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete version: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete version failed: %s - %s", resp.Status, string(body))
	}

	return nil
}

// HealthCheck 健康检查
func (c *Client) HealthCheck() error {
	url := c.baseURL + "/health"

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed: %s", resp.Status)
	}

	return nil
}

// PipeReader 管道读取器，用于流式上传
type PipeReader struct {
	reader *io.PipeReader
	done   chan error
}

// NewPipeReader 创建管道读取器
func NewPipeReader() (*PipeReader, *io.PipeWriter) {
	pr, pw := io.Pipe()
	return &PipeReader{
		reader: pr,
		done:   make(chan error, 1),
	}, pw
}

// Read 实现 io.Reader 接口
func (p *PipeReader) Read(data []byte) (int, error) {
	return p.reader.Read(data)
}

// Close 关闭读取器
func (p *PipeReader) Close() error {
	return p.reader.Close()
}
