package docker

import (
	"croox/wpclone/config"
	"croox/wpclone/local"
	"fmt"
)

func Configure(cfg *config.Config) error {
	dbSettings := map[string]string{
		"DB_HOST":     fmt.Sprintf("%s:%d", cfg.LocalDBHost(), cfg.LocalDBPort()),
		"DB_NAME":     cfg.LocalDBName(),
		"DB_USER":     cfg.LocalDBUser(),
		"DB_PASSWORD": cfg.LocalDBPassword(),
	}

	for key, value := range dbSettings {
		if err := WPCli(cfg, "config", "set", key, value); err != nil {
			return err
		}
	}

	for key, value := range cfg.LocalConfigSettings() {
		if err := WPCli(cfg, "config", "set", key, value); err != nil {
			return err
		}
	}

	// raw means: "Place the value into the wp-config.php file as is, instead of as a quoted string."
	for key, value := range cfg.LocalRAWConfigSettings() {
		if err := WPCli(cfg, "config", "set", key, value, "--raw"); err != nil {
			return err
		}
	}
	if local.HasUserINI(cfg.LocalPath()) {
		if err := local.UpdateUserINI(cfg.LocalPath(), cfg.LocalPath(), cfg.LocalPath()); err != nil {
			return err
		}
	}

	if cfg.RunInDocker() {
		if cfg.DockerSSLEnabled() {
			if err := Shell(cfg, "ensure-line", "SetEnv HTTPS on", ".htaccess"); err != nil {
				return err
			}
		} else {
			if err := Shell(cfg, "remove-line", "SetEnv HTTPS on", ".htaccess"); err != nil {
				return err
			}
		}
	}

	return nil
}
