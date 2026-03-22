# sync-client

## Purpose

Client-side data sync functionality, including data packing, encryption, upload/download and version recovery.

## Responsibilities

- Data packing and compression (tar/gzip)
- Data encryption using user's local encryption key
- Communication with sync server (upload/download)
- Version recovery from server-stored versions
- Local data integrity verification

## Dependencies

- tar/gzip compression library
- Encryption library
- ConfigService (for sync credentials)
- HTTP client for server communication

## Related Capabilities

- sync-config: Provides sync ID, user key, encryption key
- sync-server: Remote endpoint for data storage

## Requirements

### Requirement: Data restore closes database after extraction and prompts for manual restart
The system SHALL close the database connection after extracting data and before replacing files, then prompt the user to manually restart the application.

#### Scenario: Database is closed after extraction
- **WHEN** the restore process downloads and extracts data to temp directory
- **THEN** the system SHALL close the database connection
- **AND** release all file handles to the database file

#### Scenario: Data files are replaced after database closure
- **WHEN** the database is closed
- **THEN** the system SHALL successfully overwrite the database file
- **AND** replace all data files in the user data directory

#### Scenario: User is prompted to manually restart after successful restore
- **WHEN** the data restore completes successfully
- **THEN** the system SHALL display a restart prompt dialog
- **AND** inform the user that the application needs to be manually restarted to use the restored data
- **AND** the dialog SHALL only have an OK button (no automatic restart)

#### Scenario: Restore fails after database closure
- **WHEN** the restore process fails after database is closed
- **THEN** the system SHALL display an error message
- **AND** inform the user to manually restart the application
- **AND** the application SHALL remain in a non-functional state until restart

### Requirement: Data restore uses temp directory in parent folder
The system SHALL extract downloaded backup data to a temporary directory located in `userDataDir` parent folder's `temp` subdirectory.

#### Scenario: Temp directory is created in parent/temp
- **WHEN** the restore process downloads backup data
- **THEN** the system SHALL create a temporary directory at `<parent>/temp/vexo-sync-tmp-<timestamp>`
- **AND** the downloaded data SHALL be extracted to this temporary directory

#### Scenario: Extracted data overwrites userDataDir contents
- **WHEN** the backup data is extracted to the temporary directory
- **THEN** the system SHALL copy the extracted files to overwrite contents in `userDataDir`
- **AND** the backup mechanism SHALL be used to ensure data safety
