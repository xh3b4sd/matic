package collector

import (
	"go/ast"

	taskqPkg "github.com/zyndiecate/taskq"
)

var (
	Verbosef = func(f string, v ...interface{}) {}
)

type Logger func(f string, v ...interface{})

func Configure(verboseLogger Logger) {
	Verbosef = verboseLogger
}

type File struct {
	// File path of a source code file.
	Path string

	// Go code in string form of a source code file.
	Code string

	// Variable name of the imported middleware package, if any.
	PkgImport string

	// *ast.File of the current go code.
	AstFile *ast.File

	// Serve information describing which routes the middleware server provides.
	ServeStmts []ServeStmt

	// Middleware information describing data and logic used of provided routes.
	//Middlewares []Middleware
}

type Ctx struct {
	WorkingDir string

	// Variable name of the created middleware server, if any. We assume there is
	// only one created middleware server. Maybe that is not true for all cases.
	ServerName string

	Files []File
}

type ClientCollectorI interface {
	// Generate a clients source code based on the given root path.
	GenerateClient(root string) error

	// Generate a api blueprint based on the given root path.
	ApiBlueprint(root string) (string, error)

	//CreateClientWithSourceCode(path, sourceCode string) error
	//CreateClientWithApiBlueprint(path, apiBlueprint string) error
}

type GoClientCollector struct{}

func NewGoClientCollector() ClientCollectorI {
	return &GoClientCollector{}
}

func (gcg *GoClientCollector) GenerateClient(wd string) error {
	// Create task context.
	ctx := &Ctx{
		WorkingDir: wd,
	}

	// Create a new queue.
	q := taskqPkg.NewQueue(ctx)

	// Run tasks.
	err := q.RunTasks(
		taskqPkg.InSeries(
			SourceCodeTask,
			PackageImportTask,
			ServerNameTask,
			ServeStmtTask,
			// find middlewares for each route
			// find possible responses for each route
		),
	)

	if err != nil {
		return Mask(err)
	}

	return nil
}

func (gcg *GoClientCollector) ApiBlueprint(root string) (string, error) {
	return "", nil
}
