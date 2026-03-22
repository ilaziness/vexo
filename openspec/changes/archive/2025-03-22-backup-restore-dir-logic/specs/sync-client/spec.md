## MODIFIED Requirements

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
