package reader

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
)

func Parse(filename string, rc io.ReadCloser) (*ast.File, *token.FileSet, error) {

	defer rc.Close()

	fset := token.NewFileSet()

	fileNode, err := parser.ParseFile(fset, filename, rc, parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}

	return fileNode, fset, nil
}
