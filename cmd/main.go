package main

import (
	"context"
	"fmt"

	"github.com/nikkisversion/astral/embedder"
	"github.com/nikkisversion/astral/reader"
	"github.com/nikkisversion/astral/store"
)

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

	newEmbedder := embedder.NewOllamaEmbedder("http://localhost:11434", "nomic-embed-text")
	embeddings, err := newEmbedder.EmbedStrings(context.Background(), inputs)
	if err != nil {
		fmt.Printf("Error embedding strings: %v\n", err)
		return
	}

	newReader.UpdateEmbedding(embeddings)

	fmt.Printf("Collected %d chunks from %s\n", len(newReader.Collector.Chunks), filePath)
	fmt.Printf("First chunk content preview: %v\n", newReader.Collector.Chunks[0].Embedding)

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

}
