package collector

import (
	"go/ast"
	"go/printer"
	"go/token"
	"strings"
)

type ChunkType string

const (
	FunctionType ChunkType = "function"
	StructType   ChunkType = "struct"
)

type Chunk struct {
	Type       ChunkType
	Name       string
	Content    string
	StartLine  int
	EndLine    int
	SourceFile string
}

func NewChunk(chunkType ChunkType, name, content string, startLine, endLine int, sourceFile string) *Chunk {
	return &Chunk{
		Type:       chunkType,
		Name:       name,
		Content:    content,
		StartLine:  startLine,
		EndLine:    endLine,
		SourceFile: sourceFile,
	}
}

func ExtractChunk(node ast.Node, fset *token.FileSet, name string, chunkType ChunkType, sourceFile string) (*Chunk, error) {

	var builder strings.Builder

	cfg := &printer.Config{
		Mode:     printer.TabIndent | printer.UseSpaces,
		Tabwidth: 8,
	}

	if err := cfg.Fprint(&builder, fset, node); err != nil {
		return nil, err
	}

	startPos := fset.Position(node.Pos())
	endPos := fset.Position(node.End())

	return NewChunk(chunkType, name, builder.String(), startPos.Line, endPos.Line, sourceFile), nil

}
