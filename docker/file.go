package docker

import (
	"croox/wpclone/config"
	"croox/wpclone/pkg/util"
	"fmt"
)

func PullFiles(cfg *config.Config) error {
	remotePath := util.NoTrailingSlash(cfg.RemotePath())
	src := fmt.Sprintf("%s@%s:%s/", cfg.RemoteSSHUser(), cfg.RemoteSSHHost(), remotePath)
	dest := "/var/www/html/"

	return Rsync(cfg, src, dest)
}

func PushFiles(cfg *config.Config) error {
	remotePath := util.NoTrailingSlash(cfg.RemotePath())
	src := "/var/www/html/"
	dest := fmt.Sprintf("%s@%s:%s/", cfg.RemoteSSHUser(), cfg.RemoteSSHHost(), remotePath)

	// Do not override cfg.PushExcludes() files on remote
	return Rsync(cfg, src, dest, cfg.PushExcludes()...)
}

func PushWPConfig(cfg *config.Config) error {
	remotePath := util.NoTrailingSlash(cfg.RemotePath())
	src := "/var/www/html/wp-config.php"
	dest := fmt.Sprintf("%s@%s:%s/", cfg.RemoteSSHUser(), cfg.RemoteSSHHost(), remotePath)

	return Rsync(cfg, src, dest)
}
