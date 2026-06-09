package directoryhandler

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/nikkisversion/astral/embedder"
	"github.com/nikkisversion/astral/filehandler"
	"github.com/nikkisversion/astral/store"
)

type DirectoryHandler struct {
	Path     string
	Embedder embedder.Embedder
	Store    store.Store
}

func New(path string, embedder embedder.Embedder, store store.Store) (*DirectoryHandler, error) {
	return &DirectoryHandler{
		Path:     path,
		Embedder: embedder,
		Store:    store,
	}, nil
}

func (h *DirectoryHandler) ProcessDirectory(ctx context.Context) error {

	numWorkers := runtime.NumCPU()
	fileChan := make(chan string, 100)
	errChan := make(chan error, numWorkers+1)

	wg := &sync.WaitGroup{}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {

			defer wg.Done()
			for path := range fileChan {

				fh, errFH := filehandler.New(path, h.Embedder, h.Store)
				if errFH != nil {
					fmt.Printf("Error creating handler for file: %v", path)
					errChan <- errFH
					return
				}

				errProcess := fh.ProcessFile(ctx)
				if errProcess != nil {
					fmt.Printf("Error processing file: %v", path)
					errChan <- errProcess
					return
				}
			}

		}()
	}

	err := filepath.WalkDir(h.Path, func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			return err
		}

		if !d.IsDir() && filepath.Ext(path) == ".go" {
			fileChan <- path
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error while parsing directory: %v", err)
		return err
	}

	close(fileChan)
	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return <-errChan
	}

	return nil

}
