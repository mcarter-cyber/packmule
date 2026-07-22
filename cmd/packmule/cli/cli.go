package cli

import (
	"os"

	"forgejo.eng.ultra-cis.com/devops/packmule/internal/config"
	"forgejo.eng.ultra-cis.com/devops/packmule/internal/flags"
	"forgejo.eng.ultra-cis.com/devops/packmule/pkg/log"
	"github.com/alecthomas/kong"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

type CLI struct {
	flags.Global

	Version VersionCMD `cmd:"" aliases:"v" help:"show version information"`
	Pypi    PypiCMD    `cmd:"" help:"interact with the pypi mirror"`
}

func New() (*kong.Context, error) {
	cli := CLI{}

	kongCtx := kong.Parse(&cli,
		kong.Name("packmule"),
		kong.Description("Mirror packages for offline usage."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact:             true,
			NoExpandSubcommands: true,
		}),
		kong.Bind(&cli.Global),
		kong.Bind(&cli.Pypi.Pypi),
		kong.Vars{
			"config_file": "./packmule.yaml",
		},
	)

	// Initialize logger after parsing flags
	cli.Logger = log.NewLogger(os.Stdout)

	// Set debug level if debug flag is enabled
	if cli.Debug {
		cli.Logger.SetLevel("debug")
	} else {
		cli.Logger.SetLevel("info")
	}

	k := *koanf.New(".")

	// Load Default configuration
	if err := k.Load(structs.Provider(&config.DefaultConfig, "koanf"), nil); err != nil {
		cli.Logger.Errorf("error loading config: %v", err)
		return nil, err
	}
	// Load Configuration file - configs will be merged
	if err := k.Load(file.Provider(cli.ConfigFile), yaml.Parser()); err != nil {
		cli.Logger.Errorf("error loading config: %v", err)
		return nil, err
	}

	k.Unmarshal("", &cli.Config)

	return kongCtx, nil
}
