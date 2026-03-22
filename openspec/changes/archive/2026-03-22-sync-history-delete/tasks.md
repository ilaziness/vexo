## 1. Server-side Changes

- [x] 1.1 Update default max_versions from 5 to 500 in sync-backend/config.go
- [x] 1.2 Add DeleteVersion method to sync-backend/service.go (delete DB record and physical file permanently)
- [x] 1.3 Add DeleteVersionHandler to sync-backend/handler.go
- [x] 1.4 Register DELETE /versions endpoint in handler Routes
- [x] 1.5 Add pagination support (limit/offset) to ListVersions in sync-backend/service.go
- [x] 1.6 Add pagination parameters to ListVersionsHandler in sync-backend/handler.go
- [x] 1.7 Add DeleteVersion method to internal/sync/client.go
- [x] 1.8 Add DeleteSyncVersion method to services/sync_service.go
- [x] 1.9 Add pagination support to ListSyncVersions in services/sync_service.go

## 2. Client-side UI Changes - Sync Settings Page

- [x] 2.1 Limit version list to display only first 10 versions in SyncSettings.tsx
- [x] 2.2 Add delete icon button to each version item in SyncSettings.tsx
- [x] 2.3 Add delete confirmation dialog state and handlers
- [x] 2.4 Implement delete confirmation dialog UI
- [x] 2.5 Add delete version API call with error handling
- [x] 2.6 Refresh version list after successful deletion
- [x] 2.7 Add success/error message notifications for delete operation

## 3. Client-side UI Changes - Restore Dialog

- [x] 3.1 Implement infinite scroll for version list in restore dialog
- [x] 3.2 Add pagination state (page number, hasMore) for version loading
- [x] 3.3 Modify version loading to fetch 50 items per request
- [x] 3.4 Add scroll event handler to trigger next page load
- [x] 3.5 Add loading indicator when fetching more versions
- [x] 3.6 Handle all versions loaded state

## 4. Testing & Verification

- [ ] 4.1 Test server version limit of 500 versions
- [ ] 4.2 Test delete API endpoint with valid credentials
- [ ] 4.3 Test delete API with invalid credentials (should fail)
- [ ] 4.4 Test delete confirmation dialog appears
- [ ] 4.5 Test cancel delete operation
- [ ] 4.6 Test successful delete removes item from list
- [ ] 4.7 Test error handling when delete fails
- [ ] 4.8 Test sync settings page shows only 10 versions
- [ ] 4.9 Test restore dialog loads 50 versions initially
- [ ] 4.10 Test restore dialog infinite scroll loads more versions
- [ ] 4.11 Test restore dialog shows all loaded indicator at end
