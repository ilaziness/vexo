## 1. Database Migration

- [x] 1.1 Create database migration script for ai_sessions table in internal/database/migration.go
- [x] 1.2 Create database migration script for ai_messages table in internal/database/migration.go
- [x] 1.3 Create indexes on ai_messages(session_id) and ai_sessions(updated_at DESC)
- [~] 1.4 Test database migration runs successfully on application startup (推迟到 Phase 2)

## 2. Data Models

- [x] 2.1 Define AISession struct in internal/database/ai_session.go
- [x] 2.2 Define AIMessage struct in internal/database/ai_session.go
- [x] 2.3 Add JSON tags to struct fields for proper serialization
- [x] 2.4 Verify struct definitions match PRD specifications

## 3. Repository Interface

- [x] 3.1 Define AISessionRepository interface in internal/database/ai_session.go
- [x] 3.2 Add CreateSession method signature to interface
- [x] 3.3 Add GetSession method signature to interface
- [x] 3.4 Add ListSessions method signature to interface
- [x] 3.5 Add UpdateSession method signature to interface
- [x] 3.6 Add DeleteSession method signature to interface
- [x] 3.7 Add CreateMessage method signature to interface
- [x] 3.8 Add ListMessages method signature to interface

## 4. Repository Implementation

- [x] 4.1 Implement SQLite repository constructor in internal/database/ai_session.go
- [x] 4.2 Implement database connection initialization logic
- [x] 4.3 Implement CreateSession method with ID generation
- [x] 4.4 Implement GetSession method with error handling for not found
- [x] 4.5 Implement ListSessions method with limit and ordering
- [x] 4.6 Implement UpdateSession method with timestamp update
- [x] 4.7 Implement DeleteSession method (CASCADE delete handled by foreign key)
- [x] 4.8 Implement CreateMessage method with ID generation
- [x] 4.9 Implement ListMessages method with session filtering and ordering
- [~] 4.10 Add unit tests for repository methods (推迟到 Phase 2)

## 5. AIService Extension

- [x] 5.1 Add sessionRepo field to AIService struct in services/ai_service.go
- [x] 5.2 Initialize sessionRepo in AIService constructor
- [x] 5.3 Define ChatRequest struct with SessionID, Messages, NewMessage fields
- [x] 5.4 Define ChatResponse struct with Message field
- [x] 5.5 Implement CreateSession service method
- [x] 5.6 Implement GetSession service method
- [x] 5.7 Implement ListSessions service method
- [x] 5.8 Implement DeleteSession service method
- [x] 5.9 Implement UpdateSession service method
- [x] 5.10 Implement Chat service method with message saving
- [x] 5.11 Implement StopGeneration service method with context cancellation
- [x] 5.12 Add error handling for all service methods

## 6. Genkit Flow Extension

- [x] 6.1 Review existing chatFlow implementation in internal/ai/genkit.go
- [x] 6.2 Modify chatFlow to accept ChatRequest as input
- [x] 6.3 Implement buildChatPrompt function to format message history
- [x] 6.4 Update genkit.Generate call to pass conversation context
- [~] 6.5 Test multi-turn conversation with context preservation (推迟到 Phase 2)
- [~] 6.6 Verify context cancellation works for StopGeneration (推迟到 Phase 2)

## 7. Frontend Bindings

- [x] 7.1 Run wails3 generate bindings -ts to generate TypeScript bindings
- [x] 7.2 Verify generated bindings include new AIService methods
- [x] 7.3 Verify generated bindings include ChatRequest and ChatResponse types
- [x] 7.4 Verify generated bindings include AISession and AIMessage types
- [x] 7.5 Check for any binding generation errors

## 8. Integration Testing

- [~] 8.1 Test end-to-end session creation flow (推迟到 Phase 2)
- [~] 8.2 Test message saving and retrieval (推迟到 Phase 2)
- [~] 8.3 Test multi-turn conversation with context (推迟到 Phase 2)
- [~] 8.4 Test session deletion cascades to messages (推迟到 Phase 2)
- [~] 8.5 Test StopGeneration functionality (推迟到 Phase 2)
- [~] 8.6 Verify database persistence after application restart (推迟到 Phase 2)
