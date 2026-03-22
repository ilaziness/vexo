## Why

当前同步功能的恢复目录逻辑需要优化。恢复时下载和解压后的临时目录应放在 `userDataDir` 上一级的 `temp` 目录中，避免在数据目录所在目录直接创建临时目录造成混乱。

## What Changes

- **恢复逻辑修改**: 
  - 下载和解压后的临时目录改为放在 `userDataDir` 上一级的 `temp` 目录
  - 解压后的文件内容覆盖恢复 `userDataDir`
- **目录结构变更**: 临时目录位置从 `filepath.Dir(dstDir)` 改为 `filepath.Dir(dstDir)/temp`

## Capabilities

### New Capabilities
- 无

### Modified Capabilities
- `sync-client`: 修改数据恢复的目录处理逻辑
  - `Download` 方法: 临时目录改为 `userDataDir` 上一级的 `temp` 目录，解压后覆盖恢复

## Impact

- `internal/sync/sync.go`: `Download` 方法的临时目录创建逻辑
- 保持打包逻辑不变，仍然打包 `userDataDir` 目录内的文件
