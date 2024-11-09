package config

import (
	"croox/wpclone/config/global"
	"croox/wpclone/pkg/defaults"
	"croox/wpclone/pkg/util"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

var pushExcludes = []string{
	"wp-config.php",
	".user.ini",
	".htaccess",
}

var cleanupFiles = []string{
	"wp-content/uploads/omgf", // remove omgf cache
}

var localRAWConfigSettings = map[string]string{
	"DISABLE_WP_CRON": "false",
}

var localConfigSettings = map[string]string{}

//go:embed *.tmpl
var tpls embed.FS

type Empty struct{}

func NewConfig() *Config {
	return &Config{}
}

func Load(flags Flags) (*Config, error) {
	cfg := Config{
		Flags: flags,
	}

	// load gloabl config
	global, err := global.Load()
	if err != nil {
		return nil, nil
	}
	cfg.Global = global

	if err := util.EnsureDir(cfg.CertDirPath()); err != nil {
		return nil, err
	}

	body, err := os.ReadFile(cfg.WPCloneYAMLPath())
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(body, &cfg)
	if err != nil {
		return nil, err
	}

	if err := checkMandatoryFields(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func WPCloneYAMLPath(flags Flags) string {
	if flags.ConfigFilePath != "" {
		return util.AbsPath(flags.ConfigFilePath)
	}

	cfgFileDir := WPCloneYAMLDir(flags)
	cfgFileName := WPCloneYAMLName(flags)
	cfgFilePath := filepath.Join(cfgFileDir, cfgFileName)

	return cfgFilePath
}

func WPCloneYAMLName(flags Flags) string {
	projectName := flags.Project

	if projectName != "" && !strings.HasSuffix(projectName, ".") {
		projectName = fmt.Sprintf("%s.", projectName)
	}

	return fmt.Sprintf("%swpclone.yml", projectName)
}

func WPCloneYAMLDir(flags Flags) string {
	if flags.ConfigFilePath != "" {
		configFilePath := util.AbsPath(flags.ConfigFilePath)
		return filepath.Dir(configFilePath)
	}

	if flags.Project != "" {
		return ProjectDir()
	}

	return util.Getwd()
}

func ConfigDir() string {
	return filepath.Join(util.UserHome(), ".wpclone")
}

func ProjectDir() string {
	d := filepath.Join(util.UserHome(), "wpclone")
	return defaults.String(os.Getenv("WPCLONE_PROJECT_DIR"), d)
}

func CertDirPath() string {
	return filepath.Join(global.ConfigDir(), "certs")
}

func SaveInit(filePath string, cfg Config) error {
	t, err := template.ParseFS(tpls, "*")
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return t.ExecuteTemplate(file, "init.yml.tmpl", cfg)
}

func ContainerName(s string) string {
	return fmt.Sprintf("wpclone_%s", s)
}

func ContainerNameWithID(s string) string {
	id := uuid.New()
	return ContainerName(fmt.Sprintf("%s_%s", s, id))
}

func checkMandatoryFields(cfg *Config) error {
	if cfg.Local.Path == "" {
		return fmt.Errorf("local.path is mandatory")
	}

	if cfg.Local.URL == "" {
		return fmt.Errorf("local.url is mandatory")
	}

	return nil
}
