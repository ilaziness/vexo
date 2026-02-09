# Vexo — 跨平台 SSH & SFTP GUI 客户端

[![Go Version](https://img.shields.io/badge/Go-1.25%2B-blue)](https://golang.org)
[![Wails Version](https://img.shields.io/badge/Wails-v3-8A2BE2)](https://wails.io)
[![License](https://img.shields.io/badge/License-Apache%202.0-green.svg)](LICENSE)

> 一款基于 Go 和 Wails v3 构建的现代化、跨平台 SSH 和 SFTP 桌面 GUI 应用。

---

## 📖 简介

**Vexo** 是一个简洁、响应迅速的桌面应用程序，将 SSH 和 SFTP 的强大功能带到您的指尖——无需离开原生 GUI 环境。使用 [Wails v3](https://wails.io) 和 Go 构建，Vexo 原生运行于 Windows、macOS 和 Linux，为开发人员和系统管理员提供统一的远程访问和文件管理工具。

设计理念：清爽简洁、专业可靠、护眼舒适

---

## ✨ 特性

- 嵌入式 SSH 终端，支持多标签页
- 基于 XTerm.js，提供流畅的终端体验
- 支持多种终端配色方案
- SFTP 文件浏览器，支持上传、下载、删除、重命名文件
- 文件传输任务管理
- 密码认证和私钥认证（支持 RSA、Ed25519、ECDSA）
- 敏感信息安全存储
- 会话书签管理，快速连接常用服务器
- 支持文件夹分类管理书签
- 多种主题方案：亮色模式、暗色模式、蓝夜模式、护眼模式
- 敏感信息加密存储，支持系统集成
- 支持 Zmodem 文件传输协议（rz/sz）

---

## 📸 截图

> 以下截图展示了 Vexo 的主要界面：

### 主界面 - SSH 终端
![主界面](screenshots/main.png)
*多标签 SSH 终端界面*

![主界面](screenshots/terminal.png)
*多标签 SSH 终端界面*

### SFTP 文件浏览器
![SFTP](screenshots/sftp.png)
*直观的文件管理界面*

### 书签管理
![书签](screenshots/bookmark.png)
*服务器连接配置管理*

---

## 📄 许可证

本项目采用 [Apache License 2.0](LICENSE) 许可证。

```
Copyright 2024 Vexo Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
