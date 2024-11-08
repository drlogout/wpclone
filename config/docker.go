package config

import (
	"croox/wpclone/pkg/dock"
	"fmt"
)

func (cfg Config) DockerDBDumpPath() string {
	return fmt.Sprintf("/var/www/html/%s", localDumpName)
}

func (cfg Config) DockerWPContainerName() string {
	n := fmt.Sprintf("wp_%s", cfg.LocalSlug())
	return dock.ContainerName(n)
}

func (cfg Config) DockerDBHost() string {
	return "wpclone_db"
}
