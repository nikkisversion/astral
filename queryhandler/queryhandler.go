package queryhandler

import (
	"context"
	"fmt"
	"strings"

	"github.com/nikkisversion/astral/embedder"
	"github.com/nikkisversion/astral/store"
)

type QueryHandler struct {
	Embedder embedder.Embedder
	Store    store.Store
}

func New(e embedder.Embedder, s store.Store) *QueryHandler {
	return &QueryHandler{Embedder: e, Store: s}
}

func (q *QueryHandler) Handle(ctx context.Context, query string) (string, error) {

	queryEmbedding, errQE := q.Embedder.EmbedStrings(ctx, []string{query})
	if errQE != nil {
		fmt.Printf("Error embedding NLQ: %v\n", errQE)
		return "", errQE
	}

	scoredChunks, errSS := q.Store.SearchSimilar(ctx, queryEmbedding[0], 3)
	if errSS != nil {
		fmt.Printf("Error searching similar in store: %v\n", errSS)
		return "", errSS
	}

	chatMessages := convertChunksToChatMessage(scoredChunks, query)

	answer, errAns := q.Embedder.GenerateAnswer(ctx, chatMessages)
	if errAns != nil {
		fmt.Printf("Error generating answer: %v\n", errAns)
		return "", errAns
	}

	return answer, nil
}

func convertChunksToChatMessage(chunks []store.ScoredChunk, query string) []embedder.ChatMessage {

	// Combine  chunks into a single text block
	var contextBlock strings.Builder
	for _, chunk := range chunks {
		contextBlock.WriteString("\n```go\n" + chunk.Content + "\n```\n")
	}

	// 2. Format the messages
	messages := []embedder.ChatMessage{
		{Role: "system", Content: "You are a Go code assistant. Provide a descriptive, medium-length answer based only on the code provided."},
		{Role: "user", Content: fmt.Sprintf("Code Context:\n%s\n\nQuestion: %s", contextBlock.String(), query)},
	}

	return messages
}
