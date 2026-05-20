## ADDED Requirements

### Requirement: Create session
The system SHALL allow creation of new AI chat sessions.

#### Scenario: CreateSession generates ID
- **WHEN** CreateSession is called
- **THEN** the system SHALL generate a unique session ID
- **AND** it SHALL create a session with empty title
- **AND** it SHALL set created_at and updated_at to current time
- **AND** it SHALL save the session to the database
- **AND** it SHALL return the created AISession

#### Scenario: CreateSession initial state
- **WHEN** a new session is created
- **THEN** it SHALL have no associated messages
- **AND** it SHALL be ready for the first user message

### Requirement: Get session
The system SHALL allow retrieval of a specific session by ID.

#### Scenario: GetSession by ID
- **WHEN** GetSession is called with a valid session ID
- **THEN** the system SHALL retrieve the session from the database
- **AND** it SHALL return the AISession object

#### Scenario: GetSession not found
- **WHEN** GetSession is called with an invalid session ID
- **THEN** the system SHALL return an error indicating session not found

### Requirement: List sessions
The system SHALL allow listing of sessions with optional limit.

#### Scenario: ListSessions default
- **WHEN** ListSessions is called without a limit
- **THEN** it SHALL return all sessions ordered by updated_at DESC

#### Scenario: ListSessions with limit
- **WHEN** ListSessions is called with a limit parameter
- **THEN** it SHALL return up to the specified number of sessions
- **AND** the sessions SHALL be ordered by updated_at DESC

#### Scenario: ListSessions empty
- **WHEN** no sessions exist in the database
- **THEN** ListSessions SHALL return an empty array

### Requirement: Delete session
The system SHALL allow deletion of sessions.

#### Scenario: DeleteSession by ID
- **WHEN** DeleteSession is called with a valid session ID
- **THEN** the system SHALL delete the session from the database
- **AND** due to CASCADE delete, all associated messages SHALL be deleted
- **AND** it SHALL return nil on success

#### Scenario: DeleteSession not found
- **WHEN** DeleteSession is called with an invalid session ID
- **THEN** the system SHALL return an error indicating session not found

### Requirement: Update session
The system SHALL allow updating session metadata.

#### Scenario: UpdateSession title
- **WHEN** UpdateSession is called with a session ID and new title
- **THEN** the system SHALL update the title field in the database
- **AND** it SHALL update the updated_at timestamp to current time
- **AND** it SHALL return the updated AISession

#### Scenario: UpdateSession not found
- **WHEN** UpdateSession is called with an invalid session ID
- **THEN** the system SHALL return an error indicating session not found
