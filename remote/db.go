package remote

import (
	"croox/wpclone/config"
	"croox/wpclone/pkg/wp"

	_ "github.com/go-sql-driver/mysql"
)

func ImportDB(cfg *config.Config) error {
	defer Shell(cfg, "rm", "-f", cfg.RemoteDBDumpPath())

	exists, err := DBExists(cfg)
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

	if err := WPCli(cfg, "db", "import", cfg.RemoteDBDumpPath()); err != nil {
		return err
	}

	return nil
}

func SearchReplace(cfg *config.Config) error {
	variants, err := wp.URLVariants(cfg.LocalURL())
	if err != nil {
		return err
	}

	for _, url := range variants {
		if err := WPCli(cfg, "search-replace", url, cfg.RemoteURL()); err != nil {
			return err
		}
	}

	return nil
}
