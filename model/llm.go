package model

const (
	ChatRoleSystem    string = "system"
	ChatRoleUser      string = "user"
	ChatRoleAssistant string = "assistant"
)

// ChatMessage 表示一条对话消息。
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content,omitempty"`
}

// ChatRequest 是统一的聊天请求载荷。
type ChatRequest struct {
	Provider    string        `json:"provider,omitempty"` // 例如 openai/anthropic/azure
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float32       `json:"temperature,omitempty"`
	TopP        float32       `json:"top_p,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
}

// ChatResponse 是统一的聊天响应。
type ChatResponse struct {
	Provider string      `json:"provider,omitempty"`
	Model    string      `json:"model,omitempty"`
	Message  ChatMessage `json:"message"`
}
