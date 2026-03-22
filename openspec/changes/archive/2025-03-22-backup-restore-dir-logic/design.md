## Context

当前数据同步功能的恢复目录处理逻辑需要优化：

1. **打包逻辑**: 保持现有逻辑，打包 `userDataDir` 目录内的内容
2. **恢复逻辑**: 当前临时目录创建在 `filepath.Dir(dstDir)`，需要改为放在 `userDataDir` 上一级的 `temp` 目录中

相关代码位于:
- `internal/sync/packer.go`: 打包和解包功能
- `internal/sync/sync.go`: 同步管理器的上传下载逻辑

## Goals / Non-Goals

**Goals:**
- 保持打包逻辑不变，仍然打包 `userDataDir` 目录内的文件
- 恢复时临时目录放在 `userDataDir` 上一级目录的 `temp` 文件夹中
- 解压后的文件内容覆盖恢复 `userDataDir`
- 保持现有的数据库关闭和重启提示机制

**Non-Goals:**
- 不改变打包逻辑
- 不改变加密/压缩算法
- 不改变与同步服务器的通信协议
- 不改变版本管理和历史记录功能

## Decisions

### 1. 临时目录位置调整

**决策**: 临时目录改为 `filepath.Join(filepath.Dir(dstDir), "temp", "vexo-sync-tmp-<timestamp>")`

**原因**:
- 将临时文件与数据目录分离，避免污染
- 使用专门的 `temp` 目录便于管理和清理
- 符合常见的临时文件存放规范

### 2. 恢复覆盖策略

**决策**: 解压后的文件内容覆盖恢复 `userDataDir`

**原因**:
- 保持现有的备份/恢复机制作为安全网
- 解压后的文件直接覆盖到 `userDataDir` 目录内

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| 临时目录位置改变可能影响权限 | 确保 `temp` 目录有正确的读写权限 |
| 覆盖恢复失败可能导致数据丢失 | 保留现有的备份机制，失败时可以恢复 |
