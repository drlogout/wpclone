package local

import (
	"croox/wpclone/config"
	sshexec "croox/wpclone/pkg/sshexec"
	"croox/wpclone/pkg/util"
	"croox/wpclone/pkg/wp"
	"croox/wpclone/remote"

	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func PullDB(cfg *config.Config) error {
	outfile, err := os.Create(cfg.LocalDBDumpPath())
	if err != nil {
		return err
	}
	defer outfile.Close()

	opts := sshexec.RunOpts{
		Stdout: outfile,
		Dir:    cfg.RemotePath(),
	}

	if err := remote.ShellWithOpts(cfg, opts, cfg.RemoteWPCli(), "db", "export", "-"); err != nil {
		return err
	}

	if err := util.RemoveLineInFile(cfg.LocalDBDumpPath(), "^.*enable the sandbox mode.*$"); err != nil {
		return err
	}

	return err
}

func PushDB(cfg *config.Config) error {
	defer os.Remove(cfg.LocalDBDumpPath())

	if err := WPCli(cfg, "db", "export", cfg.LocalDBDumpPath()); err != nil {
		return err
	}

	if err := util.RemoveLineInFile(cfg.LocalDBDumpPath(), "^.*enable the sandbox mode.*$"); err != nil {
		return err
	}

	dest := fmt.Sprintf("%s@%s:%s", cfg.RemoteSSHUser(), cfg.RemoteSSHHost(), cfg.RemoteDBDumpPath())

	if err := Rsync(cfg, cfg.LocalDBDumpPath(), dest); err != nil {
		return err
	}

	return nil
}

func ImportDB(cfg *config.Config) error {
	defer os.Remove(cfg.LocalDBDumpPath())

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

	if err := WPCli(cfg, "db", "import", cfg.LocalDBDumpPath()); err != nil {
		return err
	}

	return nil
}

func SearchReplace(cfg *config.Config) error {
	variants, err := wp.URLVariants(cfg.RemoteURL())
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

func DBExists(cfg *config.Config) (bool, error) {
	db, err := sql.Open("mysql", cfg.LocalDBConnectionString())
	if err != nil {
		return false, fmt.Errorf("failed to open database connection: %w", err)
	}
	defer db.Close()

	// Query to check if the database exists
	query := `SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?`

	var dbName string
	err = db.QueryRow(query, cfg.LocalDBName()).Scan(&dbName)
	if err != nil {
		if err == sql.ErrNoRows {
			// Database does not exist
			return false, nil
		}

		return false, fmt.Errorf("failed to query database: %w", err)
	}

	// Database exists
	return true, nil
}
