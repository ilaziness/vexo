# sync-config

## Purpose

Sync configuration management, including sync ID, keys and encryption key configuration.

## Responsibilities

- Store and manage sync ID (unique user identifier)
- Store and manage user key (authentication)
- Store and manage encryption key (local data encryption)
- Configuration persistence via ConfigService
- Configuration validation

## Configuration Fields

| Field | Type | Description |
|-------|------|-------------|
| sync_id | string | Unique user identifier for sync service |
| user_key | string | Authentication key for sync server |
| encryption_key | string | Local encryption key for data security |

## Dependencies

- ConfigService (for persistence)

## Related Capabilities

- sync-client: Uses these credentials for sync operations
- sync-server: Generates credentials via CLI tool
