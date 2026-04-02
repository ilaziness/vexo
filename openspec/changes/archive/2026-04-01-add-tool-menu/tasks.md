## 1. Go Backend - Tool Service

- [x] 1.1 Add `ToolWindow` field to `AppInstance` struct in `services/app.go`
- [x] 1.2 Create `newToolWindow()` function in `services/app.go`
- [x] 1.3 Create `ToolService` struct in `services/tool_service.go`
- [x] 1.4 Implement `ShowWindow()` method in `ToolService`
- [x] 1.5 Implement `CloseWindow()` method in `ToolService`
- [x] 1.6 Implement `GetTools()` method returning available tools list
- [x] 1.7 Register `ToolService` in `RegisterServices()` in `services/service.go`
- [x] 1.8 Run `wails3 generate bindings -ts` to generate frontend bindings

## 2. Frontend - Routes and Types

- [x] 2.1 Create `Tool` interface in `frontend/src/types/tool.ts`
- [x] 2.2 Add `/tools` route to `frontend/src/routes.ts`
- [x] 2.3 Add `/tools/:toolId` route to `frontend/src/routes.ts`
- [x] 2.4 Create tool store to fetch and cache tools list from backend

## 3. Frontend - Tool Dashboard Page

- [x] 3.1 Create `Tools.tsx` page component in `frontend/src/pages/`
- [x] 3.2 Implement card grid layout using MUI Grid
- [x] 3.3 Create `ToolCard` component with icon, name, description
- [x] 3.4 Add click handler to navigate to `/tools/:toolId`
- [x] 3.5 Style cards with hover effects

## 4. Frontend - Tool Layout Component

- [x] 4.1 Create `ToolLayout.tsx` component in `frontend/src/components/`
- [x] 4.2 Implement left sidebar with tool list
- [x] 4.3 Add active tool highlighting in sidebar
- [x] 4.4 Add click handler for tool switching
- [x] 4.5 Add back button to return to dashboard
- [x] 4.6 Style sidebar with fixed 200px width

## 5. Go Backend - Port Check Logic

- [x] 5.1 Define `PortCheckResult` struct in `tool_service.go`
- [x] 5.2 Implement `CheckPort(host string, port int) PortCheckResult` method
- [x] 5.3 Implement TCP connection with timeout
- [x] 5.4 Measure and return connection response time

## 6. Frontend - Port Check Tool

- [x] 6.1 Create `PortCheckTool.tsx` component
- [x] 6.2 Add host input field with validation
- [x] 6.3 Add port input field (1-65535 range)
- [x] 6.4 Call `ToolService.CheckPort()` on button click
- [x] 6.5 Display connection result (success/failure)
- [x] 6.6 Show response time for successful connections

## 7. Go Backend - Encoder Logic

- [x] 7.1 Define `EncodeRequest` and `EncodeResponse` structs
- [x] 7.2 Implement `Encode(toolType string, input string) string` method
- [x] 7.3 Implement `Decode(toolType string, input string) string` method
- [x] 7.4 Support Base64, URL, HTML encoding types

## 8. Frontend - Encoder Tool

- [x] 8.1 Create `EncoderTool.tsx` component
- [x] 8.2 Add encoding type selector (Base64, URL, HTML)
- [x] 8.3 Add source text input area
- [x] 8.4 Add result output area (read-only)
- [x] 8.5 Call `ToolService.Encode()` / `ToolService.Decode()` on button click
- [x] 8.6 Add copy result button with feedback

## 9. Go Backend - Regex Logic

- [x] 9.1 Define `RegexMatchResult` and `Match` structs
- [x] 9.2 Implement `RegexMatch(pattern string, text string, flags string) RegexMatchResult` method
- [x] 9.3 Parse regex flags (i, m, g) and apply to pattern
- [x] 9.4 Return match positions and capture groups
- [x] 9.5 Return error for invalid patterns

## 10. Frontend - Regex Tool

- [x] 10.1 Create `RegexTool.tsx` component
- [x] 10.2 Add regex pattern input
- [x] 10.3 Add test text input area
- [x] 10.4 Add regex flags toggles (i, m, g)
- [x] 10.5 Call `ToolService.RegexMatch()` on input change
- [x] 10.6 Display match count
- [x] 10.7 Display capture groups for each match
- [x] 10.8 Show error message for invalid patterns

## 11. Frontend - Rich Text Editor Tool

- [x] 11.1 Install `react-quill` and `quill` packages
- [x] 11.2 Create `RichTextTool.tsx` component
- [x] 11.3 Integrate ReactQuill editor with toolbar
- [x] 11.4 Implement visual editing mode
- [x] 11.5 Add toggle to switch between visual and source (HTML) mode
- [x] 11.6 Implement source mode with HTML textarea
- [x] 11.7 Sync content between visual and source modes
- [x] 11.8 Add copy HTML source button
- [x] 11.9 Style editor to match application theme

## 12. Frontend - Header Integration

- [x] 12.1 Import `Handyman` icon from MUI icons
- [x] 12.2 Add tool button to `Header.tsx`
- [x] 12.3 Add click handler to call `ToolService.ShowWindow()`
- [x] 12.4 Add tooltip for the tool button
- [x] 12.5 Position button appropriately in the header

## 13. Testing and Verification

- [x] 13.1 Verify tool window opens correctly
- [x] 13.2 Verify dashboard displays tools from backend
- [x] 13.3 Verify navigation from dashboard to tool works
- [x] 13.4 Verify sidebar tool switching works
- [x] 13.5 Verify back button returns to dashboard
- [x] 13.6 Test port check tool functionality via backend
- [x] 13.7 Test encoder tool with all encoding types via backend
- [x] 13.8 Test regex tool with various patterns via backend
- [x] 13.9 Test rich text editor visual and source modes
- [x] 13.10 Run `wails3 dev` to verify full integration