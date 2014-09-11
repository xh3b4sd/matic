package main

import (
	cliPkg "github.com/zyndiecate/matic/cli"
)

var projectVersion string

func main() {
	cliPkg.NewMaticCmd(projectVersion)
}
