package cli

import (
	pypicli "forgejo.eng.ultra-cis.com/devops/packmule/cmd/packmule/cli/pypicli"
	"forgejo.eng.ultra-cis.com/devops/packmule/internal/flags"
)

type PypiCMD struct {
	flags.Pypi

	Init  pypicli.InitCMD  `cmd:"" help:"initialise mirror directory"`
	Stats pypicli.StatsCMD `cmd:"" help:"Fetch Pypi statistics"`
}
