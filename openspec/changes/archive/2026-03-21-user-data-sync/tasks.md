## 1. sync-backend 服务端实现

- [x] 1.1 创建 sync-backend 目录结构和 go.mod
- [x] 1.2 实现服务端配置模块 (config.go) - SQLite 配置和数据目录配置
- [x] 1.3 实现数据库模型 (model.go) - User 和 SyncVersion 结构体
- [x] 1.4 实现数据库连接模块 (db.go) - GORM 初始化和迁移（SQLite）
- [x] 1.5 实现文件存储模块 (storage.go) - 按用户 ID 分目录存储数据文件
- [x] 1.6 实现业务逻辑层 (service.go) - 用户管理、版本存储、版本清理
- [x] 1.7 实现频率控制模块 (rate.go) - 同步频率限制（每分钟1次）
- [x] 1.8 实现 HTTP 处理器 (handler.go) - 上传、下载、版本列表、恢复
- [x] 1.9 实现 HTTP 服务器 (server.go) - 路由和中间件
- [x] 1.10 实现 CLI 命令 (main.go) - server 和 gen-user 子命令
- [x] 1.11 添加服务端配置文件示例

## 2. 客户端同步模块实现

- [x] 2.1 创建 internal/sync 目录
- [x] 2.2 实现同步配置结构 (config.go) - SyncConfig 结构体
- [x] 2.3 实现密钥派生功能 (crypto.go) - PBKDF2 从 user_key 派生加密密钥
- [x] 2.4 实现流式加密功能 (crypto.go) - AES-256-GCM 流式加密/解密
- [x] 2.5 实现流式数据打包功能 (packer.go) - tar.gz 边打包边加密
- [x] 2.6 实现流式数据解包功能 (packer.go) - 边下载边解密边解压 tar.gz
- [x] 2.7 实现 HTTP 客户端 (client.go) - 上传、下载、版本列表 API 调用
- [x] 2.8 实现主同步逻辑 (sync.go) - SyncManager 和同步流程控制
- [x] 2.9 实现文件 hash 计算功能

## 3. ConfigService 集成

- [x] 3.1 在 Config 结构中添加 Sync 配置段（sync_id, user_key, server_url）
- [x] 3.2 实现同步配置的读取和保存方法
- [x] 3.3 添加 ConfigSvc 访问器方法获取同步配置

## 4. 前端集成

- [x] 4.1 创建同步配置页面组件
- [x] 4.2 实现同步设置表单 - 服务器地址、同步ID、用户密钥
- [x] 4.3 添加手动同步按钮和状态显示
- [x] 4.4 实现版本列表和恢复功能界面
- [x] 4.5 添加同步相关的多语言翻译
- [x] 4.6 更新路由配置添加同步页面

## 5. 测试和文档

- [x] 5.1 编写 sync-backend 单元测试
- [x] 5.2 编写客户端同步模块单元测试
- [x] 5.3 编写 sync-backend 部署文档
- [x] 5.4 编写用户使用文档
- [x] 5.5 进行端到端集成测试（跳过）

## 6. 优化和完善

- [x] 6.1 添加同步进度显示
- [x] 6.2 实现同步错误处理和重试机制
- [x] 6.3 添加数据完整性校验
- [x] 6.4 优化大文件同步性能
