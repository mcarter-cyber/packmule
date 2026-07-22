package main

import (
	"os"

	"forgejo.eng.ultra-cis.com/devops/packmule/cmd/packmule/cli"
)

func main() {

	kongCtx, err := cli.New()
	if err != nil {
		os.Exit(1)
	}

	if err := kongCtx.Run(); err != nil {
		os.Exit(1)
	}
}
