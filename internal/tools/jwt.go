package tools

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"strings"
)

// JWTResult JWT 解码结果
type JWTResult struct {
	Success   bool   `json:"success"`
	Header    string `json:"header,omitempty"`
	Payload   string `json:"payload,omitempty"`
	Signature string `json:"signature,omitempty"`
	Error     string `json:"error,omitempty"`
}

const bearerPrefix = "bearer "

// DecodeJWT 本地拆分 JWT，不校验签名
func (ts *Service) DecodeJWT(token string) JWTResult {
	token = strings.TrimSpace(token)
	if token == "" {
		return JWTResult{Success: false, Error: "JWT 不能为空"}
	}

	if len(token) >= len(bearerPrefix) && strings.EqualFold(token[:len(bearerPrefix)], bearerPrefix) {
		token = strings.TrimSpace(token[len(bearerPrefix):])
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return JWTResult{Success: false, Error: "无效的 JWT：应为 header.payload.signature 三段"}
	}
	if parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return JWTResult{Success: false, Error: "无效的 JWT：header / payload / signature 不能为空"}
	}

	header, err := decodeJWTPart(parts[0])
	if err != nil {
		return JWTResult{Success: false, Error: "Header 解码失败: " + err.Error()}
	}
	payload, err := decodeJWTPart(parts[1])
	if err != nil {
		return JWTResult{Success: false, Error: "Payload 解码失败: " + err.Error()}
	}

	return JWTResult{
		Success:   true,
		Header:    header,
		Payload:   payload,
		Signature: parts[2],
	}
}

func decodeJWTPart(part string) (string, error) {
	raw, err := base64.RawURLEncoding.DecodeString(part)
	if err != nil {
		raw, err = base64.URLEncoding.DecodeString(part)
		if err != nil {
			return "", err
		}
	}
	var buf bytes.Buffer
	if err := json.Indent(&buf, raw, "", "  "); err != nil {
		return string(raw), nil
	}
	return buf.String(), nil
}
