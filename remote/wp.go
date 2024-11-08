package remote

import (
	"croox/wpclone/config"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func WPInstalled(cfg *config.Config) (bool, error) {
	err := WPCli(cfg, "core", "is-installed")
	if err != nil {
		if _, ok := err.(*ssh.ExitError); !ok {
			return false, fmt.Errorf("Failed to check if WordPress is installed: %w", err)
		}

		return false, nil
	}

	return true, nil
}

func WPConfigExists(cfg *config.Config) (bool, error) {
	err := WPCli(cfg, "config", "list")
	if err != nil {
		if _, ok := err.(*ssh.ExitError); !ok {
			return false, fmt.Errorf("Failed to check if WordPress config exists: %w", err)
		}

		return false, nil
	}

	return true, nil
}
