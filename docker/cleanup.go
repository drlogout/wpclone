package docker

import (
	"croox/wpclone/config"
	"os"
)

func CleanUp(cfg *config.Config) error {
	for _, file := range cfg.LocalCleanupPaths() {
		if err := os.RemoveAll(file); err != nil {
			return err
		}
	}

	if err := WPCli(cfg, "cache", "flush"); err != nil {
		return err
	}

	return nil
}
