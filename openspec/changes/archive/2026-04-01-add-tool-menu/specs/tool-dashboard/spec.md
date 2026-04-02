## ADDED Requirements

### Requirement: Tool dashboard display
The system SHALL display a dashboard page as the default view of the tool window.

#### Scenario: Dashboard on tool window open
- **WHEN** the tool window is opened
- **THEN** the system SHALL display the tool dashboard at route `/tools`

### Requirement: Tool cards layout
The system SHALL display available tools in a card-based grid layout.

#### Scenario: Display tool cards
- **WHEN** the dashboard page is displayed
- **THEN** the system SHALL show a card for each available tool
- **AND** each card SHALL display the tool's icon, name, and description
- **AND** cards SHALL be arranged in a responsive grid layout

### Requirement: Tool card interaction
The system SHALL navigate to the selected tool when a card is clicked.

#### Scenario: Click tool card
- **WHEN** user clicks on a tool card
- **THEN** the system SHALL navigate to `/tools/:toolId`
- **AND** the system SHALL display the selected tool's interface
