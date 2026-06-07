package filehandler

import (
	"context"
	"fmt"

	"github.com/nikkisversion/astral/embedder"
	"github.com/nikkisversion/astral/reader"
	"github.com/nikkisversion/astral/store"
)

type FileHandler struct {
	Reader   reader.Reader
	Embedder embedder.Embedder
	Store    store.Store
}

func New(filePath string, embedder embedder.Embedder, store store.Store) (*FileHandler, error) {

	newReader, err := reader.New(filePath)
	if err != nil {
		fmt.Printf("Failed to initialise new reader for file handler: %v", err.Error())
		return nil, err
	}

	return &FileHandler{
		Reader:   newReader,
		Embedder: embedder,
		Store:    store,
	}, nil
}

func (h *FileHandler) ProcessFile(ctx context.Context) error {

	errRead := h.Reader.Read()
	if errRead != nil {
		fmt.Printf("Error while reading file: %v", errRead.Error())
		return errRead
	}

	inputs := h.Reader.GenerateInputForEmbedding()

	embeddings, errEmbed := h.Embedder.EmbedStrings(ctx, inputs)
	if errEmbed != nil {
		fmt.Printf("Error while generating embeddings: %v", errEmbed.Error())
		return errEmbed
	}

	h.Reader.UpdateEmbedding(embeddings)

	errStore := h.Store.BatchUpsert(ctx, h.Reader.GetChunks())
	if errStore != nil {
		fmt.Printf("Error storing chunks: %v", errStore)
		return errStore
	}

	fmt.Println("File processed successfully")
	return nil

}
