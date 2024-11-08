package config

import (
	"croox/wpclone/pkg/util"
	"fmt"
	"net/url"
	"strings"
)

const (
	remoteSQLDumpName = "wpclone_remote_db.sql"
)

func (cfg Config) RemoteSlug() string {
	return util.Slug(cfg.RemoteFQDN())
}

func (cfg Config) RemotePath() string {
	return cfg.Remote.Path
}

func (cfg Config) RemoteURL() string {
	return util.NoTrailingSlash(cfg.Remote.URL)
}

func (cfg Config) RemoteFQDN() string {
	remoteURL := cfg.Remote.URL
	if !strings.HasPrefix(remoteURL, "http") {
		remoteURL = fmt.Sprintf("http://%s", remoteURL)
	}

	u, err := url.Parse(remoteURL)
	if err != nil {
		panic(err)
	}

	return u.Hostname()
}

func (cfg Config) RemoteDBDumpPath() string {
	remotePath := util.NoTrailingSlash(cfg.RemotePath())
	return fmt.Sprintf("%s/%s", remotePath, remoteSQLDumpName)
}

func (cfg Config) RemoteWWWUser() string {
	if cfg.Remote.WWWUser != "" {
		return cfg.Remote.WWWUser
	}

	return cfg.RemoteSSHUser()
}

func (cfg Config) RemoteWPCli() string {
	if cfg.Remote.WPCli != "" {
		return cfg.Remote.WPCli
	}

	return "wp"
}

func (cfg Config) RemoteSSHUser() string {
	return cfg.Remote.SSH.User
}

func (cfg Config) RemoteSSHHost() string {
	return cfg.Remote.SSH.Host
}

func (cfg Config) RemoteSSHPassword() string {
	return cfg.Remote.SSH.Password
}

func (cfg Config) RemoteSSHPort() int {
	if cfg.Remote.SSH.Port != 0 {
		return cfg.Remote.SSH.Port
	}

	return 22
}

func (cfg Config) RemoteCleanupPaths() []string {
	cleanupPaths := []string{}

	remotePath := util.NoTrailingSlash(cfg.RemotePath())

	for _, f := range cleanupFiles {
		cleanupPaths = append(cleanupPaths, fmt.Sprintf("%s/%s", remotePath, f))
	}

	return cleanupPaths
}

func (cfg Config) RemoteDBHost() string {
	if cfg.Remote.DB.Host != "" {
		return cfg.Remote.DB.Host
	}

	return "127.0.0.1"
}

func (cfg Config) RemoteDBName() string {
	return cfg.Remote.DB.Name
}

func (cfg Config) RemoteDBUser() string {
	return cfg.Remote.DB.User
}

func (cfg Config) RemoteDBPassword() string {
	return cfg.Remote.DB.Password
}

func (cfg Config) RemoteDBPort() int {
	if cfg.Remote.DB.Port != 0 {
		return cfg.Remote.DB.Port
	}

	return 3306
}

func (cfg Config) RemoteDBCredentailsSet() bool {
	return cfg.Remote.DB.Name != "" && cfg.Remote.DB.User != "" && cfg.Remote.DB.Password != ""
}
