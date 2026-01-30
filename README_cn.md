# Vexo — 跨平台 SSH & FTP 图形客户端

[![Go Version](https://img.shields.io/badge/Go-1.22%2B-blue)](https://golang.org)
[![Wails Version](https://img.shields.io/badge/Wails-v3-8A2BE2)](https://wails.io)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

> 基于 Go 与 Wails v3 构建的现代化跨平台 SSH 与 FTP/SFTP 图形客户端。

---

## 🚀 简介

**Vexo** 是一款简洁高效的桌面应用，将 SSH 与 FTP/SFTP 功能集成到直观的图形界面中。基于 [Wails v3](https://wails.io) 和 Go 语言开发，支持 Windows、macOS 和 Linux，为开发者与运维人员提供一体化的远程连接与文件管理体验。

### ✨ 主要功能

- **SSH 终端**：内嵌终端，支持多标签页
- **认证方式**：密码或私钥登录（支持 RSA、Ed25519、ECDSA 等格式）
- **SFTP 文件浏览器**：上传、下载、删除、重命名、预览文件
- **会话管理**：保存、分类并快速连接常用服务器
- **端口转发**：支持本地/远程端口转发（TCP 隧道）
- **主题切换**：暗色/亮色模式，支持自动跟随系统
- **跨平台**：一套代码，原生运行于三大桌面系统
- **安全存储**：凭据本地加密（未来支持系统钥匙串集成）

### 🛠 技术栈

- **前端**：Vue 3 + TypeScript + Tailwind CSS（由 Wails v3 模板生成）
- **后端**：Go 1.22+
- **核心依赖**：
  - `golang.org/x/crypto/ssh` – 官方 SSH 实现
  - `github.com/pkg/sftp` – SFTP 协议支持
  - `github.com/go-ftp/ftp` – FTP 客户端功能
  - `github.com/atotto/clipboard` – 剪贴板交互
- **打包发布**：Wails v3 原生构建（计划支持自动更新）

### 📦 安装方式

#### 从源码构建

1. 确保已安装 [Go 1.22+](https://golang.org/dl/) 和 [Node.js 18+](https://nodejs.org/)
2. 克隆仓库：
   ```bash
   git clone https://github.com/your-username/vexo.git
   cd vexo
   ```
3. 安装依赖：
   ```bash
   go mod tidy
   npm install
   ```
4. 构建：
   ```bash
   # Build for Windows (from any OS)
   wails3 build -platform windows/amd64

   # Build for macOS
   wails3 build -platform darwin/amd64
   wails3 build -platform darwin/arm64

   # Build for Linux
   wails3 build -platform linux/amd64
   ```

#### 预编译版本

即将发布！请关注本仓库的 Release 页面。

### 🖼️ 截图预览

*（首次发布后添加截图）*  
- 多标签 SSH 终端（带语法高亮）  
- 支持拖拽的 SFTP 文件管理器  
- 会话配置面板  

### 🤝 参与贡献

欢迎提交 Issue 或 Pull Request！

### 📄 开源协议

MIT © [你的名字/组织] – 详见 [LICENSE](LICENSE)
