package remote

import (
	"croox/wpclone/config"
	"fmt"
)

func Configure(cfg *config.Config) error {
	dbSettings := map[string]string{
		"DB_HOST":     fmt.Sprintf("%s:%d", cfg.RemoteDBHost(), cfg.RemoteDBPort()),
		"DB_NAME":     cfg.RemoteDBName(),
		"DB_USER":     cfg.RemoteDBUser(),
		"DB_PASSWORD": cfg.RemoteDBPassword(),
	}

	for key, value := range dbSettings {
		if err := WPCli(cfg, "config", "set", key, value); err != nil {
			return err
		}
	}

	// delete local wpclone settings
	for key := range cfg.LocalConfigSettings() {
		if err := WPCli(cfg, "config", "delete", key); err != nil {
			return err
		}
	}

	for key := range cfg.LocalRAWConfigSettings() {
		if err := WPCli(cfg, "config", "delete", key); err != nil {
			return err
		}
	}

	return nil
}
