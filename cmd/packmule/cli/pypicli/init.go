package pypicli

import (
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"

	"forgejo.eng.ultra-cis.com/devops/packmule/internal/flags"
)

type InitCMD struct {
}

func (i *InitCMD) Run(ctx *kong.Context, globals *flags.Global, pypi_flags *flags.Pypi) error {
	globals.Logger.Infof("Initialising PyPi mirror.")

	mirror_dir := filepath.Join(globals.Config.Packmule.Directory, globals.Config.Pypi.Directory)
	mirror_dir, err := filepath.Abs(mirror_dir)
	if err != nil {
		return err
	}

	globals.Logger.Infof("Mirror directory: %s", mirror_dir)

	if _, err := os.Stat(mirror_dir); err != nil {
		globals.Logger.Infof("Creating target mirror directories.")
		if err := os.MkdirAll(mirror_dir, 0755); err != nil {
			return err
		}
	} else {
		globals.Logger.Infof("Mirror directories exist.")
	}

	globals.Logger.Infof("Collecting simple index.")

	// pypi := pypi.New(globals.Config.Pypi.PypiURL)

	return nil
}
