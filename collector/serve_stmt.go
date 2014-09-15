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

func ServeStmtTask(ctx interface{}) error {
	Verbosef("Searching serve statements")

	for i, file := range ctx.(*Ctx).Files {
		if file.PkgImport == "" {
			continue
		}

		serveStmts := []ServeStmt{}

		ast.Inspect(file.AstFile, func(n ast.Node) bool {
			if n == nil {
				return true
			}

			switch callExpr := n.(type) {
			case *ast.CallExpr:
				switch selExpr := callExpr.Fun.(type) {
				case *ast.SelectorExpr:
					if selExpr.Sel.Name != "Serve" {
						return true
					}

					// Get serve statement of Serve() methods, where the callers expression
					// is directly assigned by the packages NewServer() method.
					isServeDirectlyAssigned := selExpr.Sel.Obj == nil
					isServerServe := selExpr.X.(*ast.Ident).Name == ctx.(*Ctx).ServerName
					// TODO && packageImport.FilePath == serverName.FilePath
					if isServeDirectlyAssigned && isServerServe {
						serveStmt := serveStmtByCallExpr(callExpr)
						serveStmts = append(serveStmts, serveStmt)
						return true
					}

					// Get serve statement of Serve() methods, where the callers expression
					// is NOT directly assigned by the packages NewServer() method, but
					// referenced as method parameter.
					serveStmt(callExpr, file.PkgImport, func(serveStmt ServeStmt) {
						serveStmts = append(serveStmts, serveStmt)
					})
				}
			}

			return true
		})

		ctx.(*Ctx).Files[i].ServeStmts = serveStmts
	}

	fmt.Printf("%#v\n", ctx.(*Ctx).Files)

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// private

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
