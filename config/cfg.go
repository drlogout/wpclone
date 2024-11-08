package config

import (
	"croox/wpclone/config/global"
	"croox/wpclone/pkg/util"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Remote Remote        `yaml:"remote"`
	Local  Local         `yaml:"local"`
	Flags  Flags         `yaml:"-"`
	Global global.Config `yaml:"-"`
}

func (cfg Config) WPCloneYAMLPath() string {
	return WPCloneYAMLPath(cfg.Flags)
}
func (cfg Config) WPCloneYAMLDir() string {
	return WPCloneYAMLDir(cfg.Flags)
}

func (cfg Config) CertDirPath() string {
	return CertDirPath()
}

func (cfg Config) SSHKeyPath() string {
	if cfg.Remote.SSH.Key == "" {
		userHome := util.UserHome()

		ed25519Path := fmt.Sprintf("%s/.ssh/id_ed25519", userHome)
		if _, err := os.Stat(ed25519Path); err == nil {
			return ed25519Path
		}

		return fmt.Sprintf("%s/.ssh/id_rsa", userHome)
	}

	if strings.HasPrefix(cfg.Remote.SSH.Key, "~/") {
		userHome := util.UserHome()

		return filepath.Join(userHome, strings.TrimLeft(cfg.Remote.SSH.Key, "~/"))
	}

	if filepath.IsAbs(cfg.Remote.SSH.Key) {
		return cfg.Remote.SSH.Key
	}

	configDir := filepath.Dir((cfg.WPCloneYAMLPath()))

	return util.JoinPath(configDir, cfg.Remote.SSH.Key)
}

func (cfg Config) RunInDocker() bool {
	if cfg.RunAllInDocker() {
		return true
	}

	if cfg.Local.Docker != nil && (cfg.Local.Docker.All || !cfg.Local.Docker.DBOnly) {
		return true
	}

	return false
}

func (cfg Config) DockerDBOnly() bool {
	return cfg.Local.Docker != nil && cfg.Local.Docker.DBOnly
}

func (cfg Config) RunAllInDocker() bool {
	return cfg.Global.DockerOnly
}

func (cfg Config) Mode() string {
	if cfg.RunInDocker() {
		return "docker"
	}

	return "system"
}

func (cfg Config) PushExcludes() []string {
	return pushExcludes
}

type Remote struct {
	Path           string `yaml:"path"`
	URL            string `yaml:"url"`
	SSH            SSH    `yaml:"ssh"`
	WWWUser        string `yaml:"www_user"`
	WPCli          string `yaml:"wp_cli"`
	DB             DB     `yaml:"db"`
	WPConfigExists bool   `yaml:"-"`
}

type Local struct {
	Path   string  `yaml:"path"`
	URL    string  `yaml:"url"`
	DB     DB      `yaml:"db"`
	WPCli  string  `yaml:"wp_cli"`
	Docker *Docker `yaml:"docker"`
}

type Docker struct {
	All    bool `yaml:"all"`
	SSL    bool `yaml:"ssl"`
	DBOnly bool `yaml:"db_only"`
}

type SSH struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Key      string `yaml:"key"`
	Password string `yaml:"password"`
}

type DB struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type Flags struct {
	Verbose           bool   `yaml:"verbose"`
	Debug             bool   `yaml:"debug"`
	ConfigFilePath    string `yaml:"config"`
	Dir               string `yaml:"dir"`
	Project           string `yaml:"project"`
	SkipRsync         bool   `yaml:"skip_rsync"`
	ListAllContainers bool   `yaml:"list_all_containers"`
	Locale            string `yaml:"locale"`
	AdminEmail        string `yaml:"admin_email"`
	AdminPassword     string `yaml:"admin_password"`
	PushWPConfig      bool   `yaml:"push_wp_config"`
}
