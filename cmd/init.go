package cmd

import (
	"croox/wpclone/cmd/common"
	"croox/wpclone/config"
	"croox/wpclone/pkg/exec"
	"croox/wpclone/pkg/message"
	"croox/wpclone/pkg/util"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

var Init = &cli.Command{
	Name:  "init",
	Usage: "Initialize wpclone.yml",
	Action: func(ctx *cli.Context) error {
		flags := common.ConfigFlags(ctx)

		// if project is set create dir and .gitignore
		if flags.Project != "" {
			projectDir := config.WPCloneYAMLDir(flags)
			if !util.FileExists(projectDir) {
				if err := os.MkdirAll(projectDir, 0755); err != nil {
					return err
				}
			}

			gitignorePath := filepath.Join(projectDir, ".gitignore")
			if !util.FileExists(gitignorePath) {
				t, err := template.ParseFS(tpls, "*")
				if err != nil {
					return err
				}
				file, err := os.Create(gitignorePath)
				if err != nil {
					return err
				}
				defer file.Close()

				if err := t.ExecuteTemplate(file, "gitignore.tmpl", nil); err != nil {
					return err
				}
			}

			gitPath := filepath.Join(projectDir, ".git")
			if !util.FileExists(gitPath) {
				opts := exec.RunOpts{
					Dir: projectDir,
				}
				if err := exec.RunWithOpts(opts, "git", "init"); err != nil {
					return err
				}
			}
		}

		configFilePath := config.WPCloneYAMLPath(flags)
		baseDir := filepath.Dir(configFilePath)

		if !util.FileExists(baseDir) {
			return message.ExitError(fmt.Errorf("config directory %s does not exist", baseDir), "init failed")
		}

		if util.FileExists(configFilePath) {
			ok := util.YesNoPrompt("❓ "+configFilePath+" already exists, do you want to overwrite it?", false)
			if !ok {
				return message.Exit("Aborted")
			}
		}

		if err := initWPClone(configFilePath, flags); err != nil {
			return message.ExitError(err, "init failed")
		}

		message.Successf("Successfully initialized %s", configFilePath)
		return nil
	},
}

func initWPClone(configFilePath string, flags config.Flags) error {
	cfg := setInitialDefaults(configFilePath, flags)

	if err := config.SaveInit(configFilePath, cfg); err != nil {
		return err
	}

	return nil
}

func setInitialDefaults(configFilePath string, flags config.Flags) config.Config {
	cfg := config.NewConfig()

	baseDir := filepath.Dir(configFilePath)
	if flags.Project != "" {
		cfg.Local.Path = fmt.Sprintf("./%s", flags.Project)
	} else {
		cfg.Local.Path = filepath.Join(baseDir, "local")
	}

	cfg.Local.URL = "http://example.test"

	return *cfg
}
