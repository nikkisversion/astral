package reader

import "go/ast"

type Reader interface {
	Read() error
	GetChunks() []Chunk
	GenerateInputForEmbedding() []string
	UpdateEmbedding(embeddings [][]float32)
}

type GoReader struct {
	SourceFile string
	Collector  *Collector
}

func New(filePath string) (*GoReader, error) {
	return &GoReader{SourceFile: filePath}, nil
}

func (r *GoReader) Read() error {

	fileReader, err := NewReadCloser(r.SourceFile)
	if err != nil {
		return err
	}

	node, fset, err := Parse(r.SourceFile, fileReader)
	if err != nil {
		return err
	}

	collector := NewCollector(fset, node, r.SourceFile)
	ast.Walk(collector, node)

	r.Collector = collector

	return nil

}

func (r *GoReader) GetChunks() []Chunk {
	return r.Collector.Chunks
}

func (r *GoReader) GenerateInputForEmbedding() []string {
	return GetInputForEmbedding(r.Collector.Chunks)
}

func (r *GoReader) UpdateEmbedding(embeddings [][]float32) {
	UpdateChunksWithEmbeddings(r.Collector.Chunks, embeddings)
}
