package local

import (
	"croox/wpclone/config"
)

func WPInstall(cfg *config.Config) error {
	downloadArgs := []string{
		"core", "download",
		"--locale=" + cfg.Flags.Locale,
	}

	if err := WPCli(cfg, downloadArgs...); err != nil {
		return err
	}

	createArgs := []string{
		"config", "create",
		"--dbhost=" + cfg.LocalDBHost(),
		"--dbname=" + cfg.LocalDBName(),
		"--dbuser=" + cfg.LocalDBUser(),
		"--dbpass=" + cfg.LocalDBPassword(),
	}

	if err := WPCli(cfg, createArgs...); err != nil {
		return err
	}

	installArgs := []string{
		"core", "install",
		"--url=" + cfg.LocalURL(),
		"--title=" + "My wpclone Site",
		"--admin_user=" + "admin",
		"--admin_password=" + cfg.Flags.AdminPassword,
		"--admin_email=" + cfg.Flags.AdminEmail,
	}

	if err := WPCli(cfg, installArgs...); err != nil {
		return err
	}

	return nil
}
