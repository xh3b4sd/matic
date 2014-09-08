package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	maticPkg "github.com/zyndiecate/matic/src"
)

func NewCli(version string) *cobra.Command {
	cli := &cobra.Command{
		Use:   "matic",
		Short: "Autogenerate clients for web-services",
		Long:  "Autogenerate clients for web-services",
		Run:   maticRun,
	}

	cli.Execute()

	return cli
}

func maticRun(cmd *cobra.Command, args []string) {
	fmt.Printf("%#v\n", args)
	fmt.Printf("%#v\n", maticPkg.Foo())
}
