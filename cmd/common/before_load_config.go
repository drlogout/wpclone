package common

import (
	"croox/wpclone/config"
	"fmt"

	"github.com/urfave/cli/v2"
)

var BeforeLoadConfig = func(ctx *cli.Context) error {
	flags := ConfigFlags(ctx)

	cfg, err := config.Load(flags)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	SaveConfigToCTX(ctx, cfg)

	return nil
}
