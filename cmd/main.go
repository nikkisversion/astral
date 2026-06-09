package main

import (
	"context"
	"fmt"

	"github.com/nikkisversion/astral/directoryhandler"
	"github.com/nikkisversion/astral/embedder"
	"github.com/nikkisversion/astral/queryhandler"
	"github.com/nikkisversion/astral/store"
)

/* Steps:
 1. Run Ollama application and on terminal, pull the nomic-embed-text model using:
 `ollama pull nomic-embed-text` and
 `ollama pull qwen2.5-coder:1.5b`
 2. Run the docker command to start qdrant: docker run -p 6333:6333 -p 6334:6334 \
    -v "$(pwd)/qdrant_storage:/qdrant/storage:z" \
    qdrant/qdrant
3. Run this main.go file using `go run cmd/main.go` and see the results in the terminal.
4. Open 'http://localhost:6333/dashboard' in your browser to see the QDrant dashboard and verify that the collection and points have been created.
*/

func main() {
	fmt.Println("Hello from Astral! LFG!")

	//filePath := "/Users/nikitarai/Documents/projects/astral/reader/reader.go"
	path := "/Users/nikitarai/Documents/projects/astral/reader"
	ctx := context.Background()

	e := embedder.NewOllamaEmbedder("http://localhost:11434", "nomic-embed-text", "qwen2.5-coder:1.5b")

	vs, err := store.NewQDrantStore(ctx, "code-chunks", 768)
	if err != nil {
		fmt.Printf("Error creating vector store: %v\n", err)
		return
	}

	dirHandler, errDH := directoryhandler.New(path, e, vs)
	if errDH != nil {
		fmt.Printf("Error creating DH: %v", errDH)
		return
	}

	errPD := dirHandler.ProcessDirectory(ctx)
	if errPD != nil {
		fmt.Printf("Error processing directory: %v", errPD)
		return
	}

	fmt.Println("Directory processed successfully!")

	nlQuery := "What is the Visit function doing? Walk me through the steps."
	qh := queryhandler.New(e, vs)
	answer, errAns := qh.Handle(ctx, nlQuery)
	if errAns != nil {
		return
	}

	fmt.Printf("\nAnswer:\n%v", answer)

}
