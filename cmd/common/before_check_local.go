package common

import (
	"croox/wpclone/pkg/util"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

var BeforeCheckLocalFolder = func(ctx *cli.Context) error {
	if ctx.Bool("force") {
		return nil
	}

	cfg := ConfigFromCTX(ctx)

	if !util.FileExists(cfg.LocalPath()) {
		if err := os.MkdirAll(cfg.LocalPath(), 0755); err != nil {
			return err
		}
	}

	empty, err := util.FolderEmpty(cfg.LocalPath())
	if err != nil {
		return err
	}

	if !empty && !cfg.Flags.SkipRsync {
		label := fmt.Sprintf("❓ Local folder %s is not empty, do you want to proceed?", cfg.LocalPath())
		ok := util.YesNoPrompt(label, false)
		if !ok {
			return fmt.Errorf("Aborted")
		}
	}

	return nil
}
