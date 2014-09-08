package main

import (
	"github.com/zyndiecate/matic/cli"
)

var clientMaticVersion string

func main() {
	cli.NewCli(clientMaticVersion)
}
