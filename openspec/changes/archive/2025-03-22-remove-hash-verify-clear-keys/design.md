## Context

当前同步系统使用 AES-GCM 加密数据传输，同时在客户端和服务器端实现了额外的 SHA-256 hash 验证机制。AES-GCM 本身提供认证加密功能，能够检测数据篡改，因此独立的 hash 验证是冗余的。

此外，密钥派生使用 PBKDF2 生成加密密钥后，密钥在内存中留存直到垃圾回收，存在被内存转储攻击的风险。

## Goals / Non-Goals

**Goals:**
- 移除客户端和服务器端的冗余 hash 验证逻辑
- 在密钥派生后立即清除内存中的派生密钥
- 简化同步协议，减少网络传输开销
- 提高同步性能（减少 hash 计算时间）

**Non-Goals:**
- 不改变加密算法（仍使用 AES-GCM）
- 不改变密钥派生算法（仍使用 PBKDF2）
- 不修改用户认证机制
- 不添加新的功能特性

## Decisions

### 1. 移除 Hash 验证
**Decision**: 完全移除客户端和服务器端的 hash 计算和验证逻辑。

**Rationale**:
- AES-GCM 提供认证加密，能够检测密文篡改
- 独立的 hash 验证增加了计算开销和网络传输
- 简化代码逻辑，减少潜在 bug 点

**Alternatives Considered**:
- 保留 hash 作为可选功能：增加复杂性，无实际安全收益
- 使用更轻量的校验和（如 CRC32）：AES-GCM 已足够，无需额外校验

### 2. 内存密钥清除
**Decision**: 使用 `memset` 模式（通过重新赋值为零值）清除派生密钥。

**Rationale**:
- Go 语言没有内置的 `secureZeroMemory` 函数
- 通过将字节数组重新赋值为零值来实现清除
- 在 `NewEncryptStream` 和 `NewDecryptStream` 中密钥使用完毕后立即清除

**Implementation**:
```go
key := DeriveKey(userKey, salt)
// ... 使用 key ...
for i := range key {
    key[i] = 0
}
```

### 3. 数据库 Schema 变更
**Decision**: 直接删除 `file_hash` 数据库字段，不考虑兼容性。

**Rationale**:
- 不需要兼容旧版本，直接清理无用字段
- 简化数据库结构
- 减少存储开销

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| 旧版本客户端无法与新版服务器通信 | 这是一次 **BREAKING** 变更，需要同时升级客户端和服务器 |
| 内存清除可能被编译器优化跳过 | 使用 `for` 循环赋值，Go 编译器通常不会优化这种写法 |
| 失去 hash 用于调试的能力 | AES-GCM 解密失败会返回错误，可用于检测数据问题 |

## Migration Plan

1. **Phase 1**: 部署新版服务器（数据库字段已删除）
2. **Phase 2**: 发布新版客户端

**注意**: 不兼容旧版本，需要同时升级客户端和服务器。

### 前端变更
- 移除 `VersionInfo` 接口中的 `file_hash` 字段
- 移除版本列表中显示的 hash 标签
- 移除 `VerifyDataIntegrity` 调用和相关 UI
- 移除 `calculating_hash` 阶段标签
