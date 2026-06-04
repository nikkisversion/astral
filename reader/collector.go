package reader

import (
	"go/ast"
	"go/token"
)

type Collector struct {
	FSet       *token.FileSet
	Node       ast.Node
	Chunks     []Chunk
	SourceFile string
}

func NewCollector(fset *token.FileSet, node ast.Node, sourceFile string) *Collector {
	return &Collector{
		FSet:       fset,
		Node:       node,
		Chunks:     []Chunk{},
		SourceFile: sourceFile,
	}
}

func (c *Collector) Visit(node ast.Node) ast.Visitor {

	if node == nil {
		return c
	}

	switch n := node.(type) {
	case *ast.FuncDecl:
		chunk, err := ExtractChunk(node, c.FSet, n.Name.Name, FunctionType, c.SourceFile)
		if err != nil {
			return c
		}
		c.Chunks = append(c.Chunks, *chunk)
		return nil
	case *ast.GenDecl:
		if n.Tok == token.TYPE {
			for _, spec := range n.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					chunk, err := ExtractChunk(node, c.FSet, typeSpec.Name.Name, StructType, c.SourceFile)
					if err != nil {
						return c
					}
					c.Chunks = append(c.Chunks, *chunk)
				}
			}
			return nil
		}
	}

	return c
}
