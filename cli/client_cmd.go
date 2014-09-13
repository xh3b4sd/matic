package cli

import (
	"github.com/spf13/cobra"

	generatorPkg "github.com/zyndiecate/matic/generator"
	goPkg "github.com/zyndiecate/matic/generator/go"
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
	clientCmd.PersistentFlags().StringVarP(&lang, "lang", "l", "go", "Which language to use to generate client")
}

func clientRun(cmd *cobra.Command, args []string) {
	// Setup client generator.
	switch lang {
	case "go":
		clientGenerator = goPkg.NewGoClientGenerator()
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
	sourceCodeList, err := clientGenerator.GenerateClient(root)
	if err != nil {
		ExitStderr(Mask(err))
	}

	for _, item := range sourceCodeList {
		Verbosef("### %s ####", item.Path)
		Verbosef(item.Code)
	}
}
