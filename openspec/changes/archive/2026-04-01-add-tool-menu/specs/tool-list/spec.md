## ADDED Requirements

### Requirement: Tool list provided by backend
The system SHALL provide an API to retrieve the list of available tools from the backend.

#### Scenario: Get tools from backend
- **WHEN** the frontend requests the tool list
- **THEN** the system SHALL return a list of available tools
- **AND** each tool SHALL have: id, name, description, icon, category

### Requirement: Tool data structure
The system SHALL define a consistent tool data structure.

#### Scenario: Tool struct definition
- **WHEN** tools are returned from the backend
- **THEN** each tool SHALL contain:
  - `id`: unique string identifier
  - `name`: display name
  - `description`: brief description
  - `icon`: icon identifier
  - `category`: tool category

### Requirement: Dynamic tool list
The system SHALL support dynamic tool list from backend.

#### Scenario: Frontend displays tools from backend
- **WHEN** the tool dashboard is displayed
- **THEN** the frontend SHALL call the backend to get the tool list
- **AND** the dashboard SHALL display tools returned by the backend
