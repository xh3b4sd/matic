package cli

import (
	"github.com/spf13/cobra"
)

var (
	globalFlags = struct {
		debug   bool
		verbose bool
	}{}

	maticCmd = &cobra.Command{
		Use:   "matic",
		Short: "Autogenerate clients for web-services",
		Long:  "Autogenerate clients for web-services",
		Run:   maticRun,
	}
)

func init() {
	maticCmd.PersistentFlags().BoolVarP(&globalFlags.debug, "debug", "d", false, "Print debug output")
	maticCmd.PersistentFlags().BoolVarP(&globalFlags.verbose, "verbose", "v", false, "Print verbose output")
}

func NewMaticCmd(version string) *cobra.Command {
	maticCmd.AddCommand(clientCmd)
	maticCmd.Execute()
	return maticCmd
}

func maticRun(cmd *cobra.Command, args []string) {
	cmd.Help()
}
