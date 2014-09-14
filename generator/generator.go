package generator

import (
	taskqPkg "github.com/zyndiecate/taskq"
)

var (
	Verbosef = func(f string, v ...interface{}) {}
)

type Logger func(f string, v ...interface{})

func Configure(verboseLogger Logger) {
	Verbosef = verboseLogger
}

type Ctx struct {
	SourceCode    SourceCodeCtx
	PackageImport PackageImportCtx
	ServerName    ServerNameCtx
	ServeStmt     ServeStmtCtx
}

type ClientGeneratorI interface {
	// Generate a clients source code based on the given root path.
	GenerateClient(root string) ([]SourceCode, error)

	// Generate a api blueprint based on the given root path.
	ApiBlueprint(root string) (string, error)

	//CreateClientWithSourceCode(path, sourceCode string) error
	//CreateClientWithApiBlueprint(path, apiBlueprint string) error
}

type GoClientGenerator struct{}

func NewGoClientGenerator() ClientGeneratorI {
	return &GoClientGenerator{}
}

func (gcg *GoClientGenerator) GenerateClient(root string) ([]SourceCode, error) {
	// Create task context.
	ctx := &Ctx{
		SourceCode: SourceCodeCtx{
			Ext:  "go",
			Root: root,
		},
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
		return []SourceCode{}, Mask(err)
	}

	//for _, item := range ctx.PackageImport.PackageImportList {
	//	Verbosef("### %s ####", item.FilePath)
	//	Verbosef(item.PkgName)
	//}

	return ctx.SourceCode.SourceCodeList, nil
}

func (gcg *GoClientGenerator) ApiBlueprint(root string) (string, error) {
	return "", nil
}
