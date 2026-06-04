package embedder

import "context"

type Embedder interface {
	EmbedStrings(ctx context.Context, inputs []string) ([][]float32, error)
}
