## REMOVED Requirements

### Requirement: File hash verification during upload
The system SHALL verify the file hash provided by the client during upload.

**Reason**: AES-GCM encryption provides authenticated encryption, making separate hash verification redundant.

**Migration**: Remove X-File-Hash header from upload requests. Server will no longer validate or store file hashes. Database file_hash column will be dropped.

#### Scenario: Upload without hash header
- **WHEN** client uploads data without X-File-Hash header
- **THEN** server accepts the upload and stores the version

### Requirement: Client-side data integrity verification
The system SHALL provide a way for clients to verify local data integrity against server data.

**Reason**: Hash-based integrity verification is removed. AES-GCM authentication is sufficient.

**Migration**: Remove VerifyDataIntegrity function and related UI.

#### Scenario: No integrity check after sync
- **WHEN** upload or download completes
- **THEN** no hash-based integrity verification is performed

## MODIFIED Requirements

### Requirement: Upload data to server
The system SHALL accept encrypted data uploads from authenticated clients.

#### Scenario: Successful upload
- **WHEN** an authenticated client sends a POST request to /upload with encrypted data
- **THEN** the server stores the data as a new version
- **AND** returns the version number and creation timestamp
- **AND** the response does not include file_hash field

### Requirement: Download data from server
The system SHALL allow authenticated clients to download encrypted data.

#### Scenario: Successful download
- **WHEN** an authenticated client sends a GET request to /download
- **THEN** the server returns the encrypted data for the requested version
- **AND** the response does not include X-File-Hash header

### Requirement: Version information display
The system SHALL display version history to users.

#### Scenario: Display version list
- **WHEN** user views sync settings
- **THEN** the version list shows version number, file size, and creation time
- **AND** the version list does not show file hash

## ADDED Requirements

### Requirement: Secure key derivation
The system SHALL clear derived encryption keys from memory immediately after use.

#### Scenario: Key derivation during encryption
- **WHEN** the system derives an encryption key from user key
- **THEN** the derived key SHALL be cleared from memory after creating the AES-GCM cipher

#### Scenario: Key derivation during decryption
- **WHEN** the system derives a decryption key from user key
- **THEN** the derived key SHALL be cleared from memory after creating the AES-GCM cipher
