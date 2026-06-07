package store

import (
	"context"

	"github.com/nikkisversion/astral/reader"
)

type ScoredChunk struct {
	Content string  `json:"content"`
	Score   float32 `json:"score"`
}

type Store interface {
	Dimension() int
	BatchUpsert(ctx context.Context, chunks []reader.Chunk) error
	SearchSimilar(ctx context.Context, embedding []float32, limit uint64) ([]ScoredChunk, error)
}
