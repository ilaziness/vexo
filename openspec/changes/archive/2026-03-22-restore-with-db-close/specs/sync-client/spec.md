## MODIFIED Requirements

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
