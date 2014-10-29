package collector

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
			ServeCallTask,
			ServeInfoTask,
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
