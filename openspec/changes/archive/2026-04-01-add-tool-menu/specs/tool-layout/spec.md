## ADDED Requirements

### Requirement: Tool interface layout
The system SHALL use a consistent layout for all tool interfaces.

#### Scenario: Tool page layout structure
- **WHEN** a specific tool page is displayed
- **THEN** the system SHALL show a left sidebar with the tool list
- **AND** the system SHALL show a main content area on the right
- **AND** the left sidebar width SHALL be 200 pixels

### Requirement: Tool list navigation
The system SHALL display all available tools in the left sidebar for quick switching.

#### Scenario: Tool list display
- **WHEN** any tool page is displayed
- **THEN** the left sidebar SHALL list all available tools
- **AND** the currently active tool SHALL be visually highlighted

#### Scenario: Switch tool via sidebar
- **WHEN** user clicks on a tool in the left sidebar
- **THEN** the system SHALL navigate to the selected tool
- **AND** the main content area SHALL update to show the selected tool

### Requirement: Back to dashboard
The system SHALL provide a way to return to the tool dashboard.

#### Scenario: Click back button
- **WHEN** user clicks the back/home button
- **THEN** the system SHALL navigate back to `/tools` (the dashboard)
