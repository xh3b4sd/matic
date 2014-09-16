package collector

import (
	_ "fmt"
	"go/ast"
	"strconv"
)

type MiddlewareExpr struct {
	// The type a middleware is assigned to, e.g. V1 for v1.HelloWorldTwo
	Type string

	// The package where a middleware is defined, e.g. v1 for Foo
	Pkg string

	// The name of the middleware method, eg. HelloWorldTwo or Foo
	Name string
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

					// Get serve statement of Serve() methods, where the callers
					// expression is directly assigned by the packages NewServer()
					// method. That is, each Serve() statement in simple.go
					isServeDirectlyAssigned := selExpr.Sel.Obj == nil
					isServerServe := selExpr.X.(*ast.Ident).Name == ctx.(*Ctx).ServerName

					if isServeDirectlyAssigned && isServerServe {
						serveStmt := serveStmtByCallExpr(file.AstFile, callExpr)
						serveStmts = append(serveStmts, serveStmt)
						return true
					}

					// Get serve statement of Serve() methods, where the callers
					// expression is NOT directly assigned by the packages NewServer()
					// method, but referenced as method parameter. That is, each Serve()
					// statement in middleware/v1/middleware.go
					if callExprHasType(file.PkgImport, callExpr) && callExprHasName("Serve", callExpr) {
						serveStmt := serveStmtByCallExpr(file.AstFile, callExpr)
						serveStmts = append(serveStmts, serveStmt)
						return true
					}
				}
			}

			return true
		})

		ctx.(*Ctx).Files[i].ServeStmts = serveStmts
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// private

// Return true if e.g. strName is Serve and callExpr is srv.Serve(...)
func callExprHasName(strName string, callExpr *ast.CallExpr) bool {
	selExpr := starTypeSelExprByCallExpr(callExpr)

	if selExpr == nil {
		return false
	}

	if selExpr.Sel.Name == strName {
		return true
	}

	return false
}

// Return true if e.g. strType is srvPkg and callExpr is srv.Serve()
func callExprHasType(strType string, callExpr *ast.CallExpr) bool {
	selExpr := starTypeSelExprByCallExpr(callExpr)

	if selExpr == nil {
		return false
	}

	return selExpr.X.(*ast.Ident).Name == strType
}

func starExprBySelExpr(selExpr *ast.SelectorExpr) *ast.StarExpr {
	switch ai := selExpr.X.(type) {
	case *ast.Ident:
		if ai.Obj == nil {
			return nil
		}

		switch x := ai.Obj.Decl.(type) {
		case *ast.Field:
			switch aStarExpr := x.Type.(type) {
			case *ast.StarExpr:
				return aStarExpr
			}
		}
	}

	return nil
}

// Return true if e.g. strType is V1 and selExpr is v1.MiddlewareOne
func selExprHasType(strType string, selExpr *ast.SelectorExpr) bool {
	return selExprType(selExpr) == strType
}

// Return V1 if e.g. selExpr is v1.MiddlewareOne
func selExprType(selExpr *ast.SelectorExpr) string {
	starExpr := starExprBySelExpr(selExpr)

	if starExpr == nil {
		switch ai := selExpr.X.(type) {
		case *ast.Ident:
			if ai.Obj == nil {
				return ""
			}

			switch x := ai.Obj.Decl.(type) {
			case *ast.AssignStmt:
				return assignStmtType(x)
			}
		}
	}

	switch aIdent := starExpr.X.(type) {
	case *ast.Ident:
		return aIdent.Name
	}

	return ""
}

// Return v1Pkg if e.g. selExpr is v1.MiddlewareOne
func selExprPkg(selExpr *ast.SelectorExpr) string {
	selExprType := selExprType(selExpr)

	if selExprType == "" {
		return selExpr.X.(*ast.Ident).Name
	}

	return ""
}

func assignStmtType(assignStmt *ast.AssignStmt) string {
	if len(assignStmt.Lhs) != 1 {
		return ""
	}

	switch aIdent := assignStmt.Lhs[0].(type) {
	case *ast.Ident:
		switch assign := aIdent.Obj.Decl.(type) {
		case *ast.AssignStmt:
			if len(assign.Lhs) == 1 {
				switch aIdent2 := assign.Lhs[0].(type) {
				case *ast.Ident:
					switch assign2 := aIdent2.Obj.Decl.(type) {
					case *ast.AssignStmt:
						if len(assign2.Rhs) == 1 {
							switch unary := assign2.Rhs[0].(type) {
							case *ast.UnaryExpr:
								switch comp := unary.X.(type) {
								case *ast.CompositeLit:
									switch selExpr2 := comp.Type.(type) {
									case *ast.SelectorExpr:
										return selExpr2.Sel.Name
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return ""
}

func starTypeSelExprByCallExpr(callExpr *ast.CallExpr) *ast.SelectorExpr {
	selExpr := callExpr.Fun.(*ast.SelectorExpr)
	return starTypeSelExprBySelExpr(selExpr)
}

func starTypeSelExprBySelExpr(selExpr *ast.SelectorExpr) *ast.SelectorExpr {
	switch ai := selExpr.X.(type) {
	case *ast.Ident:
		switch af := ai.Obj.Decl.(type) {
		case *ast.Field:
			switch aStarExpr := af.Type.(type) {
			case *ast.StarExpr:
				switch aSelExpr := aStarExpr.X.(type) {
				case *ast.SelectorExpr:
					return aSelExpr
				}
			}
		}
	}

	return nil
}

func serveStmtByCallExpr(astFile *ast.File, callExpr *ast.CallExpr) ServeStmt {
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
				Type: selExprType(v),
				Pkg:  selExprPkg(v),
				Name: v.Sel.Name,
			}

			serveStmt.Middlewares = append(serveStmt.Middlewares, middlewareExpr)
		case *ast.Ident:
			// Middlewares of Serve() statements in e.g. middleware/v1/v1.go

			middlewareExpr := MiddlewareExpr{
				Type: "",
				Pkg:  astFile.Name.Name,
				Name: v.Name,
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
