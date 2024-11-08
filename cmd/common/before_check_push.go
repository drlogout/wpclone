package common

import (
	"croox/wpclone/pkg/message"
	"croox/wpclone/remote"

	"github.com/urfave/cli/v2"
)

var BeforeCheckPush = func(ctx *cli.Context) error {
	cfg := ConfigFromCTX(ctx)

	exists, err := remote.WPConfigExists(cfg)
	if err != nil {
		return err
	}

	cfg.Remote.WPConfigExists = exists
	SaveConfigToCTX(ctx, cfg)

	if exists && !cfg.Flags.PushWPConfig {
		return nil
	}

	if !exists {
		message.Info("wp-config.php not found on remote, pushing local wp-config.php")
	}

	if cfg.Flags.PushWPConfig {
		message.Info("Forcing push of wp-config.php")
	}

	if !cfg.RemoteDBCredentailsSet() {
		return message.Exit("Remote DB credentials not set in wpclone.yml")
	}

	return nil
}
