package collector

import (
	_ "fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func astTreeByFile(filePath string, sourceCodeList []SourceCode) (*ast.File, error) {
	code := ""

	for _, sourceCode := range sourceCodeList {
		if filePath == sourceCode.FilePath {
			code = sourceCode.Code
			break
		}
	}

	// That should never happen.
	if code == "" {
		return nil, Mask(ErrSourceCodeNotFoundByFilePath)
	}

	// Create an ast tree for the given piece of source code.
	fset := token.NewFileSet()
	astTree, err := parser.ParseFile(fset, filePath, code, 0)
	if err != nil {
		return nil, Mask(err)
	}

	return astTree, nil
}
