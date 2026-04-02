## ADDED Requirements

### Requirement: Encoder tool interface
The system SHALL provide an encoding/decoding tool.

#### Scenario: Display encoder tool
- **WHEN** user navigates to `/tools/encoder`
- **THEN** the system SHALL display the encoder tool interface

### Requirement: Encoding type selection
The system SHALL support multiple encoding/decoding operations.

#### Scenario: Select encoding type
- **WHEN** the encoder tool is displayed
- **THEN** the system SHALL provide options for: Base64, URL, HTML
- **AND** user SHALL be able to select one encoding type at a time

### Requirement: Text input
The system SHALL provide an input area for the source text.

#### Scenario: Enter source text
- **WHEN** the encoder tool is displayed
- **THEN** the system SHALL provide a multi-line text input field
- **AND** the input field SHALL have a placeholder indicating the expected input

### Requirement: Backend encode API
The system SHALL provide a backend API for encoding operations.

#### Scenario: Backend encode method
- **WHEN** the frontend calls the encode API with tool type and input
- **THEN** the backend SHALL encode the input using the specified encoding type
- **AND** the backend SHALL return the encoded result

### Requirement: Backend decode API
The system SHALL provide a backend API for decoding operations.

#### Scenario: Backend decode method
- **WHEN** the frontend calls the decode API with tool type and input
- **THEN** the backend SHALL decode the input using the specified encoding type
- **AND** the backend SHALL return the decoded result
- **AND** the backend SHALL return an error if the input is invalid

### Requirement: Encode operation
The system SHALL encode the input text via backend API.

#### Scenario: Click encode button
- **WHEN** user clicks the encode button
- **THEN** the frontend SHALL call the backend encode API
- **AND** the system SHALL display the encoded result

#### Scenario: Base64 encode
- **WHEN** user selects Base64 encoding and clicks encode
- **THEN** the backend SHALL encode the input using Base64 algorithm

#### Scenario: URL encode
- **WHEN** user selects URL encoding and clicks encode
- **THEN** the backend SHALL URL-encode the input text

### Requirement: Decode operation
The system SHALL decode the input text via backend API.

#### Scenario: Click decode button
- **WHEN** user clicks the decode button
- **THEN** the frontend SHALL call the backend decode API
- **AND** the system SHALL display the decoded result

### Requirement: Copy result
The system SHALL allow users to copy the result to clipboard.

#### Scenario: Click copy button
- **WHEN** user clicks the copy button
- **THEN** the system SHALL copy the result text to the clipboard
- **AND** the system SHALL provide visual feedback indicating success
