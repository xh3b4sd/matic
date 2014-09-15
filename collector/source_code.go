package collector

import (
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Find all important source code files.
func SourceCodeTask(ctx interface{}) error {
	Verbosef("Reading source code")

	err := filepath.Walk(ctx.(*Ctx).WorkingDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return Mask(err)
		}

		// Skip uninteresting files.
		if skipFile(path, info) {
			return nil
		}

		// Read file.
		code, err := readFile(path)
		if err != nil {
			return Mask(err)
		}

		// Create an ast tree for the given piece of source code.
		fset := token.NewFileSet()
		astFile, err := parser.ParseFile(fset, path, code, 0)
		if err != nil {
			return Mask(err)
		}

		// Create context file.
		file := File{
			Path:    path,
			Code:    code,
			AstFile: astFile,
		}

		ctx.(*Ctx).Files = append(ctx.(*Ctx).Files, file)

		return nil
	})

	if err != nil {
		return Mask(err)
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// private

func skipFile(path string, info os.FileInfo) bool {
	// Skip directories.
	if info.IsDir() {
		return true
	}

	// Skip none go files.
	if filepath.Ext(path) != ".go" {
		return true
	}

	return false
}

func readFile(path string) (string, error) {
	reader, err := os.Open(path)
	if err != nil {
		return "", Mask(err)
	}

	byteSlice, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", Mask(err)
	}

	return string(byteSlice), nil
}
