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

func RemoveWP(cfg *config.Config) (*Site, error) {
	wpContainer, err := dock.EnsureRemovedContainer(cfg.DockerWPContainerName())
	if err != nil {
		return nil, err
	}

	if wpContainer.State == "missing" {
		return nil, nil
	}

	if err := DBRemove(cfg.LocalDBName()); err != nil {
		return nil, err
	}

	if err := UpdateProxy(); err != nil {
		return nil, err
	}

	return &Site{
		WP: getWPCloneContainerInfo(wpContainer),
	}, nil
}

func StopWP(cfg *config.Config) error {
	client, err := dock.GetClient()
	if err != nil {
		return err
	}

	if err := dock.StopAndRemoveContainer(client, cfg.DockerWPContainerName()); err != nil {
		return err
	}

	return nil
}

func IsWPRunning(name string) (bool, error) {
	client, err := dock.GetClient()
	if err != nil {
		return false, err
	}

	container, err := dock.GetContainer(client, name)
	if err != nil {
		return false, err
	}

	if container == nil {
		return false, nil
	}

	return true, nil
}
