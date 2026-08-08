package helix

// TokenUsage mirrors Helix runtime token accounting.
type TokenUsage struct {
	PromptTokens     int     `json:"promptTokens"`
	CompletionTokens int     `json:"completionTokens"`
	TotalTokens      int     `json:"totalTokens"`
	EstimatedCostUSD float64 `json:"estimatedCostUsd"`
}

// ApprovalRequest is a parked tool call waiting on a human.
type ApprovalRequest struct {
	ID        string `json:"id"`
	ToolName  string `json:"toolName"`
	Input     any    `json:"input"`
	CreatedAt string `json:"createdAt"`
	Status    string `json:"status"`
}

// ChatMessage is one turn in a durable session transcript.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

// SessionRecord is a durable Helix session.
type SessionRecord struct {
	ID               string            `json:"id"`
	CreatedAt        string            `json:"createdAt"`
	UpdatedAt        string            `json:"updatedAt"`
	Messages         []ChatMessage     `json:"messages"`
	Usage            TokenUsage        `json:"usage"`
	Status           string            `json:"status"`
	PendingApprovals []ApprovalRequest `json:"pendingApprovals"`
	WorkflowID       string            `json:"workflowId,omitempty"`
	Channel          string            `json:"channel,omitempty"`
}

// RuntimeEvent is one forensic timeline entry.
type RuntimeEvent struct {
	Type      string         `json:"type"`
	At        string         `json:"at"`
	SessionID string         `json:"sessionId"`
	Data      map[string]any `json:"data,omitempty"`
}

// ToolCallResult is a tool invocation recorded on a turn.
type ToolCallResult struct {
	Name   string `json:"name"`
	Input  any    `json:"input"`
	Output any    `json:"output"`
}

// RunResult is the response from POST /helix/v1/sessions.
type RunResult struct {
	SessionID  string           `json:"sessionId"`
	WorkflowID string           `json:"workflowId,omitempty"`
	Reply      string           `json:"reply"`
	Usage      TokenUsage       `json:"usage"`
	ToolCalls  []ToolCallResult `json:"toolCalls"`
	Events     []RuntimeEvent   `json:"events"`
	Parked     bool             `json:"parked,omitempty"`
	ModelUsed  string           `json:"modelUsed,omitempty"`
}

// WorkflowRun is a durable workflow under .helix/workflows.
type WorkflowRun struct {
	ID        string         `json:"id"`
	SessionID string         `json:"sessionId"`
	Status    string         `json:"status"`
	CreatedAt string         `json:"createdAt"`
	UpdatedAt string         `json:"updatedAt"`
	Steps     []any          `json:"steps"`
	Meta      map[string]any `json:"meta,omitempty"`
}

// RunOptions configures a chat turn.
type RunOptions struct {
	Message     string
	SessionID   string
	AutoApprove bool
	Channel     string
}
