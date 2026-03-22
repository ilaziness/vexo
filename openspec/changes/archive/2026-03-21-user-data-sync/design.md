## Context

Vexo 是一个基于 Go 语言和 Wails v3 框架构建的 SSH 桌面 GUI 应用。当前用户数据（SSH 配置、书签等）仅存储在本地 UserDataDir 目录中，缺乏跨设备同步能力。

当前配置系统：
- ConfigService 管理应用配置，使用 TOML 格式存储
- UserDataDir 路径可通过配置自定义，默认为可执行目录下的 data 文件夹
- 数据库使用 SQLite 存储书签等结构化数据

## Goals / Non-Goals

**Goals:**
- 实现客户端数据流式打包压缩加密并同步到远程服务端
- 实现服务端多版本存储（最多5个版本）和版本恢复
- 支持多用户通过 sync ID 隔离数据
- 服务端使用 SQLite 存储用户信息和版本元数据
- 同步数据文件保存到文件系统
- 提供 CLI 工具生成用户凭证（sync_id 和 user_key）
- 实现同步频率控制，防止滥用

**Non-Goals:**
- 实时同步（仅支持手动触发同步）
- 冲突自动解决（以客户端最新为准）
- 增量同步（全量打包同步）
- 服务端数据解密（服务端只存储加密数据）
- 数据导入导出（不通过服务端）

## Decisions

### 1. 加密方案选择

**决策**: 使用 AES-256-GCM 进行流式数据加密，加密密钥从 user_key 派生

**理由**:
- 用户只需保存 user_key，加密密钥通过 PBKDF2 派生，减少配置项
- 流式加密支持边打包压缩边加密，内存占用低
- GCM 模式提供认证加密，防止数据篡改
- Go 标准库 crypto/aes 和 crypto/cipher 原生支持

**密钥派生**:
```go
// 使用 PBKDF2 从 user_key 派生 32 字节加密密钥
func DeriveKey(userKey string, salt []byte) []byte {
    return pbkdf2.Key([]byte(userKey), salt, 100000, 32, sha256.New)
}
```

**替代方案考虑**:
- ChaCha20-Poly1305: 同样安全，但 AES 在 Go 中更常用
- 独立加密密钥: 需要用户保存两个密钥，体验较差

### 2. 压缩格式选择

**决策**: 使用 tar.gz 格式打包数据，配合流式处理

**理由**:
- Go 标准库 archive/tar 和 compress/gzip 原生支持
- 支持流式读写，可与加密管道串联
- 跨平台兼容性好，保留文件权限和目录结构

**流式处理流程**:

上传流程：
```
UserDataDir → tar → gzip → AES 加密 → HTTP 上传
```

下载流程：
```
HTTP 下载 → AES 解密 → gzip 解压 → tar 解包 → UserDataDir
```

### 3. 服务端数据存储

**决策**: 元数据存 SQLite，数据文件存文件系统

**理由**:
- SQLite 存储用户信息和版本元数据（版本号、hash、创建时间等）
- 文件系统存储加密后的数据文件，便于管理和备份
- 按用户 ID 分目录存储，避免单目录文件过多

**存储结构**:
```
data/
├── sync.db              # SQLite 数据库
└── files/
    ├── {sync_id}/
    │   ├── v1.dat       # 版本 1 数据文件
    │   ├── v2.dat       # 版本 2 数据文件
    │   └── ...
    └── ...
```

**表结构**:
```sql
-- users 表
CREATE TABLE users (
    id VARCHAR(32) PRIMARY KEY,  -- sync_id
    user_key VARCHAR(64) NOT NULL,
    last_sync_at DATETIME,       -- 上次同步时间（用于频率控制）
    created_at DATETIME,
    updated_at DATETIME
);

-- sync_versions 表
CREATE TABLE sync_versions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id VARCHAR(32),
    version_number INTEGER,
    file_hash VARCHAR(64),
    file_size BIGINT,
    file_path VARCHAR(255),      -- 数据文件路径
    created_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### 4. 认证机制

**决策**: 使用 sync_id + user_key 双因素认证

**理由**:
- sync_id 标识用户，user_key 验证身份并派生加密密钥
- 简单有效，无需 JWT 等复杂机制
- 通过 HTTPS 传输，安全性足够

**请求头设计**:
```
X-Sync-ID: <sync_id>
X-User-Key: <user_key>
```

### 5. 版本管理策略

**决策**: 保留最多5个版本，超出时删除最旧版本

**理由**:
- 5个版本足够应对误操作恢复
- 限制存储成本
- 简单实现，无需复杂版本合并

### 6. 同步频率控制

**决策**: 限制每个用户每分钟最多同步 1 次

**理由**:
- 防止恶意用户频繁上传导致资源耗尽
- 数据同步通常不需要高频操作
- 在 users 表记录 last_sync_at 实现简单控制

**实现方式**:
```go
// 检查上次同步时间
if time.Since(user.LastSyncAt) < time.Minute {
    return errors.New("sync too frequent, please wait")
}
```

### 7. 项目结构

**客户端** (`internal/sync/`):
```
internal/sync/
├── client.go      # HTTP 客户端，与服务器通信
├── crypto.go      # 流式加密解密功能，密钥派生
├── packer.go      # 数据流式打包/解包
├── config.go      # 同步配置结构
└── sync.go        # 主同步逻辑
```

**服务端** (`sync-backend/`):
```
sync-backend/
├── main.go        # 入口，CLI 命令处理
├── config.go      # 服务端配置
├── server.go      # HTTP 服务器
├── handler.go     # HTTP 处理器
├── model.go       # 数据模型
├── db.go          # 数据库连接
├── service.go     # 业务逻辑
└── rate.go        # 频率控制
```

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| user_key 泄露导致数据被解密 | 使用 HTTPS 传输；泄露后可重新生成凭证 |
| 全量同步导致网络开销大 | 数据量通常较小（配置+书签）；流式处理减少内存占用 |
| 服务端存储被暴力破解 | 使用强 user_key 派生加密密钥；服务端仅存储加密数据 |
| 版本冲突（多设备同时修改） | 以最新上传为准；后续可考虑添加冲突检测 |
| 流式处理中断导致数据损坏 | GCM 模式提供完整性校验；下载时验证 hash |
| 同步频率限制影响用户体验 | 限制较宽松（1分钟1次）；仅针对上传操作 |

## Migration Plan

1. **Phase 1**: 实现 sync-backend 服务端
   - 数据库模型和连接（SQLite）
   - 文件存储逻辑
   - HTTP API 实现
   - CLI 工具实现（生成 sync_id 和 user_key）
   - 频率控制实现

2. **Phase 2**: 实现客户端同步模块
   - 流式打包/压缩/加密功能
   - HTTP 客户端
   - 配置管理

3. **Phase 3**: 前端集成
   - 同步配置界面
   - 手动同步按钮
   - 版本恢复界面

4. **Phase 4**: 测试和文档
   - 集成测试
   - 部署文档

## Open Questions

无
