package collector

import (
	"go/parser"
	"go/token"
	"path"
	"strings"
)

const (
	MiddlewarePkg = "github.com/catalyst-zero/middleware-server"
)

type PackageImport struct {
	// File path of the analysed source code file.
	FilePath string

	// Package name as imported and used in the source code.
	PkgName string
}

type PackageImportCtx struct {
	PackageImportList []PackageImport
}

// Find all files where the middleware package is imported.
func PackageImportTask(ctx interface{}) error {
	Verbosef("Searching imports of '%s'", MiddlewarePkg)

	for _, item := range ctx.(*Ctx).SourceCode.SourceCodeList {
		fset := token.NewFileSet()

		f, err := parser.ParseFile(fset, item.FilePath, item.Code, parser.ImportsOnly)
		if err != nil {
			return Mask(err)
		}

		// TODO what happens if a package is imported twice?
		filePath := item.FilePath
		pkgName := ""

		for _, spec := range f.Imports {
			// Just go further if we have the import spec of the middleware package.
			if !strings.Contains(spec.Path.Value, MiddlewarePkg) {
				continue
			}

			// If spec name if not nil, the package has an alternative name assigned.
			// E.g. here spec name would be 'srvPkg'
			//   import srvPkg "github.com/catalyst-zero/middleware-server"
			if spec.Name != nil {
				pkgName = spec.Name.Name
				break
			}

			// If there is no alternative package name defined, the package name is
			// the last element of the import path.
			pkgName = path.Base(spec.Path.Value)
			break
		}

		// Go ahead if we found no middleware package.
		if pkgName == "" {
			continue
		}

		ctx.(*Ctx).PackageImport.PackageImportList = append(
			ctx.(*Ctx).PackageImport.PackageImportList,
			PackageImport{FilePath: filePath, PkgName: pkgName},
		)
	}

	return nil
}
