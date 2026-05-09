## Context

当前数据库初始化逻辑在 `internal/database/database.go` 中，采用硬编码方式：

- `createTables()` 创建所有表和索引（使用 `IF NOT EXISTS`）
- `migrateAddProxyJumpID()` 单独判断并添加 `proxy_jump_id` 列
- 版本号 `CurrentDBVersion = 2` 和 `DBVersion` 存储在 `config.toml` 中
- `CheckAndUpdateDBVersion()` 在 `services/service.go` 中决定是否跳过建表

这种设计的缺陷：

1. 版本号在配置文件中，与数据库分离，易不一致
2. 迁移逻辑分散，无法扩展
3. 无法记录迁移历史
4. 新安装和升级的逻辑耦合在 service 层

## Goals / Non-Goals

**Goals:**

- 建立数据库内部的迁移记录机制（`schema_migrations` 表）
- 所有迁移幂等性，可安全重复执行
- 移除配置文件中的数据库版本管理逻辑
- 现有用户数据不受影响，升级安全

**Non-Goals:**

- 不引入第三方迁移框架（如 golang-migrate）
- 不支持迁移回滚（down migration）
- 不改变现有表结构和业务数据模型

## Decisions

### 1. 版本号使用递增整数（非时间戳）

**决策**：继续使用整数版本号，从当前 `CurrentDBVersion=2` 延续。
**理由**：

- 项目当前是单人维护，整数版本号足够简单
- 与现有 `CurrentDBVersion` 概念兼容，迁移成本低
- 时间戳版本号的优势（避免冲突）在此场景不明显

### 2. 迁移记录表结构

```sql
CREATE TABLE IF NOT EXISTS schema_migrations (
    version INTEGER PRIMARY KEY,
    applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

**理由**：

- 仅记录版本号即可满足需求（用户选择）
- `applied_at` 用于审计，不增加复杂度
- 以版本号为唯一主键，天然去重

### 3. 迁移定义方式：函数数组

```go
type Migration struct {
    Version int
    Name    string
    Up      func(*sql.DB) error
}

var migrations = []Migration{
    {Version: 1, Name: "init schema", Up: migrateInitSchema},
    {Version: 2, Name: "add proxy_jump_id", Up: migrateAddProxyJumpID},
}
```

**理由**：

- 不引入外部依赖，保持项目轻量
- 函数内使用 `IF NOT EXISTS` 或列存在性检查实现幂等
- 易于添加新迁移，只需在数组末尾追加

### 4. 初始化流程重构

**新流程**：

```
打开数据库连接
  → 创建 schema_migrations 表
  → 遍历 migrations 数组
      → 检查版本是否已记录
      → 未记录则执行 Up 函数
      → 成功则插入版本记录
  → 初始化 Repository
```

**理由**：

- 统一处理首次安装和升级
- 不再需要 `skipCreateTables` 参数
- `services/service.go` 中移除 `CheckAndUpdateDBVersion` 调用

### 5. 幂等性实现策略

- **建表/建索引**：使用 `CREATE TABLE IF NOT EXISTS` / `CREATE INDEX IF NOT EXISTS`
- **加列**：SQLite 不支持 `ADD COLUMN IF NOT EXISTS`，需先查询 `pragma_table_info` 判断列是否存在
- **其他变更**：同理，先检查再执行

## Risks / Trade-offs

| Risk                                              | Mitigation                                                            |
| ------------------------------------------------- | --------------------------------------------------------------------- |
| 现有用户的 `config.toml` 中残留 `db_version` 字段 | 不影响功能，只是冗余数据；可选择清理但不必须                          |
| 迁移执行中途失败导致状态不一致                    | 每个迁移在独立事务中执行，失败不记录版本，下次启动会重试              |
| 未来可能需要回滚迁移                              | 当前设计不支持 down migration，如需回滚需手动写 SQL；项目当前无此需求 |
| 多实例同时启动竞争执行迁移                        | 桌面应用单实例运行，无竞争问题                                        |

## Migration Plan

1. **实现迁移系统**：
   - 在 `internal/database/` 下创建 `migration.go`
   - 定义 `Migration` 结构体、迁移数组、执行引擎
   - 将现有 `createTables` 和 `migrateAddProxyJumpID` 改造为幂等迁移函数

2. **重构 `database.go`**：
   - `Initialize()` 中调用迁移引擎
   - 移除 `skipCreateTables` 参数

3. **清理配置逻辑**：
   - `services/config_service.go`：移除 `CurrentDBVersion`、`DBVersion`、`CheckAndUpdateDBVersion`
   - `services/service.go`：移除 `skipCreateTables` 相关逻辑

4. **测试验证**：
   - 新安装：空数据库能正确执行所有迁移
   - 升级：已有数据库（版本1或2）能正确识别已执行迁移，不重复操作
   - 重复执行：所有迁移函数多次执行不报错

## Open Questions

- 是否需要清理现有 `config.toml` 中的 `db_version` 字段？（建议不清理，避免额外逻辑）
- 未来若迁移数量增多，是否需要将迁移函数分到单独文件？（当前数量少，可集中管理）
