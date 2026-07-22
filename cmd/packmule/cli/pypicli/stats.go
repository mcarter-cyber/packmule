package pypicli

import (
	"context"

	"forgejo.eng.ultra-cis.com/devops/packmule/internal/flags"
	"forgejo.eng.ultra-cis.com/devops/packmule/internal/pypi"
	"github.com/alecthomas/kong"
)

type StatsCMD struct {
}

func (s *StatsCMD) Run(ctx *kong.Context, globals *flags.Global, pypi_flags *flags.Pypi) error {
	store := pypi.NewFileStore("index.json")
	collector := pypi.NewCollector(store, pypi.CollectorConfig{
		FetchDetails: false, // just the index, no per-project crawl
	})
	if err := collector.Run(context.Background()); err != nil {
		globals.Logger.Errorf("%s", err)
	}
	return nil
}
