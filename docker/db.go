package docker

import (
	"bytes"
	"croox/wpclone/config"
	"croox/wpclone/pkg/dock"
	"croox/wpclone/pkg/wp"
	"fmt"
	"io"
	"strings"

	log "github.com/sirupsen/logrus"
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
	exists, err := DBExists(cfg.LocalDBName())
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

func ListDBs() ([]string, error) {
	var buffer bytes.Buffer
	var writer io.Writer = &buffer
	dbs := []string{}

	_, err := EnsureDB()
	if err != nil {
		return dbs, err
	}

	opts := dock.ExecOptions{
		ContainerName: config.ContainerName("db"),
		Cmd:           []string{"dbs-list"},
		Stdout:        writer,
	}

	status, err := dock.Exec(opts)
	if err != nil {
		return dbs, err
	}

	if status != 0 {
		return dbs, fmt.Errorf("failed to list databases")
	}

	output := strings.Split(buffer.String(), "\n")
	for _, db := range output {
		if db != "" && noSystemDB(db) {
			dbs = append(dbs, db)
		}
	}

	return dbs, nil
}

func DBCreate(dbName string) error {
	opts := dock.ExecOptions{
		ContainerName: config.ContainerName("db"),
		Cmd:           []string{"db-create", dbName},
	}

	_, err := dock.Exec(opts)
	if err != nil {
		return err
	}

	log.Debugf("Database %s created", dbName)

	return nil
}

func DBExists(dbname string) (bool, error) {
	var buffer bytes.Buffer
	var writer io.Writer = &buffer

	_, err := dock.Exec(dock.ExecOptions{
		ContainerName: config.ContainerName("db"),
		Cmd:           []string{"db-exists", dbname},
		Stdout:        writer,
	})

	return buffer.String() != "", err
}

func DBRemove(dbName string) error {
	opts := dock.ExecOptions{
		ContainerName: config.ContainerName("db"),
		Cmd:           []string{"db-remove", dbName},
	}

	_, err := dock.Exec(opts)
	if err != nil {
		return err
	}

	log.Debugf("Database %s removed", dbName)

	return nil
}

func DBWait(dbName string) error {
	opts := dock.ExecOptions{
		ContainerName: config.ContainerName("db"),
		Cmd:           []string{"db-wait", dbName},
	}

	_, err := dock.Exec(opts)
	if err != nil {
		return fmt.Errorf("database %s is not ready", dbName)
	}

	log.Debugf("Database %s is ready", dbName)

	return nil
}

func noSystemDB(db string) bool {
	return db != "information_schema" && db != "performance_schema" && db != "mysql" && db != "sys"
}
