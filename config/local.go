package config

import (
	"croox/wpclone/pkg/util"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

const (
	localDumpName = "wpclone_db.sql"
)

func (cfg Config) LocalSlug() string {
	return util.Slug(cfg.LocalFQDN())
}

func (cfg Config) LocalPath() string {
	path := cfg.Local.Path

	if strings.HasPrefix(path, "./") {
		return filepath.Join(cfg.WPCloneYAMLDir(), path)
	}

	return util.AbsPath(path)
}

func (cfg Config) LocalURL() string {
	if cfg.RunInDocker() {
		if cfg.DockerSSLEnabled() {
			return fmt.Sprintf("https://%s", cfg.LocalFQDN())
		}

		return fmt.Sprintf("http://%s", cfg.LocalFQDN())
	}

	return util.NoTrailingSlash(cfg.Local.URL)
}

func (cfg Config) LocalFQDN() string {
	u, err := url.Parse(cfg.Local.URL)
	if err != nil {
		panic(err)
	}

	return u.Hostname()
}

func (cfg Config) LocalDBHost() string {
	if cfg.RunInDocker() {
		return cfg.DockerDBHost()
	}

	if cfg.DockerDBOnly() {
		return "127.0.0.1"
	}

	if cfg.Local.DB.Host != "" {
		return cfg.Local.DB.Host
	}

	return "127.0.0.1"
}

func (cfg Config) LocalDBName() string {
	if cfg.RunInDocker() || cfg.DockerDBOnly() {
		return cfg.LocalSlug()
	}

	return cfg.Local.DB.Name
}

func (cfg Config) LocalDBUser() string {
	if cfg.RunInDocker() || cfg.DockerDBOnly() {
		return cfg.LocalSlug()
	}

	return cfg.Local.DB.User
}

func (cfg Config) LocalDBPassword() string {
	if cfg.RunInDocker() || cfg.DockerDBOnly() {
		return cfg.LocalSlug()
	}

	return cfg.Local.DB.Password
}

func (cfg Config) LocalDBPort() int {
	if cfg.RunInDocker() || cfg.DockerDBOnly() {
		return 3306
	}

	if cfg.Local.DB.Port != 0 {
		return cfg.Local.DB.Port
	}

	return 3306
}

func (cfg Config) LocalDBConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/", cfg.LocalDBUser(), cfg.LocalDBPassword(), cfg.LocalDBHost(), cfg.LocalDBPort())
}

func (cfg Config) LocalDBDumpPath() string {
	return filepath.Join(cfg.LocalPath(), localDumpName)
}

func (cfg Config) LocalWPCli() string {
	if cfg.Local.WPCli != "" {
		return cfg.Local.WPCli
	}

	return "wp"
}

func (cfg Config) LocalCleanupPaths() []string {
	cleanupPaths := []string{}

	for _, f := range cleanupFiles {
		cleanupPaths = append(cleanupPaths, filepath.Join(cfg.LocalPath(), f))
	}

	return cleanupPaths
}

func (cfg Config) LocalConfigSettings() map[string]string {
	return localConfigSettings
}

func (cfg Config) LocalRAWConfigSettings() map[string]string {
	return localRAWConfigSettings
}

func (cfg Config) LocalConfigSettingKeys() []string {
	keys := []string{}

	for key := range localConfigSettings {
		keys = append(keys, key)
	}

	return keys
}

func (cfg Config) LocalRAWConfigSettingKeys() []string {
	keys := []string{}

	for key := range localRAWConfigSettings {
		keys = append(keys, key)
	}

	return keys
}

func (cfg Config) DockerSSLEnabled() bool {
	return cfg.Local.Docker != nil && cfg.Local.Docker.SSL
}
