package cli

import (
	"github.com/spf13/cobra"

	collectorPkg "github.com/zyndiecate/matic/collector"
)

var (
	// Library.
	clientCollector collectorPkg.ClientCollectorI

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
	// Setup client collector.
	switch lang {
	case "go":
		clientCollector = collectorPkg.NewGoClientCollector()
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
	err := clientCollector.GenerateClient(root)
	if err != nil {
		ExitStderr(Mask(err))
	}
}
