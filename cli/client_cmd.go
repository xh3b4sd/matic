package cli

import (
	"github.com/spf13/cobra"

	generatorPkg "github.com/zyndiecate/matic/generator"
)

var (
	// Library.
	clientGenerator generatorPkg.ClientGeneratorI

	// CLI.
	clientCmd = &cobra.Command{
		Use:   "client",
		Short: "Autogenerate client for web-service",
		Long:  "Autogenerate client for web-service",
		Run:   clientRun,
	}

	// Flags.
	lang string
	root string
)

func init() {
	clientCmd.Flags().StringVarP(&lang, "lang", "l", "go", "Which language to use to generate client")
}

func clientRun(cmd *cobra.Command, args []string) {
	// Setup client generator.
	switch lang {
	case "go":
		clientGenerator = generatorPkg.NewGoClientGenerator()
	default:
		cmd.Help()
		ExitStderr(ErrWrongInput)
	}

	// Setup source code root.
	switch len(args) {
	case 0:
		root = "."
	case 1:
		root = args[0]
	default:
		cmd.Help()
		ExitStderr(ErrWrongInput)
	}

	// Generate client based on given directory.
	_, err := clientGenerator.GenerateClient(root)
	if err != nil {
		ExitStderr(Mask(err))
	}
}
