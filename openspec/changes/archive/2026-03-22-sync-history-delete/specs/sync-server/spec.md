## MODIFIED Requirements

### Requirement: Server stores up to 500 versions per user
The system SHALL retain up to 500 historical versions per user.

#### Scenario: Store new version within limit
- **WHEN** a user has fewer than 500 versions stored
- **AND** the user uploads new data
- **THEN** the system SHALL store the new version
- **AND** all existing versions SHALL be retained

#### Scenario: Cleanup old versions when exceeding limit
- **WHEN** a user has 500 versions stored
- **AND** the user uploads new data
- **THEN** the system SHALL store the new version
- **AND** the oldest version(s) SHALL be deleted to maintain the 500 version limit

#### Scenario: Version limit configuration
- **WHEN** the server configuration is loaded
- **THEN** the max_versions setting SHALL default to 500

## ADDED Requirements

### Requirement: Server supports permanently deleting a specific version
The system SHALL support permanently deleting a specific version by version number. The deletion is irreversible with no backup or recycle bin.

#### Scenario: Delete specific version permanently
- **WHEN** a delete request is received with a valid version number
- **THEN** the system SHALL permanently delete the version file from storage
- **AND** the system SHALL permanently delete the version record from database
- **AND** the system SHALL return success status
- **AND** the deleted version SHALL NOT be recoverable

#### Scenario: Delete non-existent version
- **WHEN** a delete request is received for a non-existent version
- **THEN** the system SHALL return an error indicating version not found

#### Scenario: Delete version with invalid credentials
- **WHEN** a delete request is received with invalid sync credentials
- **THEN** the system SHALL return an authentication error

### Requirement: Server supports paginated version listing
The system SHALL support paginated retrieval of versions with limit and offset parameters.

#### Scenario: List versions with pagination
- **WHEN** a list request is received with limit=50 and offset=0
- **THEN** the system SHALL return the first 50 versions ordered by version_number DESC

#### Scenario: List versions with offset
- **WHEN** a list request is received with limit=50 and offset=50
- **THEN** the system SHALL return versions 51-100 ordered by version_number DESC

#### Scenario: List versions without pagination parameters
- **WHEN** a list request is received without limit/offset parameters
- **THEN** the system SHALL return all versions for backward compatibility
