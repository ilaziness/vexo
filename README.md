# Vexo — Cross-Platform SSH & FTP GUI Client

[![Go Version](https://img.shields.io/badge/Go-1.22%2B-blue)](https://golang.org)
[![Wails Version](https://img.shields.io/badge/Wails-v3-8A2BE2)](https://wails.io)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

> A modern, cross-platform GUI client for SSH and FTP/SFTP built with Go and Wails v3.

---

## 🚀 Overview

**Vexo** is a sleek, responsive desktop application that brings the power of SSH and FTP/SFTP to your fingertips—without leaving a native GUI. Built with [Wails v3](https://wails.io) and Go, Vexo runs natively on Windows, macOS, and Linux, offering developers and system administrators a unified tool for remote access and file management.

### ✨ Features

- **SSH Terminal**: Embedded terminal with multi-tab support
- **Authentication**: Password or private key (RSA, Ed25519, ECDSA)
- **SFTP File Browser**: Upload, download, delete, rename, and preview files
- **Session Manager**: Save, organize, and quickly connect to your servers
- **Port Forwarding**: Local and remote port forwarding (TCP tunneling)
- **Dark/Light Theme**: Auto-switch or manual toggle
- **Cross-Platform**: Single codebase, native experience everywhere
- **Secure**: All credentials encrypted at rest (optional OS keychain integration)

### 🛠 Tech Stack

- **Frontend**: Vue 3 + TypeScript + Tailwind CSS (via Wails v3 template)
- **Backend**: Go 1.22+
- **Core Libraries**:
  - `golang.org/x/crypto/ssh` – Robust SSH implementation
  - `github.com/pkg/sftp` – SFTP protocol support
  - `github.com/go-ftp/ftp` – FTP client functionality
  - `github.com/atotto/clipboard` – Clipboard integration
- **Packaging**: Wails v3 native bundling (with auto-updater support planned)

### 📦 Installation

#### From Source

1. Ensure you have [Go 1.22+](https://go.dev/dl/) and [Node.js 18+](https://nodejs.org/) installed.
2. Clone the repo:
   ```bash
   git clone https://github.com/your-username/vexo.git
   cd vexo
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   npm install
   ```
4. Build & run:
   ```bash
   wails build
   ./build/vexo  # or vexo.exe on Windows
   ```

#### Pre-built Binaries

Coming soon! Watch this repo for releases.

### 🖼️ Screenshots

*(Add screenshots after first release)*  
- Terminal tab with syntax-colored output  
- SFTP file browser with drag-and-drop  
- Session management panel  

### 🤝 Contributing

Contributions are welcome! Please open an issue or submit a PR.

### 📄 License

MIT © [Your Name/Org] – see [LICENSE](LICENSE) for details.
