## ADDED Requirements

### Requirement: Regex tool interface
The system SHALL provide a regular expression testing tool.

#### Scenario: Display regex tool
- **WHEN** user navigates to `/tools/regex`
- **THEN** the system SHALL display the regex testing tool interface

### Requirement: Pattern input
The system SHALL provide an input field for the regular expression pattern.

#### Scenario: Enter regex pattern
- **WHEN** the regex tool is displayed
- **THEN** the system SHALL provide an input field for the regex pattern

### Requirement: Test text input
The system SHALL provide an input area for the test text.

#### Scenario: Enter test text
- **WHEN** the regex tool is displayed
- **THEN** the system SHALL provide a multi-line text area for test input
- **AND** the system SHALL support entering multiple lines of text

### Requirement: Regex flags
The system SHALL support common regex flags.

#### Scenario: Toggle regex flags
- **WHEN** the regex tool is displayed
- **THEN** the system SHALL provide toggles for: case insensitive, multiline, global
- **AND** the flags SHALL be applied to the regex matching

### Requirement: Backend regex match API
The system SHALL provide a backend API for regex matching.

#### Scenario: Backend regex match method
- **WHEN** the frontend calls the regex match API with pattern, text, and flags
- **THEN** the backend SHALL compile the regex pattern with the specified flags
- **AND** the backend SHALL perform matching against the provided text
- **AND** the backend SHALL return all matches including positions and capture groups
- **AND** the backend SHALL return an error if the pattern is invalid

### Requirement: Real-time matching
The system SHALL perform regex matching via backend API.

#### Scenario: Display matches
- **WHEN** user enters a regex pattern and test text
- **THEN** the frontend SHALL call the backend regex match API
- **AND** the system SHALL highlight all matches in the test text
- **AND** the system SHALL display the number of matches found

#### Scenario: Display match groups
- **WHEN** the regex pattern contains capture groups
- **THEN** the system SHALL display each match with its captured groups
- **AND** the system SHALL show the group index and matched content

### Requirement: Invalid pattern handling
The system SHALL handle invalid regex patterns gracefully.

#### Scenario: Invalid regex pattern
- **WHEN** user enters an invalid regex pattern
- **THEN** the backend SHALL return an error
- **AND** the frontend SHALL display an error message indicating the syntax error
