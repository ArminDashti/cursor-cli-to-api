package openai

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ChatMessage struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"`
}

type ChatCompletionRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

type ChatCompletionResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChatCompletionChoice `json:"choices"`
	Usage   Usage                  `json:"usage"`
}

type ChatCompletionChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ResponsesRequest struct {
	Model  string          `json:"model"`
	Input  json.RawMessage `json:"input"`
	Stream bool            `json:"stream"`
}

type ResponseObject struct {
	ID        string           `json:"id"`
	Object    string           `json:"object"`
	CreatedAt int64            `json:"created_at"`
	Status    string           `json:"status"`
	Model     string           `json:"model"`
	Output    []ResponseOutput `json:"output"`
	Usage     ResponseUsage    `json:"usage"`
}

type ResponseOutput struct {
	ID      string               `json:"id"`
	Type    string               `json:"type"`
	Status  string               `json:"status"`
	Role    string               `json:"role"`
	Content []ResponseContentBit `json:"content"`
}

type ResponseContentBit struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ResponseUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

func NewChatID() string {
	return "chatcmpl-" + strings.ReplaceAll(uuid.NewString(), "-", "")
}

func NewResponseID() string {
	return "resp_" + strings.ReplaceAll(uuid.NewString(), "-", "")
}

func EstimateTokens(s string) int {
	n := len(strings.Fields(s))
	if n == 0 && s != "" {
		return 1
	}
	return n
}

func ContentToString(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	var parts []map[string]any
	if err := json.Unmarshal(raw, &parts); err == nil {
		var b strings.Builder
		for _, p := range parts {
			if t, ok := p["type"].(string); ok && t == "text" {
				if txt, ok := p["text"].(string); ok {
					b.WriteString(txt)
				}
			} else if txt, ok := p["text"].(string); ok {
				b.WriteString(txt)
			}
		}
		return b.String()
	}
	return string(raw)
}

func PromptFromChatMessages(msgs []ChatMessage) string {
	var b strings.Builder
	for _, m := range msgs {
		role := strings.TrimSpace(m.Role)
		if role == "" {
			role = "user"
		}
		b.WriteString(strings.ToUpper(role))
		b.WriteString(":\n")
		b.WriteString(ContentToString(m.Content))
		b.WriteString("\n\n")
	}
	return strings.TrimSpace(b.String())
}

func PromptFromResponsesInput(raw json.RawMessage) (string, error) {
	if len(raw) == 0 {
		return "", fmt.Errorf("input is required")
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s, nil
	}
	var arr []json.RawMessage
	if err := json.Unmarshal(raw, &arr); err == nil {
		var b strings.Builder
		for _, item := range arr {
			var obj map[string]any
			if err := json.Unmarshal(item, &obj); err != nil {
				b.WriteString(string(item))
				b.WriteString("\n")
				continue
			}
			role, _ := obj["role"].(string)
			if role == "" {
				role = "user"
			}
			b.WriteString(strings.ToUpper(role))
			b.WriteString(":\n")
			switch c := obj["content"].(type) {
			case string:
				b.WriteString(c)
			default:
				rawC, _ := json.Marshal(c)
				b.WriteString(ContentToString(rawC))
			}
			b.WriteString("\n\n")
		}
		return strings.TrimSpace(b.String()), nil
	}
	return string(raw), nil
}

func BuildChatCompletion(model, text string) ChatCompletionResponse {
	if model == "" {
		model = "auto"
	}
	promptTok := 0
	compTok := EstimateTokens(text)
	content, _ := json.Marshal(text)
	return ChatCompletionResponse{
		ID:      NewChatID(),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []ChatCompletionChoice{{
			Index: 0,
			Message: ChatMessage{
				Role:    "assistant",
				Content: content,
			},
			FinishReason: "stop",
		}},
		Usage: Usage{
			PromptTokens:     promptTok,
			CompletionTokens: compTok,
			TotalTokens:      promptTok + compTok,
		},
	}
}

func BuildResponse(model, text string) ResponseObject {
	if model == "" {
		model = "auto"
	}
	id := NewResponseID()
	outTok := EstimateTokens(text)
	return ResponseObject{
		ID:        id,
		Object:    "response",
		CreatedAt: time.Now().Unix(),
		Status:    "completed",
		Model:     model,
		Output: []ResponseOutput{{
			ID:     "msg_" + strings.ReplaceAll(uuid.NewString(), "-", ""),
			Type:   "message",
			Status: "completed",
			Role:   "assistant",
			Content: []ResponseContentBit{{
				Type: "output_text",
				Text: text,
			}},
		}},
		Usage: ResponseUsage{
			InputTokens:  0,
			OutputTokens: outTok,
			TotalTokens:  outTok,
		},
	}
}

func ChatStreamChunk(id, model, delta string, finish *string) map[string]any {
	chunk := map[string]any{
		"id":      id,
		"object":  "chat.completion.chunk",
		"created": time.Now().Unix(),
		"model":   model,
		"choices": []map[string]any{{
			"index": 0,
			"delta": map[string]any{
				"content": delta,
			},
			"finish_reason": finish,
		}},
	}
	return chunk
}
