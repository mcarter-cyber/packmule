package flags

import (
	"forgejo.eng.ultra-cis.com/devops/packmule/internal/config"
	"forgejo.eng.ultra-cis.com/devops/packmule/pkg/log"
)

type Global struct {
	Logger log.Logger    `kong:"-"`
	Config config.Config `kong:"-"`

	ConfigFile string `short:"c" help:"configuration file" default:"${config_file}" placeholder:"${config_file}"`
	Debug      bool   `help:"Enable debug mode"`
}
