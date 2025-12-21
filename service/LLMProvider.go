package service

import (
	"context"

	"mytoy/model"
)

// ChatStreamEventType 表示流式事件类型。
type ChatStreamEventType string

const (
	// ChatEventContent 表示增量内容。
	ChatEventContent ChatStreamEventType = "content"
	// ChatEventError 表示流式过程中的错误。
	ChatEventError ChatStreamEventType = "error"
	// ChatEventEnd 表示流式结束。
	ChatEventEnd ChatStreamEventType = "end"
)

// ChatStreamEvent 是流式聊天的事件，供 SSE 写出使用。
type ChatStreamEvent struct {
	Type         ChatStreamEventType `json:"type"`
	Delta        string              `json:"delta,omitempty"`         // 内容增量
	Message      *model.ChatMessage  `json:"message,omitempty"`       // 汇总消息（可选）
	FinishReason string              `json:"finish_reason,omitempty"` // 结束原因
	Err          error               `json:"-"`                       // 内部错误，不对外序列化
	Raw          any                 `json:"raw,omitempty"`           // 透传底层片段（可选）
}

// ChatStreamCloser 用于在上层主动中断流时释放资源。
type ChatStreamCloser func() error

// LLMProvider 定义统一的多 LLM Provider 接口。
//
// 中断与资源回收推荐组合：
// 1) controller 监听 HTTP 请求的 ctx（含客户端断连），传入 ChatStream；ctx.Done() 触发时 provider 应尽快停止拉流。
// 2) ChatStream 返回的 closer 由上层在循环结束或异常退出时调用（可幂等），内部可二次防护：关闭 SDK 流、释放 goroutine。
// 3) provider 在生产事件时 select ctx.Done()，并在结束/错误时关闭 channel；上层收到 end/error 后立刻退出并调用 closer。
// 该组合覆盖客户端断连、超时、手动终止，避免悬挂流与资源泄露。
type LLMProvider interface {
	Name() string
	ChatStream(ctx context.Context, req *model.ChatRequest) (<-chan ChatStreamEvent, ChatStreamCloser, error)
}
