package reader

import (
	"fmt"
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
	Embedding  []float32
}

func NewChunk(chunkType ChunkType, name, content string, startLine, endLine int, sourceFile string) *Chunk {
	return &Chunk{
		Type:       chunkType,
		Name:       name,
		Content:    content,
		StartLine:  startLine,
		EndLine:    endLine,
		SourceFile: sourceFile,
		Embedding:  nil, // Placeholder for future embedding data
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

func GetInputForEmbedding(chunks []Chunk) []string {

	inputs := make([]string, len(chunks))

	for i, chunk := range chunks {
		inputs[i] = chunk.Content
	}

	return inputs
}

func UpdateChunksWithEmbeddings(chunks []Chunk, embeddings [][]float32) []Chunk {

	for i := range chunks {
		chunks[i].Embedding = embeddings[i]
	}
	return chunks

}

func PrintChunks(chunks []Chunk) {
	fmt.Printf("\n%-25s | %-10s | %-10s | %s\n", "NAME", "TYPE", "LINES", "CONTENT PREVIEW")
	fmt.Println("----------------------------------------------------------------------------------------------------")
	for _, chunk := range chunks {
		// 1. Split content by lines to isolate and drop inline comments
		lines := strings.Split(chunk.Content, "\n")
		var cleanLines []string

		for _, line := range lines {
			trimmed := strings.TrimSpace(line)

			// Skip entirely commented lines within the body
			if strings.HasPrefix(trimmed, "//") {
				continue
			}

			// If a line has code followed by an inline comment, strip the comment portion
			if idx := strings.Index(line, "//"); idx != -1 {
				line = line[:idx]
			}

			cleanLines = append(cleanLines, line)
		}

		// 2. Re-join and collapse all spaces/newlines into a single line preview
		preview := strings.Join(cleanLines, " ")
		preview = strings.Join(strings.Fields(preview), " ") // Normalizes whitespace noise

		if len(preview) > 50 {
			preview = preview[:47] + "..."
		}

		fmt.Printf("%-25s | %-10s | %d-%-6d | %s\n",
			chunk.Name,
			chunk.Type,
			chunk.StartLine,
			chunk.EndLine,
			preview,
		)
	}
	fmt.Println()
}
