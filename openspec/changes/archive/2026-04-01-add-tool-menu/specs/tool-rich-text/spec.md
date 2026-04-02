## ADDED Requirements

### Requirement: Rich text editor tool interface
The system SHALL provide a rich text editor tool using react-quill.

#### Scenario: Display rich text editor tool
- **WHEN** user navigates to `/tools/rich-text`
- **THEN** the system SHALL display the rich text editor interface

### Requirement: Visual editing mode
The system SHALL provide a visual WYSIWYG editing mode.

#### Scenario: Visual editing
- **WHEN** the rich text tool is displayed
- **THEN** the system SHALL show a react-quill editor with toolbar
- **AND** the toolbar SHALL include formatting options: bold, italic, underline, strikethrough
- **AND** the toolbar SHALL include heading styles (H1, H2, H3)
- **AND** the toolbar SHALL include list options (ordered, unordered)
- **AND** the toolbar SHALL include link and image insertion
- **AND** the toolbar SHALL include alignment options
- **AND** the toolbar SHALL include color options (text and background)

### Requirement: Source code mode
The system SHALL provide an HTML source code editing mode.

#### Scenario: Toggle source mode
- **WHEN** user clicks the source mode toggle button
- **THEN** the system SHALL switch from visual mode to source mode
- **AND** the source mode SHALL display the HTML content in a textarea
- **AND** the HTML SHALL be editable

#### Scenario: Toggle back to visual mode
- **WHEN** user clicks the visual mode toggle button
- **THEN** the system SHALL switch from source mode back to visual mode
- **AND** the visual editor SHALL reflect any changes made in source mode

### Requirement: Mode synchronization
The system SHALL keep visual and source modes synchronized.

#### Scenario: Content sync from visual to source
- **WHEN** user edits content in visual mode and switches to source mode
- **THEN** the source mode SHALL display the updated HTML

#### Scenario: Content sync from source to visual
- **WHEN** user edits HTML in source mode and switches to visual mode
- **THEN** the visual editor SHALL render the updated content
- **AND** invalid HTML SHALL be handled gracefully

### Requirement: Copy HTML source
The system SHALL allow users to copy the HTML source to clipboard.

#### Scenario: Click copy button
- **WHEN** user clicks the copy HTML button
- **THEN** the system SHALL copy the current HTML content to clipboard
- **AND** the system SHALL provide visual feedback indicating success

### Requirement: Clear content
The system SHALL provide a way to clear all content.

#### Scenario: Click clear button
- **WHEN** user clicks the clear button
- **THEN** the system SHALL clear all content from the editor
- **AND** the system SHALL require confirmation if content is not empty

### Requirement: Editor theming
The system SHALL support the application's theme (light/dark mode).

#### Scenario: Theme adaptation
- **WHEN** the application theme changes
- **THEN** the rich text editor SHALL adapt to the current theme
- **AND** the editor background and text colors SHALL match the theme
