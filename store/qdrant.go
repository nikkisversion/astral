package store

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/nikkisversion/astral/reader"
	"github.com/qdrant/go-client/qdrant"
)

type QDrantStore struct {
	client         *qdrant.Client
	collectionName string
	dimension      int
}

func NewQDrantStore(ctx context.Context, collectionName string, dimension int) (*QDrantStore, error) {

	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334, //for gRPC
	})

	if err != nil {
		return nil, errors.New("Failed to create QDrant Client: " + err.Error())
	}

	// 1. Ask Qdrant if the collection already exists
	exists, errExist := client.CollectionExists(ctx, collectionName)
	if errExist != nil {
		return nil, errors.New("failed checking collection existence: " + errExist.Error())
	}

	if exists {
		errDel := client.DeleteCollection(ctx, collectionName)
		if errDel != nil {
			return nil, errors.New("Failed to delete existing QDrant Collection: " + errDel.Error())
		}
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
			// Generate a unique ID for each chunk to avoid collision while storing in Qdrant
			Id:      qdrant.NewIDUUID(GenerateNativeID(chunk.SourceFile, chunk.Content)),
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

func (s *QDrantStore) SearchSimilar(ctx context.Context, embedding []float32, limit uint64) ([]ScoredChunk, error) {

	req := &qdrant.QueryPoints{
		CollectionName: s.collectionName,
		Query:          qdrant.NewQuery(embedding...),
		Limit:          &limit,
		WithPayload:    qdrant.NewWithPayload(true),
	}

	resp, err := s.client.Query(ctx, req)
	if err != nil {
		return nil, err
	}

	scoredChunks := make([]ScoredChunk, 0, len(resp))
	for _, point := range resp {

		if content, ok := point.Payload["content"]; ok {
			sChunk := ScoredChunk{
				Content: content.GetStringValue(),
				Score:   point.Score,
			}
			scoredChunks = append(scoredChunks, sChunk)
		}

	}

	return scoredChunks, nil
}

// GenerateNativeID produces a deterministic 36-character string matching UUID format
func GenerateNativeID(filePath string, chunkContent string) string {
	// Combine properties to form an absolute unique identifier
	uniqueSalt := filePath + "::" + chunkContent

	// Create a SHA-256 hash
	hash := sha256.Sum256([]byte(uniqueSalt))
	encoded := hex.EncodeToString(hash[:])

	// Format the raw hex into a valid UUID pattern: 8-4-4-4-12 characters
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		encoded[0:8],
		encoded[8:12],
		encoded[12:16],
		encoded[16:20],
		encoded[20:32],
	)
}
