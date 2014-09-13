package generator

import (
	generatorPkg "github.com/zyndiecate/matic/generator"
	taskqPkg "github.com/zyndiecate/taskq"
)

var (
	Verbosef = func(f string, v ...interface{}) {}
)

func Configure(verboseLogger generatorPkg.VerboseLogger) {
	Verbosef = verboseLogger
}

type GoClientGenerator struct{}

func NewGoClientGenerator() generatorPkg.ClientGeneratorI {
	return &GoClientGenerator{}
}

func (gcg *GoClientGenerator) GenerateClient(root string) ([]generatorPkg.SourceCode, error) {
	// Create task context.
	ctx := &generatorPkg.Ctx{
		SourceCodeCtx: generatorPkg.SourceCodeCtx{
			Ext:  "go",
			Root: root,
		},
	}

	// Create a new queue.
	q := taskqPkg.NewQueue(ctx)

	// Run tasks.
	err := q.RunTasks(
		taskqPkg.InSeries(
			generatorPkg.SourceCodeTask,
		),
	)

	if err != nil {
		return []generatorPkg.SourceCode{}, Mask(err)
	}

	return ctx.SourceCodeCtx.SourceCodeList, nil
}

func (gcg *GoClientGenerator) ApiBlueprint(root string) (string, error) {
	return "", nil
}
