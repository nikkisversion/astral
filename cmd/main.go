package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/nikkisversion/astral/embedder"
	"github.com/nikkisversion/astral/reader"
	"github.com/nikkisversion/astral/store"
)

/* Steps:
 1. Run Ollama application and on terminal, pull the nomic-embed-text model using: `ollama pull nomic-embed-text`
 2. Run the docker command to start qdrant: docker run -p 6333:6333 -p 6334:6334 \
    -v "$(pwd)/qdrant_storage:/qdrant/storage:z" \
    qdrant/qdrant
3. Run this main.go file using `go run cmd/main.go` and see the results in the terminal.
4. Open 'http://localhost:6333/dashboard' in your browser to see the QDrant dashboard and verify that the collection and points have been created.
*/

func main() {

	fmt.Println("Hello from Astral! LFG!")

	filePath := "testFile.go"

	newReader, err := reader.New(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	errRead := newReader.Read()
	if errRead != nil {
		fmt.Println("Error reading file:", errRead)
		return
	}

	inputs := newReader.GenerateInputForEmbedding()

	newEmbedder := embedder.NewOllamaEmbedder("http://localhost:11434", "nomic-embed-text", "qwen2.5-coder:1.5b")
	embeddings, err := newEmbedder.EmbedStrings(context.Background(), inputs)
	if err != nil {
		fmt.Printf("Error embedding strings: %v\n", err)
		return
	}

	newReader.UpdateEmbedding(embeddings)

	fmt.Printf("Collected %d chunks from %s\n", len(newReader.Collector.Chunks), filePath)

	ctx := context.Background()
	// TODO: Add logic to get dimension of current model in Ollama dynamically instead of hardcoding it
	newVS, err := store.NewVectorStore(ctx, 768)
	if err != nil {
		fmt.Printf("Error creating vector store: %v\n", err)
		return
	}

	errUpsert := newVS.Client.BatchUpsert(ctx, newReader.Collector.Chunks)
	if errUpsert != nil {
		fmt.Printf("Error upserting chunks to vector store: %v\n", errUpsert)
		return
	}

	fmt.Println("Successfully upserted chunks to vector store!")

	nlQuery := "What does the NLTestFunc do?"
	queryEmbedding, errQE := newEmbedder.EmbedStrings(ctx, []string{nlQuery})
	if errQE != nil {
		fmt.Printf("Error embedding NLQ: %v\n", errQE)
		return
	}

	scoredChunks, errSS := newVS.Client.SearchSimilar(ctx, queryEmbedding[0], 3)
	if errSS != nil {
		fmt.Printf("Error searching similar in store: %v\n", errSS)
		return
	}

	chatMessages := convertChunksToChatMessage(scoredChunks, nlQuery)

	answer, errAns := newEmbedder.GenerateAnswer(ctx, chatMessages)
	if errAns != nil {
		fmt.Printf("Error generating answer: %v\n", errAns)
		return
	}

	fmt.Printf("\nAnswer:\n%v", answer)

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
