## ADDED Requirements

### Requirement: Header displays backup upload button
The Header component SHALL display a backup upload button positioned before the settings button.

#### Scenario: User views header
- **WHEN** the Header component is rendered
- **THEN** a backup upload button with CloudUpload icon is visible
- **AND** the button is positioned between the command panel button and the settings button

### Requirement: Backup button shows confirmation dialog
The system SHALL display a confirmation dialog when the user clicks the backup upload button.

#### Scenario: User clicks backup button
- **WHEN** the user clicks the backup upload button
- **THEN** a confirmation dialog appears asking "确定要上传备份吗？"
- **AND** the dialog provides "取消" and "确定" options

### Requirement: Confirmation dialog can be dismissed
The user SHALL be able to cancel the backup operation from the confirmation dialog.

#### Scenario: User cancels backup
- **GIVEN** the confirmation dialog is open
- **WHEN** the user clicks "取消"
- **THEN** the dialog closes
- **AND** no backup operation is initiated

### Requirement: Backup upload executes on confirmation
The system SHALL execute the backup upload when the user confirms in the dialog.

#### Scenario: User confirms backup
- **GIVEN** the confirmation dialog is open
- **WHEN** the user clicks "确定"
- **THEN** the dialog closes
- **AND** the backup upload process begins
- **AND** a full-screen loading overlay is displayed

### Requirement: Full-screen loading during backup
The system SHALL display a full-screen loading overlay during the backup upload process.

#### Scenario: Backup is in progress
- **GIVEN** the backup upload has been initiated
- **WHEN** the upload is in progress
- **THEN** a full-screen loading overlay covers the entire application
- **AND** the loading message indicates "正在上传备份..."

### Requirement: Loading overlay hides on completion
The system SHALL hide the loading overlay when the backup upload completes.

#### Scenario: Backup completes successfully
- **GIVEN** the backup upload is in progress
- **WHEN** the upload completes successfully
- **THEN** the loading overlay is hidden
- **AND** a success message is displayed

#### Scenario: Backup fails
- **GIVEN** the backup upload is in progress
- **WHEN** the upload fails
- **THEN** the loading overlay is hidden
- **AND** an error message is displayed with the failure reason

### Requirement: Check sync configuration before backup
The system SHALL verify that sync configuration is set before initiating backup.

#### Scenario: Sync not configured
- **GIVEN** the sync configuration is not complete (missing serverUrl, syncId, or userKey)
- **WHEN** the user clicks the backup upload button
- **THEN** an error message is displayed: "请先前往设置页面配置同步参数"
- **AND** no confirmation dialog is shown
