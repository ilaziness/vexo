## 1. Modify Restore Logic

- [x] 1.1 Update `Download` method in `internal/sync/sync.go` to create temp directory in `parent/temp/`
- [x] 1.2 Update restore logic to copy extracted files to overwrite `userDataDir` contents
- [x] 1.3 Ensure backup/restore safety mechanism still works

## 2. Testing

- [x] 2.1 Test restore extracts to correct temp location (`parent/temp/`)
- [x] 2.2 Test restore overwrites `userDataDir` contents correctly
- [x] 2.3 Test failure recovery (backup restore) still works

## 3. Verification

- [x] 3.1 Run all sync-related unit tests
- [x] 3.2 Manual end-to-end test of upload and restore flow
