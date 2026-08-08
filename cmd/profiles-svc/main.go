package main

import (
	"os"

	"github.com/netbill/profiles-svc/internal/build/cli"
)

func main() {
	cli.Run(os.Args)
}
