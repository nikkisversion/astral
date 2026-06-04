package store

import (
	"context"
	"errors"

	"github.com/nikkisversion/astral/reader"
	"github.com/qdrant/go-client/qdrant"
)

type QDrantStore struct {
	client         *qdrant.Client
	collectionName string
	dimension      int
}

func NewQDrantClient(ctx context.Context, collectionName string, dimension int) (*QDrantStore, error) {

	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334, //for gRPC
	})

	if err != nil {
		return nil, errors.New("Failed to create QDrant Client: " + err.Error())
	}

	// create a new collection with the specified name and dimension
	errCl := client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     uint64(dimension),
			Distance: qdrant.Distance_Cosine,
		}),
	})

	if errCl != nil {
		return nil, errors.New("Failed to create QDrant Collection: " + errCl.Error())
	}

	return &QDrantStore{
		client:         client,
		collectionName: collectionName,
		dimension:      dimension,
	}, nil

}

func (s *QDrantStore) Dimension() int {
	return s.dimension
}

func (s *QDrantStore) BatchUpsert(ctx context.Context, chunks []reader.Chunk) error {

	points := make([]*qdrant.PointStruct, len(chunks))

	for i, chunk := range chunks {
		// Convert local Go types into Qdrant's payload value map
		payload := map[string]any{
			"name":        chunk.Name,
			"type":        string(chunk.Type),
			"content":     chunk.Content,
			"start_line":  chunk.StartLine,
			"end_line":    chunk.EndLine,
			"source_file": chunk.SourceFile,
		}

		// 2. Build the structural point representation
		points[i] = &qdrant.PointStruct{
			// Generating deterministic sequential IDs for simplicity;
			// in production, consider using UUIDs or another robust ID generation strategy
			Id:      qdrant.NewIDNum(uint64(i + 1)),
			Vectors: qdrant.NewVectors(chunk.Embedding...),
			Payload: qdrant.NewValueMap(payload),
		}
	}

	// 3. Dispatch the bulk request across the local gRPC network line
	_, err := s.client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: s.collectionName,
		Points:         points,
	})

	return err

}
