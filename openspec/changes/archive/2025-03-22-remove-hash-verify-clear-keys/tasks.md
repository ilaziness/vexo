## 1. Client-side Hash Removal

- [ ] 1.1 Remove `FileHash` field from `VersionInfo` struct in `internal/sync/client.go`
- [ ] 1.2 Remove `X-File-Hash` header from `Upload` method in `internal/sync/client.go`
- [ ] 1.3 Remove file hash return value from `Download` method in `internal/sync/client.go`
- [ ] 1.4 Remove `calculateDirHash` method from `internal/sync/sync.go`
- [ ] 1.5 Remove `packForHash` method from `internal/sync/sync.go`
- [ ] 1.6 Remove `simplePack` method from `internal/sync/sync.go`
- [ ] 1.7 Remove hash verification logic from `Download` method in `internal/sync/sync.go`

## 2. Server-side Hash Removal

- [ ] 2.1 Remove `X-File-Hash` header validation from `UploadHandler` in `sync-backend/handler.go`
- [ ] 2.2 Remove `X-File-Hash` header from `DownloadHandler` response in `sync-backend/handler.go`
- [ ] 2.3 Update `SaveVersion` method signature in `sync-backend/service.go` to remove `fileHash` parameter
- [ ] 2.4 Remove `FileHash` field from `SyncVersion` model in `sync-backend/model.go`
- [ ] 2.5 Remove `FileHash` field from `VersionInfo` struct in `sync-backend/service.go`

## 3. Service Layer Hash Removal

- [ ] 3.1 Remove `VerifyDataIntegrity` method from `services/sync_service.go`
- [ ] 3.2 Remove `calculateDirHash` method from `services/sync_service.go`
- [ ] 3.3 Remove `calculating_hash` stage handling from progress reporting

## 4. Frontend Hash Removal

- [ ] 4.1 Remove `file_hash` field from `VersionInfo` interface in `frontend/src/components/SyncSettings.tsx`
- [ ] 4.2 Remove hash display Chip from version list in `SyncSettings.tsx`
- [ ] 4.3 Remove `VerifyDataIntegrity` import and call from `SyncSettings.tsx`
- [ ] 4.4 Remove `integrityResult` state and related UI from `SyncSettings.tsx`
- [ ] 4.5 Remove `checkIntegrity` function from `SyncSettings.tsx`
- [ ] 4.6 Remove `calculating_hash` from `stageLabels` in `SyncSettings.tsx`
- [ ] 4.7 Remove `VerifyDataIntegrity` binding function from sync service

## 3. Crypto Security Enhancement

- [ ] 3.1 Add `clearKey` helper function in `internal/sync/crypto.go` to zero out key memory
- [ ] 3.2 Update `NewEncryptStream` to clear derived key after creating AES-GCM cipher
- [ ] 3.3 Update `NewDecryptStream` to clear derived key after creating AES-GCM cipher

## 4. Testing

- [ ] 4.1 Update `internal/sync/crypto_test.go` to verify key clearing behavior
- [ ] 4.2 Update `internal/sync/packer_test.go` to remove hash-related tests
- [ ] 4.3 Update `internal/sync/config_test.go` if affected by changes
- [ ] 4.4 Run all sync-related tests to ensure no regressions

## 5. Cleanup

- [ ] 5.1 Remove unused `CalculateHash` function from `internal/sync/crypto.go` if no longer needed
- [ ] 5.2 Remove unused `CalculateFileHash` function from `internal/sync/crypto.go` if no longer needed
- [ ] 5.3 Remove unused `CalculateFileHashForDir` function from `internal/sync/crypto.go` if no longer needed
- [ ] 5.4 Update imports in affected files to remove unused dependencies
