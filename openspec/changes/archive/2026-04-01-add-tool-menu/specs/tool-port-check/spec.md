## ADDED Requirements

### Requirement: Port check tool interface
The system SHALL provide a port connectivity testing tool.

#### Scenario: Display port check tool
- **WHEN** user navigates to `/tools/port-check`
- **THEN** the system SHALL display the port check tool interface

### Requirement: Host input
The system SHALL allow users to enter a target host address.

#### Scenario: Enter host address
- **WHEN** the port check tool is displayed
- **THEN** the system SHALL provide an input field for the host address
- **AND** the input SHALL accept IP addresses or hostnames

### Requirement: Port input
The system SHALL allow users to specify the port number to test.

#### Scenario: Enter port number
- **WHEN** the port check tool is displayed
- **THEN** the system SHALL provide an input field for the port number
- **AND** the input SHALL accept numeric values between 1 and 65535

### Requirement: Backend port check API
The system SHALL provide a backend API for port connectivity checking.

#### Scenario: Backend port check method
- **WHEN** the frontend calls the port check API with host and port
- **THEN** the backend SHALL attempt a TCP connection to the specified host and port
- **AND** the backend SHALL measure the connection time
- **AND** the backend SHALL return the result including success status and response time

### Requirement: Execute port check
The system SHALL test the connectivity to the specified host and port via backend.

#### Scenario: Click check button
- **WHEN** user clicks the check button
- **THEN** the frontend SHALL call the backend port check API
- **AND** the system SHALL display the result (success or failure)
- **AND** the system SHALL display the response time if successful

#### Scenario: Port is open
- **WHEN** the port check is executed and the port is accessible
- **THEN** the system SHALL display a success indicator
- **AND** the system SHALL show the connection time

#### Scenario: Port is closed or unreachable
- **WHEN** the port check is executed and the port is not accessible
- **THEN** the system SHALL display a failure indicator
- **AND** the system SHALL show an appropriate error message
