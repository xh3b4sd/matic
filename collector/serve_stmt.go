package collector

import (
	"fmt"
	"go/ast"
	"strconv"
)

type MiddlewareExpr struct {
	FuncExpr     string
	FuncExprType string
	FuncSel      string
}

type ServeStmt struct {
	// Http method a route provides.
	Method string

	// Url path of the http route.
	Path string

	// Middleware specs used for a route.
	Middlewares []MiddlewareExpr
}

type IndirectServeStmt struct{}

type ServeStmtCtx struct {
	ServeStmtList []ServeStmt
	Indirect      IndirectServeStmt
}

func ServeStmtTask(ctx interface{}) error {
	Verbosef("Searching serve statements")

	serveStmtList := []ServeStmt{}
	serverName := ctx.(*Ctx).ServerName

	for _, packageImport := range ctx.(*Ctx).PackageImport.PackageImportList {
		astTree, err := astTreeByFile(packageImport.FilePath, ctx.(*Ctx).SourceCode.SourceCodeList)
		if err != nil {
			return Mask(err)
		}

		ast.Inspect(astTree, func(n ast.Node) bool {
			if n == nil {
				return true
			}

			switch callExpr := n.(type) {
			case *ast.CallExpr:
				selExpr := callExpr.Fun.(*ast.SelectorExpr)

				if selExpr.Sel.Name != "Serve" {
					return true
				}

				// Get serve statement of Serve() methods, where the callers expression
				// is directly assigned by the packages NewServer() method.
				if selExpr.Sel.Obj == nil &&
					selExpr.X.(*ast.Ident).Name == serverName.Name &&
					packageImport.FilePath == serverName.FilePath {

					serveStmt := serveStmtByCallExpr(callExpr)
					serveStmtList = append(serveStmtList, serveStmt)
					return true
				}

				// Get serve statement of Serve() methods, where the callers expression
				// is NOT directly assigned by the packages NewServer() method, but
				// referenced as method parameter.
				serveStmt(callExpr, packageImport.PkgName, func(serveStmt ServeStmt) {
					serveStmtList = append(serveStmtList, serveStmt)
				})
			}

			return true
		})
	}

	fmt.Printf("%#v\n", serveStmtList)

	ctx.(*Ctx).ServeStmt.ServeStmtList = serveStmtList

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

func serveStmt(callExpr *ast.CallExpr, pkgName string, cb func(ServeStmt)) {
	selExpr := callExpr.Fun.(*ast.SelectorExpr)

	switch ai := selExpr.X.(type) {
	case *ast.Ident:
		switch af := ai.Obj.Decl.(type) {
		case *ast.Field:
			switch aStarE := af.Type.(type) {
			case *ast.StarExpr:
				switch aSelE := aStarE.X.(type) {
				case *ast.SelectorExpr:
					if aSelE.X.(*ast.Ident).Name == pkgName && aSelE.Sel.Name == "Server" {
						serveStmt := serveStmtByCallExpr(callExpr)
						cb(serveStmt)
					}
				}
			}
		}
	}
}

func serveStmtByCallExpr(callExpr *ast.CallExpr) ServeStmt {
	serveStmt := ServeStmt{}

	for i, arg := range callExpr.Args {
		switch v := arg.(type) {
		case *ast.BasicLit:
			switch i {
			case 0:
				serveStmt.Method = unquote(v.Value)
			case 1:
				serveStmt.Path = unquote(v.Value)
			}
		case *ast.SelectorExpr:
			middlewareExpr := MiddlewareExpr{
				FuncExpr:     v.X.(*ast.Ident).Name, // v1
				FuncExprType: "V1",                  // TODO,
				FuncSel:      v.Sel.Name,            // MiddlewareOne
			}

			serveStmt.Middlewares = append(serveStmt.Middlewares, middlewareExpr)
		}
	}

	return serveStmt
}

func unquote(s string) string {
	if us, err := strconv.Unquote(s); err != nil {
		return s
	} else {
		return us
	}
}
