## ADDED Requirements

### Requirement: AI session persistence
The system SHALL persist AI chat sessions and messages to a local SQLite database.

#### Scenario: Create session table
- **WHEN** the application starts
- **THEN** the system SHALL create the ai_sessions table if it does not exist
- **AND** the table SHALL contain columns: id (TEXT PRIMARY KEY), title (TEXT NOT NULL), created_at (INTEGER NOT NULL), updated_at (INTEGER NOT NULL)

#### Scenario: Create message table
- **WHEN** the application starts
- **THEN** the system SHALL create the ai_messages table if it does not exist
- **AND** the table SHALL contain columns: id (TEXT PRIMARY KEY), session_id (TEXT NOT NULL), role (TEXT NOT NULL), content (TEXT NOT NULL), timestamp (INTEGER NOT NULL)
- **AND** session_id SHALL be a foreign key referencing ai_sessions(id) with CASCADE delete

#### Scenario: Create indexes
- **WHEN** the database is initialized
- **THEN** the system SHALL create an index on ai_messages(session_id)
- **AND** the system SHALL create an index on ai_sessions(updated_at DESC)

### Requirement: Session data model
The system SHALL define AISession and AIMessage data structures for database operations.

#### Scenario: AISession structure
- **WHEN** defining AISession
- **THEN** it SHALL have fields: ID (string), Title (string), CreatedAt (time.Time), UpdatedAt (time.Time)

#### Scenario: AIMessage structure
- **WHEN** defining AIMessage
- **THEN** it SHALL have fields: ID (string), SessionID (string), Role (string), Content (string), Timestamp (time.Time)
- **AND** Role SHALL be either 'user' or 'assistant'

### Requirement: Repository interface
The system SHALL provide AISessionRepository interface for data access operations.

#### Scenario: Repository session methods
- **WHEN** defining AISessionRepository
- **THEN** it SHALL have methods: CreateSession, GetSession, ListSessions, UpdateSession, DeleteSession

#### Scenario: Repository message methods
- **WHEN** defining AISessionRepository
- **THEN** it SHALL have methods: CreateMessage, ListMessages

### Requirement: SQLite repository implementation
The system SHALL implement AISessionRepository using modernc.org/sqlite driver.

#### Scenario: Initialize database connection
- **WHEN** creating the repository
- **THEN** it SHALL open or create the SQLite database file
- **AND** it SHALL run schema migration to create tables

#### Scenario: Create session in database
- **WHEN** CreateSession is called
- **THEN** the system SHALL insert a new row into ai_sessions table
- **AND** it SHALL return the created AISession with generated ID

#### Scenario: Get session by ID
- **WHEN** GetSession is called with a valid session ID
- **THEN** the system SHALL query ai_sessions table by id
- **AND** it SHALL return the AISession if found
- **AND** it SHALL return error if not found

#### Scenario: List sessions with limit
- **WHEN** ListSessions is called with a limit
- **THEN** the system SHALL query ai_sessions table ordered by updated_at DESC
- **AND** it SHALL return up to the specified number of sessions

#### Scenario: Delete session cascades messages
- **WHEN** DeleteSession is called
- **THEN** the system SHALL delete the session from ai_sessions table
- **AND** due to CASCADE delete, all related messages SHALL be automatically deleted

#### Scenario: Create message in database
- **WHEN** CreateMessage is called
- **THEN** the system SHALL insert a new row into ai_messages table
- **AND** it SHALL return the created AIMessage with generated ID

#### Scenario: List messages by session
- **WHEN** ListMessages is called with a session ID
- **THEN** the system SHALL query ai_messages table by session_id ordered by timestamp ASC
- **AND** it SHALL return all messages for that session
