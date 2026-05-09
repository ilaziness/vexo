## ADDED Requirements

### Requirement: Migration record table exists

The system SHALL create a `schema_migrations` table in the database to track executed migrations.

#### Scenario: New database

- **WHEN** the application starts with a new database
- **THEN** the `schema_migrations` table is created with columns `version` (INTEGER PRIMARY KEY) and `applied_at` (DATETIME)

#### Scenario: Existing database

- **WHEN** the application starts with an existing database
- **THEN** the `schema_migrations` table is created if it does not already exist

### Requirement: Migrations are idempotent

The system SHALL ensure all migration functions can be safely executed multiple times without causing errors or data corruption.

#### Scenario: Repeated table creation

- **WHEN** a migration that creates tables is executed on a database where those tables already exist
- **THEN** the migration completes successfully without error and existing data is preserved

#### Scenario: Repeated column addition

- **WHEN** a migration that adds a column is executed on a database where that column already exists
- **THEN** the migration completes successfully without error

### Requirement: Migration execution is tracked

The system SHALL record the version number in `schema_migrations` after a migration is successfully executed.

#### Scenario: Successful migration

- **WHEN** a migration is executed successfully
- **THEN** its version number is inserted into `schema_migrations`

#### Scenario: Skipping executed migrations

- **WHEN** the application starts and a migration version already exists in `schema_migrations`
- **THEN** that migration is skipped and not executed again

### Requirement: Migrations run in order

The system SHALL execute migrations in ascending version order.

#### Scenario: Multiple pending migrations

- **WHEN** the application starts and multiple migrations have not been executed
- **THEN** they are executed sequentially from lowest version to highest version

### Requirement: Configuration-based version tracking is removed

The system SHALL no longer use `config.toml` to track database version.

#### Scenario: Application startup

- **WHEN** the application starts
- **THEN** it does not read or write `db_version` from any configuration file
- **AND** database migration state is determined solely from `schema_migrations`

### Requirement: Backward compatibility for existing installations

The system SHALL safely handle upgrades from the current version-tracking system without data loss.

#### Scenario: Upgrade from config-based version

- **WHEN** an existing user upgrades the application
- **THEN** the database initialization detects the actual database state via `schema_migrations`
- **AND** all necessary migrations are applied correctly
- **AND** existing user data remains intact
