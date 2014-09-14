package generator

import (
	"os"
	"path/filepath"
)

type SourceCode struct {
	// File path to the current source code file.
	FilePath string

	// Content of the current source code file.
	Code string
}

type SourceCodeCtx struct {
	// Extension of source code files to analyse, to generate a client.
	Ext string

	// The root path source code will be collected from.
	Root string

	// List of found source code.
	SourceCodeList []SourceCode
}

// Find all important source code files.
func SourceCodeTask(ctx interface{}) error {
	root := ctx.(*Ctx).SourceCode.Root
	ext := ctx.(*Ctx).SourceCode.Ext

	Verbosef("Reading source code from '%s'", root)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return Mask(err)
		}

		// Skip uninteresting files.
		if skipFile(ext, path, info) {
			return nil
		}

		// Read file.
		code, err := readFile(path)
		if err != nil {
			return Mask(err)
		}

		// Fill context.
		ctx.(*Ctx).SourceCode.SourceCodeList = append(
			ctx.(*Ctx).SourceCode.SourceCodeList,
			SourceCode{FilePath: path, Code: code},
		)

		return nil
	})

	if err != nil {
		return Mask(err)
	}

	return nil
}
