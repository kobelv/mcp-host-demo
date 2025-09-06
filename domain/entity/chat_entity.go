package entity

type ChatInputEntity struct {
	Query   string `json:"query"`
	Remarks string `json:"remarks"`
}

type LLMTool struct {
	ToolName    string                 `json:"tool_name"`
	ToolDesc    string                 `json:"tool_desc"`
	InputParams map[string]interface{} `json:"input_params"`
	Required    []string               `json:"required"`
}
