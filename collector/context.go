package collector

import (
	"go/ast"
)

type File struct {
	// File path of a source code file.
	Path string

	// Go code in string form of a source code file.
	Code string

	// Variable name of the imported middleware package, if any.
	PkgImport string

	// *ast.File of the current go code.
	AstFile *ast.File

	// Call expressions of all serve calls.
	ServeCalls []*ast.CallExpr

	// Information of all serve calls.
	ServeInfos []*ServeInfo

	ReturnStmts map[string][]*ast.ReturnStmt
}

type Ctx struct {
	WorkingDir string

	// Variable name of the created middleware server, if any. We assume there is
	// only one created middleware server. Maybe that is not true for all cases.
	ServerName string

	Files []*File
}

func ctx(v interface{}) *Ctx {
	return v.(*Ctx)
}
