## 1. Import Dependencies

- [x] 1.1 Import CloudUpload icon from @mui/icons-material
- [x] 1.2 Import Dialog, DialogTitle, DialogContent, DialogActions, Button, Backdrop from @mui/material
- [x] 1.3 Import UploadSync, GetSyncProgress from syncservice bindings
- [x] 1.4 Import SyncProgress, LogService from services bindings
- [x] 1.5 Import Loading component from ./Loading

## 2. Add State Management

- [x] 2.1 Add `confirmDialogOpen` state for confirmation dialog visibility
- [x] 2.2 Add `uploading` state for tracking upload progress
- [x] 2.3 Add `progress` state for SyncProgress data

## 3. Add Backup Button to Header

- [x] 3.1 Add backup button with CloudUpload icon between command panel and settings buttons
- [x] 3.2 Add Tooltip with title "上传备份"
- [x] 3.3 Add onClick handler to open confirmation dialog

## 4. Implement Confirmation Dialog

- [x] 4.1 Create Dialog component with DialogTitle "确认上传"
- [x] 4.2 Add DialogContent with message "确定要上传备份吗？"
- [x] 4.3 Add DialogActions with "取消" and "确定" buttons
- [x] 4.4 Implement cancel handler to close dialog without action
- [x] 4.5 Implement confirm handler to start upload and close dialog

## 5. Implement Upload Logic

- [x] 5.1 Create `handleUpload` function that calls UploadSync
- [x] 5.2 Set uploading state to true at start
- [x] 5.3 Implement progress polling using GetSyncProgress (500ms interval)
- [x] 5.4 Stop polling when progress.isCompleted or progress.error
- [x] 5.5 Set uploading state to false on completion/error
- [x] 5.6 Show success message on completion
- [x] 5.7 Show error message on failure using parseCallServiceError

## 6. Add Full-screen Loading Overlay

- [x] 6.1 Add Backdrop component with open prop bound to uploading state
- [x] 6.2 Set Backdrop zIndex to theme.zIndex.modal + 1
- [x] 6.3 Embed Loading component inside Backdrop with message "正在上传备份..."

## 7. Add Sync Configuration Check

- [x] 7.1 Import ReadConfig from configservice bindings
- [x] 7.2 Check if sync config (serverUrl, syncId, userKey) is complete before showing dialog
- [x] 7.3 Show error message if config is incomplete: "请先前往设置页面配置同步参数"

## 8. Testing

- [x] 8.1 Verify backup button appears in correct position
- [x] 8.2 Verify confirmation dialog opens on click
- [x] 8.3 Verify cancel closes dialog without upload
- [x] 8.4 Verify upload starts and loading shows on confirm
- [x] 8.5 Verify loading hides and success message shows on completion
- [x] 8.6 Verify error handling works correctly
