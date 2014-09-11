package cli

import (
	"github.com/spf13/cobra"

	maticPkg "github.com/zyndiecate/matic/src"
)

var (
	clientCmd = &cobra.Command{
		Use:   "client",
		Short: "Autogenerate client for web-service",
		Long:  "Autogenerate client for web-service",
		Run:   clientRun,
	}
)

func clientRun(cmd *cobra.Command, args []string) {
	switch len(args) {
	case 0:
		// Generate client based on current directory.
		s, err := maticPkg.SourceCode(".")
		if err != nil {
			ExitStderr(Mask(err))
		}

		ExitStdoutf(s)
	case 1:
		// Generate client based on given directory.
		s, err := maticPkg.SourceCode(args[0])
		if err != nil {
			ExitStderr(Mask(err))
		}

		ExitStdoutf(s)
	}
}
