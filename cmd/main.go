package main

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/nikkisversion/astral/collector"
	"github.com/nikkisversion/astral/file"
	"github.com/nikkisversion/astral/parser"
)

func main() {
	fmt.Println("Hello from Astral! LFG!")

	filePath := "testFile.go"

	fileReader, err := file.NewReadCloser(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer fileReader.Close()

	astFile, fset, err := parser.Parse(filePath, fileReader)
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
	}

	collector := collector.New(fset, astFile, filePath)
	ast.Walk(collector, astFile)

	PrintChunks(collector.Chunks)

}

func PrintChunks(chunks []collector.Chunk) {
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
