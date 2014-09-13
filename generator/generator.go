package generator

var (
	Verbosef = func(f string, v ...interface{}) {}
)

func Configure(verboseLogger VerboseLogger) {
	Verbosef = verboseLogger
}

type VerboseLogger func(f string, v ...interface{})

type Ctx struct {
	SourceCodeCtx SourceCodeCtx
}

type ClientGeneratorI interface {
	// Generate a clients source code based on the given root path.
	GenerateClient(root string) ([]SourceCode, error)

	// Generate a api blueprint based on the given root path.
	ApiBlueprint(root string) (string, error)

	//CreateClientWithSourceCode(path, sourceCode string) error
	//CreateClientWithApiBlueprint(path, apiBlueprint string) error
}
