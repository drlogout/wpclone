package local

import (
	"croox/wpclone/config"
	"croox/wpclone/pkg/util"
	"fmt"
	"path/filepath"
)

func PullFiles(cfg *config.Config) error {
	remotePath := util.NoTrailingSlash(cfg.RemotePath())
	localPath := util.NoTrailingSlash(cfg.LocalPath())

	src := fmt.Sprintf("%s@%s:%s/", cfg.RemoteSSHUser(), cfg.RemoteSSHHost(), remotePath)
	dest := fmt.Sprintf("%s/", localPath)

	return Rsync(cfg, src, dest)
}

func PushFiles(cfg *config.Config) error {
	localPath := util.NoTrailingSlash(cfg.LocalPath())
	remotePath := util.NoTrailingSlash(cfg.RemotePath())

	src := fmt.Sprintf("%s/", localPath)
	dest := fmt.Sprintf("%s@%s:%s/", cfg.RemoteSSHUser(), cfg.RemoteSSHHost(), remotePath)

	// Do not override cfg.PushExcludes() files on remote
	return Rsync(cfg, src, dest, cfg.PushExcludes()...)
}

func PushWPConfig(cfg *config.Config) error {
	localPath := util.NoTrailingSlash(cfg.LocalPath())
	remotePath := util.NoTrailingSlash(cfg.RemotePath())

	src := filepath.Join(localPath, "wp-config.php")
	dest := fmt.Sprintf("%s@%s:%s/", cfg.RemoteSSHUser(), cfg.RemoteSSHHost(), remotePath)

	return Rsync(cfg, src, dest)
}
