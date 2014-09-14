package generator

import (
	_ "fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

func skipFile(ext, path string, info os.FileInfo) bool {
	// Skip directories.
	if info.IsDir() {
		return true
	}

	// Skip none go files.
	if filepath.Ext(path) != "."+ext {
		return true
	}

	return false
}

func readFile(path string) (string, error) {
	reader, err := os.Open(path)
	if err != nil {
		return "", Mask(err)
	}

	byteSlice, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", Mask(err)
	}

	return string(byteSlice), nil
}

func astTreeByFile(filePath string, sourceCodeList []SourceCode) (*ast.File, error) {
	code := ""

	for _, sourceCode := range sourceCodeList {
		if filePath == sourceCode.FilePath {
			code = sourceCode.Code
			break
		}
	}

	// That should never happen.
	if code == "" {
		return nil, Mask(ErrSourceCodeNotFoundByFilePath)
	}

	// Create an ast tree for the given piece of source code.
	fset := token.NewFileSet()
	astTree, err := parser.ParseFile(fset, filePath, code, 0)
	if err != nil {
		return nil, Mask(err)
	}

	return astTree, nil
}

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
