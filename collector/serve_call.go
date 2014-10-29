package collector

import (
	"go/ast"
)

func ServeCallTask(v interface{}) error {
	Verbosef("Searching serve calls")

	for i, file := range ctx(v).Files {
		callExprs := callExprsByAstFile(file.AstFile)
		callExprs = filterServeCalls(file.PkgImport, callExprs)

		ctx(v).Files[i].ServeCalls = append(ctx(v).Files[i].ServeCalls, callExprs...)
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// private

// Search all call expressions of a file.
func callExprsByAstFile(af *ast.File) []*ast.CallExpr {
	callExprs := []*ast.CallExpr{}

	ast.Inspect(af, func(n ast.Node) bool {
		if n == nil {
			return true
		}

		callExpr, isCallExpr := n.(*ast.CallExpr)
		if !isCallExpr {
			return true
		}

		_, isSelExpr := callExpr.Fun.(*ast.SelectorExpr)
		if !isSelExpr {
			return true
		}

		callExprs = append(callExprs, callExpr)

		return true
	})

	return callExprs
}

// Remove all call expressions that are not `srv.Serve()`.
func filterServeCalls(pkg string, callExprs []*ast.CallExpr) []*ast.CallExpr {
	filteredExprs := []*ast.CallExpr{}

	for _, callExpr := range callExprs {
		if !isNodeIdentWithName(callExpr, "Serve") {
			continue
		}

		objDecl := callExpr.Fun.(*ast.SelectorExpr).X.(*ast.Ident).Obj.Decl

		if assign, isAssign := objDecl.(*ast.AssignStmt); isAssign {
			if !isNodeIdentWithName(assign, pkg) {
				continue
			}
		}

		if field, isField := objDecl.(*ast.Field); isField {
			if !isNodeIdentWithName(field, pkg) {
				continue
			}
		}

		filteredExprs = append(filteredExprs, callExpr)
	}

	return filteredExprs
}
