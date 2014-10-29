package collector

import (
	"fmt"
)

func jsonResponseTask(v interface{}) error {
	for _, file := range ctx(v).Files {
		for mw, returnStmts := range file.ReturnStmts {
			fmt.Printf("%#v\n", mw)
			for _, returnStmt := range returnStmts {
				fmt.Printf("%#v\n", returnStmt)

				//      resCallArgs := returnStmt.Results[0].(*ast.CallExpr).Args
				//      resData := resCallArgs[0]
				//      resCode := resCallArgs[1]
				//      fmt.Printf("%#v\n", resData)
				//      fmt.Printf("%#v\n", resCode)
				//    }
				//  }
				//}

			}
		}
	}

	return nil
}
