package collector

import (
	"go/ast"
)

// Check if a node has the given name.
func isNodeIdentWithName(n ast.Node, name string) bool {
	isNodeIdent := false

	ast.Inspect(n, func(n ast.Node) bool {
		if n == nil {
			return true
		}

		ident, isIdent := n.(*ast.Ident)
		if !isIdent {
			return true
		}

		if ident.Name == name || isNodeIdent {
			isNodeIdent = true
			return false
		}

		return true
	})

	return isNodeIdent
}
