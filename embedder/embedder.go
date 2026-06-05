package embedder

import "context"

type Embedder interface {
	EmbedStrings(ctx context.Context, inputs []string) ([][]float32, error)
	GenerateAnswer(ctx context.Context, messages []ChatMessage) (string, error)
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Message ChatMessage `json:"message"`
}
