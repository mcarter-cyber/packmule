package cli

import (
	"fmt"

	"forgejo.eng.ultra-cis.com/devops/packmule/internal/consts"
	"forgejo.eng.ultra-cis.com/devops/packmule/internal/flags"
	"forgejo.eng.ultra-cis.com/devops/packmule/internal/version"
	"github.com/alecthomas/kong"
)

type VersionCMD struct {
	JSON bool `name:"json" short:"j" help:"set the output format to JSON"`
}

func (v *VersionCMD) Run(ctx *kong.Context, globals *flags.Global) error {
	vers := version.GetVersionInfo()
	vers.Name = ctx.Model.Name
	vers.Description = ctx.Model.Help
	vers.FontName = consts.FontName

	if v.JSON {
		out, err := vers.JSONString()
		if err != nil {
			return fmt.Errorf("unable to generate JSON from version info: %w", err)
		}
		fmt.Println(out)
	} else {
		fmt.Println(vers.String())
	}
	return nil
}
