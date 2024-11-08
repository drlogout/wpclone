package docker

import (
	"croox/wpclone/config"
	"croox/wpclone/pkg/dock"
	"croox/wpclone/pkg/wp"
	"fmt"
)

func PushDB(cfg *config.Config) error {
	if err := WPCli(cfg, "db", "export", cfg.DockerDBDumpPath()); err != nil {
		return err
	}

	if err := Shell(cfg, "sed", "-i", `/\/\*!999999\\- enable the sandbox mode \*\//d`, cfg.DockerDBDumpPath()); err != nil {
		return err
	}

	dest := fmt.Sprintf("%s@%s:%s", cfg.RemoteSSHUser(), cfg.RemoteSSHHost(), cfg.RemoteDBDumpPath())

	if err := Rsync(cfg, cfg.DockerDBDumpPath(), dest); err != nil {
		return err
	}

	if err := Shell(cfg, "rm", "-f", cfg.DockerDBDumpPath()); err != nil {
		return err
	}

	return nil
}

func SearchReplace(cfg *config.Config) error {
	variants, err := wp.URLVariants(cfg.RemoteURL())
	if err != nil {
		return err
	}

	variants, err = wp.AppendSSLURLVariants(variants, cfg.LocalURL(), cfg.DockerSSLEnabled())
	if err != nil {
		return err
	}

	for _, url := range variants {
		if err := WPCli(cfg, "search-replace", url, cfg.LocalURL()); err != nil {
			return err
		}
	}

	return nil
}

func ImportDB(cfg *config.Config) error {
	exists, err := dock.DBExists(cfg.LocalDBName())
	if err != nil {
		return err
	}

	if exists {
		if err := WPCli(cfg, "db", "drop", "--yes"); err != nil {
			return err
		}
	}

	if err := WPCli(cfg, "db", "create"); err != nil {
		return err
	}

	if err := WPCli(cfg, "db", "import", cfg.DockerDBDumpPath()); err != nil {
		return err
	}

	if err := Shell(cfg, "rm", "-f", cfg.DockerDBDumpPath()); err != nil {
		return err
	}

	return nil
}
