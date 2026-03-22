## Why

当前同步功能在客户端和服务器端都实现了独立的 hash 验证逻辑，但数据已经使用 AES-GCM 加密，该算法本身提供了认证加密（authenticated encryption）功能，能够检测数据篡改。此外，密钥派生后生成的加密密钥在内存中留存时间过长，存在安全风险。

## What Changes

- **移除客户端 hash 计算逻辑**：删除 `internal/sync/sync.go` 中的 `calculateDirHash` 和 `simplePack` 方法
- **移除服务器端 hash 验证逻辑**：删除 `sync-backend/handler.go` 中的 `X-File-Hash` 请求头验证
- **移除 hash 相关数据结构**：从 `VersionInfo` 结构体中删除 `FileHash` 字段
- **添加派生密钥内存清除**：在 `internal/sync/crypto.go` 的 `DeriveKey` 调用后立即清除密钥内存
- **删除前端 hash 相关代码**：移除 `SyncSettings.tsx` 中的 `file_hash` 字段和完整性校验功能
- **删除服务端完整性校验**：移除 `services/sync_service.go` 中的 `VerifyDataIntegrity` 和 `calculateDirHash` 方法
- **BREAKING**: 同步协议变更 - 上传/下载不再传递 hash，数据库 `file_hash` 字段直接删除

## Capabilities

### New Capabilities
- `sync-security-enhancement`: 移除冗余的 hash 验证，增强密钥派生安全性

### Modified Capabilities
- `sync-protocol`: 移除 hash 验证相关的请求头和响应头，简化同步协议

## Impact

- `internal/sync/crypto.go`: 修改 `NewEncryptStream` 和 `NewDecryptStream` 以清除派生密钥
- `internal/sync/sync.go`: 移除 hash 计算和验证逻辑
- `internal/sync/client.go`: 移除 `FileHash` 字段和 hash 相关请求头
- `sync-backend/handler.go`: 移除 `X-File-Hash` 请求头处理
- `sync-backend/service.go`: 移除 `FileHash` 参数和存储
- `sync-backend/model.go`: 从数据库模型中删除 `FileHash` 字段
- `services/sync_service.go`: 删除 `VerifyDataIntegrity` 和 `calculateDirHash` 方法
- `frontend/src/components/SyncSettings.tsx`: 移除 `file_hash` 字段和完整性校验 UI
- **Database Migration**: 删除 `sync_versions` 表的 `file_hash` 字段
