package main

import (
	"os"

	"github.com/netbill/profiles-svc/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
