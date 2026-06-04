package reader

import "go/ast"

type Reader struct {
	SourceFile string
	Collector  *Collector
}

func New(filePath string) (*Reader, error) {
	return &Reader{SourceFile: filePath}, nil
}

func (r *Reader) Read() error {

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

func (r *Reader) GenerateInputForEmbedding() []string {
	return GetInputForEmbedding(r.Collector.Chunks)
}

func (r *Reader) UpdateEmbedding(embeddings [][]float32) {
	UpdateChunksWithEmbeddings(r.Collector.Chunks, embeddings)
}
