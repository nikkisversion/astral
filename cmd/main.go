package main

import (
	"context"
	"fmt"

	"github.com/nikkisversion/astral/embedder"
	"github.com/nikkisversion/astral/reader"
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

}
