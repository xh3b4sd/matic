package generator

import (
	"os"
	"path/filepath"
)

type SourceCode struct {
	// Path to the current source code file.
	Path string

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

func SourceCodeTask(ctx interface{}) error {
	root := ctx.(*Ctx).SourceCodeCtx.Root
	ext := ctx.(*Ctx).SourceCodeCtx.Ext

	Verbosef("Reading source code from root '%s'", root)

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
		ctx.(*Ctx).SourceCodeCtx.SourceCodeList = append(
			ctx.(*Ctx).SourceCodeCtx.SourceCodeList,
			SourceCode{Path: path, Code: code},
		)

		return nil
	})

	if err != nil {
		return Mask(err)
	}

	return nil
}
