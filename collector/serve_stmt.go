package collector

import (
	"fmt"
	"go/ast"
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
