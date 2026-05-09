## 1. Create migration system core

- [x] 1.1 Create `internal/database/migration.go` with `Migration` struct, migration registry, and execution engine
- [x] 1.2 Implement `createSchemaMigrationsTable()` to create the migration tracking table
- [x] 1.3 Implement `getAppliedMigrations()` to query already executed versions
- [x] 1.4 Implement `runMigrations()` to execute pending migrations in order and record versions

## 2. Define idempotent migrations

- [x] 2.1 Extract `createTables()` logic into `migrateInitSchema()` migration (version 1) with `IF NOT EXISTS`
- [x] 2.2 Extract `migrateAddProxyJumpID()` into formal migration (version 2) with column existence check
- [x] 2.3 Register both migrations in the migration registry array

## 3. Refactor database initialization

- [x] 3.1 Update `Database.Initialize()` to remove `skipCreateTables` parameter
- [x] 3.2 Update `Database.Initialize()` to call migration engine instead of direct `createTables()`
- [x] 3.3 Remove old `createTables()` and `migrateAddProxyJumpID()` methods from `database.go` (or keep as migration functions)

## 4. Clean up configuration-based version tracking

- [x] 4.1 Remove `CurrentDBVersion` constant from `services/config_service.go`
- [x] 4.2 Remove `DBVersion` field from `AppConfig` struct
- [x] 4.3 Remove `CheckAndUpdateDBVersion()` method from `services/config_service.go`
- [x] 4.4 Update `services/service.go` to remove `skipCreateTables` logic and simplify `db.Initialize()` call

## 5. Verify backward compatibility

- [x] 5.1 Test fresh install: empty database correctly executes all migrations
- [x] 5.2 Test upgrade: existing database with tables and data skips executed migrations safely
- [x] 5.3 Test idempotency: run migrations multiple times without errors
- [x] 5.4 Run `go build` to ensure compilation succeeds
