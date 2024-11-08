package remote

import (
	"croox/wpclone/config"
)

func CleanUp(cfg *config.Config) error {
	for _, file := range cfg.RemoteCleanupPaths() {
		if err := Shell(cfg, "rm", "-f", file); err != nil {
			return err
		}
	}

	if err := WPCli(cfg, "cache", "flush"); err != nil {
		return err
	}

	return nil
}
