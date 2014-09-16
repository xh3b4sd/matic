package collector

import (
	"go/ast"
)

func ServerNameTask(ctx interface{}) error {
	Verbosef("Searching server name")

	// each source code of package import path
	for _, file := range ctx.(*Ctx).Files {
		if file.PkgImport == "" {
			continue
		}

		serverName := serverName(file.PkgImport, file.AstFile)
		if serverName != "" {
			ctx.(*Ctx).ServerName = serverName
			break
		}
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// private

// Find created server name, e.g. srv := srvPkg.NewServer(...)
func serverName(pkgName string, f *ast.File) string {
	srvName := ""

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Ident:
			if x.Obj != nil {
				switch y := x.Obj.Decl.(type) {
				case *ast.ValueSpec:
					// TODO
					// var srv = srvPkg.NewServer("127.0.0.1", "8080")
				case *ast.AssignStmt:
					// srv := srvPkg.NewServer("127.0.0.1", "8080")
					for i, rh := range y.Rhs {
						switch callExpr := rh.(type) {
						case *ast.CallExpr:
							switch callExpr.Fun.(type) {
							case *ast.SelectorExpr:
								funcExp := callExpr.Fun.(*ast.SelectorExpr).X.(*ast.Ident).Name // srvPkg
								funcSel := callExpr.Fun.(*ast.SelectorExpr).Sel.Name            // NewServer

								// Here we found the server name
								if funcExp == pkgName && funcSel == "NewServer" {
									srvName = y.Lhs[i].(*ast.Ident).Name
									return false
								}
							}
						}
					}
				}
			}
		}

		return true
	})

	return srvName
}
