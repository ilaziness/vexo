## ADDED Requirements

### Requirement: Tool menu button in Header
The system SHALL display a tool menu button in the Header component.

#### Scenario: Tool button visible
- **WHEN** the main application window is displayed
- **THEN** the Header component SHALL show a tool button with an appropriate icon

### Requirement: Tool window management
The system SHALL open a new window when the tool button is clicked.

#### Scenario: Open tool window
- **WHEN** user clicks the tool menu button
- **THEN** the system SHALL open a tool window loading the `/tools` route
- **AND** the window SHALL have dimensions of 1200x800 pixels
- **AND** the window SHALL be frameless

#### Scenario: Tool window singleton
- **WHEN** user clicks the tool menu button and a tool window is already open
- **THEN** the system SHALL focus the existing tool window instead of creating a new one

#### Scenario: Close tool window cleanup
- **WHEN** the tool window is closed
- **THEN** the system SHALL clean up the window instance reference
