package main

import (
	cliPkg "github.com/zyndiecate/matic/cli"
)

var clientMaticVersion string

func main() {
	cliPkg.NewCli(clientMaticVersion)
}
