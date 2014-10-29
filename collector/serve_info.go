package collector

import (
	"fmt"
	"go/ast"
	"strconv"
)

const (
  TypeAnnonymousMethod = "TypeAnnonymousMethod"
  TypeReceiverMethod = "TypeReceiverMethod"
)

type MiddlewareInfoI interface {
	MiddlewareType() string
}

type AnnonymousMethod struct {
  // The golang package name where the annonymous method is defined.
  Pkg string

  // The method name.
  Name string
}

func (am *AnnonymousMethod) MiddlewareType() string {
  return TypeAnnonymousMethod
}

type ReceiverMethod struct {
  // The receiver type.
  Type string

  // The method name.
  Name string
}

func (rm *ReceiverMethod) MiddlewareType() string {
  return TypeReceiverMethod
}

type ServeInfo struct {
	Method          string
	Path            string
	MiddlewareInfos []MiddlewareInfoI
}

func ServeInfoTask(v interface{}) error {
	for i, file := range ctx(v).Files {
    fmt.Printf("\n")

		for _, serveCall := range file.ServeCalls {
			serveInfo := serveInfoByServeCall(file.AstFile.Name.Name, serveCall)
			ctx(v).Files[i].ServeInfos = append(ctx(v).Files[i].ServeInfos, serveInfo)
		}
	}

	fmt.Printf("%#v\n", ctx(v).Files[2].ServeInfos[0])

	return nil
}

// `pkgName` is the name of the golang package of the current file.
func serveInfoByServeCall(pkgName string, serveCall *ast.CallExpr) *ServeInfo {
	serveInfo := &ServeInfo{
		Method: basicLiteralName(serveCall.Args[0]),
		Path:   basicLiteralName(serveCall.Args[1]),
	}

	for _, arg := range serveCall.Args[2:] {
    // Find middleware info by ident, e.g. Foo()
    if ident, isIdent := arg.(*ast.Ident); isIdent {
      annonymousMethod := &AnnonymousMethod{
        Pkg: pkgName,
        Name: ident.Name,
      }

      serveInfo.MiddlewareInfos = append(serveInfo.MiddlewareInfos, annonymousMethod)
      continue
    }

    // Find middleware info by selector expressions, e.g. v1.Foo()
    if selExpr, isSelExpr := arg.(*ast.SelectorExpr); isSelExpr {
      receiverMethod := &ReceiverMethod{
        Type: middlewareInfoType(selExpr.X.(*ast.Ident)),
        Name: selExpr.Sel.Name,
      }

      serveInfo.MiddlewareInfos = append(serveInfo.MiddlewareInfos, receiverMethod)
      continue
    }
	}

	return serveInfo
}

// TODO
func middlewareInfoType(ident *ast.Ident) string {
  miType := ""

  if ident.Obj == nil {
    // TODO v1Pkg.Foo()
    // check package import for package name
    // create annonymus method middleware info
  } else {
    fmt.Printf("%#v\n", ident.Obj.Decl)

    ast.Inspect(ident, func(n ast.Node) bool {
      return true
    })
  }

  return miType
}

func basicLiteralName(arg ast.Node) string {
	if basicLiteral, isBasicLiteral := arg.(*ast.BasicLit); isBasicLiteral {
		return unquote(basicLiteral.Value)
	}

	return ""
}

func unquote(s string) string {
	if us, err := strconv.Unquote(s); err != nil {
		return s
	} else {
		return us
	}
}
