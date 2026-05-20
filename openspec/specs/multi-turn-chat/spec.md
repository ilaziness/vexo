# Multi-turn Chat

## Purpose

Enable AI assistant to support multi-turn conversations with historical context awareness.

## Requirements

### Requirement: Multi-turn chat context
The system SHALL support multi-turn conversations by passing historical message context to the AI model.

#### Scenario: Chat with existing session
- **WHEN** Chat is called with an existing session ID and a new message
- **THEN** the system SHALL retrieve all historical messages for that session
- **AND** it SHALL pass the message history to the Genkit Flow
- **AND** the AI model SHALL generate a response considering the conversation context

#### Scenario: Chat with new session
- **WHEN** Chat is called with a new session ID and a new message
- **THEN** the system SHALL treat it as the first message in the conversation
- **AND** it SHALL pass only the new message to the Genkit Flow

### Requirement: Chat request structure
The system SHALL define ChatRequest and ChatResponse structures for multi-turn conversations.

#### Scenario: ChatRequest fields
- **WHEN** defining ChatRequest
- **THEN** it SHALL have fields: SessionID (string), Messages ([]AIMessage), NewMessage (string)

#### Scenario: ChatResponse fields
- **WHEN** defining ChatResponse
- **THEN** it SHALL have fields: Message (AIMessage)

### Requirement: Chat service method
The system SHALL provide Chat method in AIService to handle multi-turn conversations.

#### Scenario: Save user message
- **WHEN** Chat is called
- **THEN** the system SHALL save the user message to the database before calling AI

#### Scenario: Generate AI response
- **WHEN** Chat is called
- **THEN** the system SHALL call the Genkit Flow with the conversation context
- **AND** it SHALL wait for the AI response

#### Scenario: Save AI response
- **WHEN** the AI response is received
- **THEN** the system SHALL save the assistant message to the database
- **AND** it SHALL update the session's updated_at timestamp

#### Scenario: Return response
- **WHEN** Chat completes successfully
- **THEN** it SHALL return ChatResponse containing the assistant message

### Requirement: Genkit Flow extension
The system SHALL extend the existing Genkit Flow to support multi-turn conversations.

#### Scenario: Build chat prompt
- **WHEN** the chat Flow is executed
- **THEN** it SHALL build a prompt that includes the message history
- **AND** it SHALL format the history in a way the AI model can understand

#### Scenario: Pass history to AI
- **WHEN** calling genkit.Generate
- **THEN** it SHALL pass the constructed prompt with conversation history
- **AND** the AI model SHALL generate a response aware of previous context

### Requirement: Stop generation
The system SHALL support stopping an in-progress AI generation.

#### Scenario: StopGeneration method
- **WHEN** StopGeneration is called with a session ID
- **THEN** the system SHALL cancel the context for that session's ongoing generation
- **AND** the AI generation SHALL be interrupted

#### Scenario: Context cancellation
- **WHEN** the context is cancelled
- **THEN** the Genkit Flow SHALL receive the cancellation signal
- **AND** it SHALL stop processing and return an error or partial response
