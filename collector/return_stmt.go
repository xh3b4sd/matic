package collector

//
//import (
//  "fmt"
//	"go/ast"
//)
//
//type ReturnWalker struct {
//  Middleware string
//  ReturnStmts []*ast.ReturnStmt
//  PushReturnStmt func(string, *ast.ReturnStmt)
//}
//
//func (rw *ReturnWalker) Visit(n ast.Node) ast.Visitor {
//  if n == nil {
//    return nil
//  }
//
//  if returnStmt, ok := n.(*ast.ReturnStmt); ok {
//    rw.PushReturnStmt(rw.Middleware, returnStmt)
//  }
//
//  return rw
//}
//
//func ReturnStmtTask(v interface{}) error {
//	Verbosef("Searching return statements")
//
//	ctx(v).FileMiddlewarePair(func(file *File, mw *MiddlewareExpr) {
//    returnWalker := &ReturnWalker{
//      Middleware: mw.Name,
//      PushReturnStmt: file.PushReturnStmt,
//    }
//
//		ast.Inspect(file.AstFile, func(n ast.Node) bool {
//			fd := middlewareFuncDecl(n, mw)
//      if fd == nil {
//        return true
//      }
//
//      // The middleware is a receiver method.
//      if fd.Recv != nil {
//        // Pointer receiver.
//        if se, ok := fd.Recv.List[0].Type.(*ast.StarExpr); ok {
//          if se.X.(*ast.Ident).Name == mw.Type {
//              ast.Walk(returnWalker, fd)
//          }
//        }
//      }
//
//			return true
//		})
//	})
//
//	return nil
//}
//
//func middlewareFuncDecl(n ast.Node, mw *MiddlewareExpr) *ast.FuncDecl {
//  zeroVal := &ast.FuncDecl{}
//
//  if n == nil {
//    return zeroVal
//  }
//
//  // We found a function declaration.
//  fd, ok := n.(*ast.FuncDecl)
//  if !ok {
//    return zeroVal
//  }
//
//  // We found the middleware we searched for.
//  if fd.Name.Name != mw.Name {
//    return zeroVal
//  }
//
//  return fd
//}
//
//  // The middleware calls ctx.Next(), nothing to do here.
//  //func isCallingNext(returnStmt *ast.ReturnStmt) bool {
//  //  _, returnXIsIdent := returnStmt.X.(*ast.Ident)
//
//  //  if returnXIsIdent && returnStmt.Sel.Name == "Next" {
//  //    return true
//  //  }
//
//  //  return false
//  //}
//
//              // TODO user Walk to walk through middleware function bodies and find return statements.
//
//              //.(*ast.ReturnStmt).Results[0].(*ast.CallExpr).Fun.(*ast.SelectorExpr)
//              //returnStmt, ok := bodyStmt.(*ast.ReturnStmt)
//              //if !ok {
//              //  continue
//              //}
//
//              //if isCallingNext(returnStmt) {
//              //  continue
//              //}
//
//              //ctxStmt := returnStmt.X.(*ast.SelectorExpr).X.(*ast.Ident).Obj.Decl.(*ast.Field).Type.(*ast.StarExpr).X.(*ast.SelectorExpr)
//              //if ctxStmt.X.(*ast.Ident).Name == file.PkgImport && ctxStmt.Sel.Name == "Context" {
//              //  if returnStmt.X.(*ast.SelectorExpr).Sel.Name == "Response" {
//              //    if returnStmt.Sel.Name == "Json" {
//              //      file.PushReturnStmt(mw.Name, returnStmt)
//              //    }
//              //  }
//              //}
