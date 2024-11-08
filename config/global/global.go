package global

import (
	"croox/wpclone/pkg/util"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Network    Network `yaml:"network"`
	DockerOnly bool    `yaml:"docker_only"`
}

type Network struct {
	DB_Port    int `yaml:"db_port"`
	HTTP_Port  int `yaml:"http_port"`
	HTTPS_Port int `yaml:"https_port"`
}

func newConfig() Config {
	return Config{
		Network: Network{
			DB_Port:    3306,
			HTTP_Port:  80,
			HTTPS_Port: 443,
		},
		DockerOnly: false,
	}
}

func Load() (Config, error) {
	if err := ensureConfig(); err != nil {
		return Config{}, err
	}

	var cfg Config

	if err := util.LoadYAML(ConfigFilePath(), &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func ensureConfig() error {
	if !util.FileExists(ConfigDir()) {
		if err := os.MkdirAll(ConfigDir(), 0755); err != nil {
			return err
		}
	}

	if util.FileExists(ConfigFilePath()) {
		return nil
	}

	if err := writeConfig(newConfig()); err != nil {
		return err
	}

	log.Debug("Initializing global config")

	return nil
}

func ConfigDir() string {
	return filepath.Join(util.UserHome(), ".wpclone")
}

func ConfigFilePath() string {
	return filepath.Join(ConfigDir(), "global.yml")
}

func writeConfig(cfg Config) error {
	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}

	if err := os.WriteFile(ConfigFilePath(), data, 0644); err != nil {
		return err
	}

	return nil
}
