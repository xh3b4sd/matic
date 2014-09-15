package collector

import (
	"path"
	"strings"
)

// Find all files where the middleware package is imported.
func PackageImportTask(ctx interface{}) error {
	Verbosef("Searching package imports")

	for i, file := range ctx.(*Ctx).Files {
		pkgImport := ""

		for _, is := range file.AstFile.Imports {
			// Continue if we not have the import spec of the middleware package.
			if !strings.Contains(is.Path.Value, "github.com/catalyst-zero/middleware-server") {
				continue
			}

			// If import spec name if not nil, the package has an alternative name
			// assigned.  E.g. here spec name would be 'srvPkg'
			//   import srvPkg "github.com/catalyst-zero/middleware-server"
			if is.Name != nil {
				pkgImport = is.Name.Name
				break
			}

			// If there is no alternative package name defined, the package name is
			// the last element of the import path.
			pkgImport = path.Base(is.Path.Value)
			break
		}

		// Continue if we found no middleware package.
		if pkgImport == "" {
			continue
		}

		ctx.(*Ctx).Files[i].PkgImport = pkgImport
	}

	return nil
}
