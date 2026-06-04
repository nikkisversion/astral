package store

import (
	"context"
	"errors"

	"github.com/nikkisversion/astral/reader"
)

type Store interface {
	Dimension() int
	BatchUpsert(ctx context.Context, chunks []reader.Chunk) error
}

// VectorStore is an abstraction layer over the underlying vector database
type VectorStore struct {
	Client    Store
	Dimension int
}

type VectorStoreType string

const (
	QDrantType VectorStoreType = "qdrant"
)

func NewVectorStore(ctx context.Context, dimension int) (*VectorStore, error) {

	// Enhancement: Add logic to choose the appropriate vector store based on configuration or environment variables
	// For now, we will default to QDrant

	collectionName := "code_chunks"

	client, err := NewQDrantClient(ctx, collectionName, dimension)
	if err != nil {
		return nil, errors.New("Failed to create QDrant Store: " + err.Error())
	}

	return &VectorStore{
		Client:    client,
		Dimension: dimension,
	}, nil
}
