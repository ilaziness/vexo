package services

import "github.com/ilaziness/vexo/internal/tools"

type Tool = tools.Tool
type PortCheckResult = tools.PortCheckResult
type EncodeResponse = tools.EncodeResponse
type RegexMatchResult = tools.RegexMatchResult
type HashResult = tools.HashResult
type TimestampResult = tools.TimestampResult

type ToolService struct {
	windows *Windows
	core    *tools.Service
}

func NewToolService(windows *Windows) *ToolService {
	return &ToolService{windows: windows, core: tools.New()}
}

func (ts *ToolService) ShowWindow()  { ts.windows.ShowTool() }
func (ts *ToolService) CloseWindow() { ts.windows.CloseTool() }
func (ts *ToolService) GetTools() []Tool { return ts.core.GetTools() }
func (ts *ToolService) CheckPort(host string, port int) PortCheckResult {
	return ts.core.CheckPort(host, port)
}
func (ts *ToolService) Encode(toolType, input string) EncodeResponse {
	return ts.core.Encode(toolType, input)
}
func (ts *ToolService) Decode(toolType, input string) EncodeResponse {
	return ts.core.Decode(toolType, input)
}
func (ts *ToolService) RegexMatch(pattern, text, flags string) RegexMatchResult {
	return ts.core.RegexMatch(pattern, text, flags)
}
func (ts *ToolService) CalculateHash(input, algorithm string) HashResult {
	return ts.core.CalculateHash(input, algorithm)
}
func (ts *ToolService) ConvertTimestamp(input string, toTimestamp bool) TimestampResult {
	return ts.core.ConvertTimestamp(input, toTimestamp)
}
