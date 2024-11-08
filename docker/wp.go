package docker

import (
	"croox/wpclone/config"
	"croox/wpclone/pkg/dock"
)

func InstallWP(cfg *config.Config) error {
	arg := []string{
		"--db-host", cfg.LocalDBHost(),
		"--db-name", cfg.LocalDBName(),
		"--db-user", cfg.LocalDBUser(),
		"--db-password", cfg.LocalDBPassword(),
		"--locale", cfg.Flags.Locale,
		"--url", cfg.LocalURL(),
		"--admin-email", cfg.Flags.AdminEmail,
		"--admin-password", cfg.Flags.AdminPassword,
	}

	if err := Shell(cfg, "wpclone_install_wordpress", arg...); err != nil {
		return err
	}

	return nil
}

func RemoveWP(cfg *config.Config) (*dock.Site, error) {
	wpContainer, err := dock.EnsureRemovedContainer(cfg.DockerWPContainerName())
	if err != nil {
		return nil, err
	}

	if wpContainer.State == "missing" {
		return nil, nil
	}

	if err := dock.DBRemove(cfg.LocalDBName()); err != nil {
		return nil, err
	}

	if err := dock.UpdateProxy(); err != nil {
		return nil, err
	}

	return &dock.Site{
		WP: wpContainer,
	}, nil
}

func StopWP(cfg *config.Config) error {
	if err := dock.StopWP(cfg.DockerWPContainerName()); err != nil {
		return err
	}

	return nil
}
