## Why

当前数据库迁移功能简陋：版本号记录在配置文件(config.toml)中，迁移逻辑硬编码在database.go里，缺乏系统化的迁移管理。这导致无法追踪迁移历史、无法安全地重复执行迁移、且配置文件与数据库状态可能不一致。需要建立一个真正的数据库迁移系统，将版本记录迁移到数据库内部，并支持幂等性迁移以确保现有安装能安全升级。

## What Changes

- 在数据库中创建迁移记录表 `schema_migrations`，用于记录已执行的迁移版本
- 使用整数版本号（从当前 `CurrentDBVersion=2` 继续递增），每个迁移有唯一的版本号
- 所有迁移改造为幂等性（可重复执行），包括：
  - 初始化数据库表结构（`CREATE TABLE IF NOT EXISTS`、`CREATE INDEX IF NOT EXISTS`）
  - 添加 `proxy_jump_id` 列（`ADD COLUMN IF NOT EXISTS` 或等效判断）
- 移除 `config.toml` 中的 `db_version` 字段及相关逻辑（`CheckAndUpdateDBVersion`、`CurrentDBVersion` 等）
- 数据库初始化时自动检测并执行所有未执行的迁移
- 保持现有应用安装升级不被破坏：新系统能正确识别已有数据库状态，不会重复执行破坏性操作

## Capabilities

### New Capabilities
- `database-migration`: 数据库迁移管理，包括迁移记录表、迁移执行引擎、幂等性迁移定义

### Modified Capabilities
- （无现有spec需要修改）

## Impact

- `internal/database/database.go`：重构初始化逻辑，引入迁移系统
- `services/config_service.go`：移除 `CurrentDBVersion`、`DBVersion`、`CheckAndUpdateDBVersion` 及相关逻辑
- `services/service.go`：移除 `skipCreateTables` 参数传递，简化数据库初始化调用
- 现有用户数据：不受影响，迁移系统会自动适配已有数据库状态
