## ADDED Requirements

### Requirement: User can delete a specific version from history
The system SHALL allow users to delete a specific version from their sync history through the UI.

#### Scenario: Delete button visible for each version
- **WHEN** the user views the sync history list
- **THEN** each version entry SHALL display a delete button

#### Scenario: Confirm before delete
- **WHEN** the user clicks the delete button on a version
- **THEN** a confirmation dialog SHALL appear asking for confirmation
- **AND** the delete SHALL NOT execute until confirmed

#### Scenario: Cancel delete operation
- **WHEN** the user clicks the delete button
- **AND** the confirmation dialog appears
- **AND** the user clicks cancel
- **THEN** the dialog SHALL close
- **AND** no delete operation SHALL be performed

#### Scenario: Successful version deletion
- **WHEN** the user confirms the delete operation
- **THEN** the system SHALL send a delete request to the server
- **AND** upon success, the version SHALL be removed from the list
- **AND** a success message SHALL be displayed

#### Scenario: Failed version deletion
- **WHEN** the delete request fails
- **THEN** an error message SHALL be displayed
- **AND** the version SHALL remain in the list

## MODIFIED Requirements

### Requirement: Sync settings page shows limited version preview
The system SHALL display only the first 10 versions in the sync settings page history list.

#### Scenario: Display limited version list in settings
- **WHEN** the sync settings page loads
- **THEN** the system SHALL fetch and display up to 10 most recent versions
- **AND** each version SHALL show: version number, file size, creation date
- **AND** each version SHALL have a delete button

### Requirement: Restore dialog supports infinite scroll for versions
The system SHALL support loading versions in batches with infinite scroll in the restore data dialog.

#### Scenario: Initial load of versions in restore dialog
- **WHEN** the restore data dialog opens
- **THEN** the system SHALL load the first 50 versions
- **AND** display them in a scrollable list

#### Scenario: Scroll to load more versions
- **WHEN** the user scrolls to the bottom of the version list
- **AND** more versions are available
- **THEN** the system SHALL load the next 50 versions
- **AND** append them to the existing list

#### Scenario: All versions loaded
- **WHEN** the user scrolls to the bottom
- **AND** no more versions are available
- **THEN** the system SHALL indicate that all versions have been loaded
