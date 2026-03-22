## 1. Database Service Changes

- [x] 1.1 Add Close method to internal/database/database.go (if not exists)
- [x] 1.2 Expose database close functionality through ConfigService

## 2. Sync Service Changes

- [x] 2.1 Modify Download flow in internal/sync/sync.go: close database after UnpackStream, before file replacement
- [x] 2.2 Add error handling for restore failure (prompt restart on error)

## 3. Frontend Changes

- [x] 3.1 Add restart prompt dialog state and handlers in SyncSettings.tsx
- [x] 3.2 Implement restart prompt dialog UI (OK button only, no automatic restart)
- [x] 3.3 Modify handleDownload to show restart prompt dialog after successful restore
- [x] 3.4 Disable sync operations after successful restore (app needs restart)

## 4. Testing & Verification

- [ ] 4.1 Test restore with database closed after extraction (file can be overwritten)
- [ ] 4.2 Test restore failure shows error and restart prompt
- [ ] 4.3 Test restart prompt appears after successful restore
- [ ] 4.4 Test sync operations are disabled after restore
