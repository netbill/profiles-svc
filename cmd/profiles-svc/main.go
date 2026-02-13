package main

import (
	"os"

	"github.com/netbill/profiles-svc/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
