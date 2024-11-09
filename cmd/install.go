package cmd

import (
	"croox/wpclone/cmd/common"
	"croox/wpclone/config"
	"croox/wpclone/docker"
	"croox/wpclone/local"
	"croox/wpclone/pkg/message"
	"croox/wpclone/pkg/ternary"
	"croox/wpclone/pkg/util"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/urfave/cli/v2"
)

var Install = &cli.Command{
	Name:  "install",
	Usage: "Install local WordPress",
	Before: common.BeforeCmds([]cli.BeforeFunc{
		common.BeforeLoadConfig,
		common.BeforeCheckGeneric,
	}),
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "locale",
			Usage:   "WordPress locale (default: de_DE)",
			Value:   "de_DE",
			Aliases: []string{"l"},
		},
		&cli.StringFlag{
			Name:    "admin-email",
			Usage:   "WordPress admin email",
			Aliases: []string{"e"},
		},
		&cli.StringFlag{
			Name:    "admin-password",
			Usage:   "WordPress admin password",
			Aliases: []string{"p"},
		},
	},
	Action: func(ctx *cli.Context) error {
		cfg := common.ConfigFromCTX(ctx)

		if err := util.EnsureDir(cfg.LocalPath()); err != nil {
			return err
		}

		empty, err := util.FolderEmpty(cfg.LocalPath())
		if err != nil {
			return err
		}

		if !empty {
			label := fmt.Sprintf("❓ Local folder %s is not empty, do you want to reinstall?", cfg.LocalPath())
			ok := util.YesNoPrompt(label, false)
			if !ok {
				return fmt.Errorf("Aborted")
			}
		}

		if err := ensureWPBlankSlate(cfg); err != nil {
			return err
		}

		if err := os.MkdirAll(cfg.LocalPath(), 0755); err != nil {
			return err
		}

		if cfg.Flags.AdminEmail == "" {
			cfg.Flags.AdminEmail = "admin@example.com"
		}

		if cfg.Flags.AdminPassword == "" {
			cfg.Flags.AdminPassword = util.RandomID(16)
		}

		if err := installWP(cfg); err != nil {
			return message.ExitError(err, "install failed")
		}

		message.Success("Successfully installed WordPress")
		mode := ternary.String(cfg.RunInDocker(), "docker", "local")

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Local path", "Local URL", "Mode", "Admin user", "Admin email", "Admin password"})
		t.AppendRow(table.Row{cfg.LocalPath(), cfg.LocalURL(), mode, "admin", cfg.Flags.AdminEmail, cfg.Flags.AdminPassword})
		t.AppendSeparator()
		t.Render()

		return nil
	},
}

func installWP(cfg *config.Config) error {
	if cfg.RunInDocker() {
		return installWithContainerWPCli(cfg)
	}

	return installWithLocalWPCli(cfg)
}

func ensureWPBlankSlate(cfg *config.Config) error {
	if cfg.RunInDocker() || cfg.DockerDBOnly() {
		_, err := docker.EnsureDB()
		if err != nil {
			return err
		}

		exists, err := docker.DBExists(cfg.LocalDBName())
		if err != nil {
			return err
		}

		if exists {
			if err := docker.DBRemove(cfg.LocalDBName()); err != nil {
				return err
			}
		}

		if err := docker.DBCreate(cfg.LocalDBName()); err != nil {
			return err
		}

		if err := docker.DBWait(cfg.LocalDBName()); err != nil {
			return err
		}
	} else {
		exists, err := local.DBExists(cfg)
		if err != nil {
			return err
		}

		if exists {
			mysqlDropArgs := []string{
				"-h", cfg.LocalDBHost(),
				"-u", cfg.LocalDBUser(),
				"-p" + cfg.LocalDBPassword(),
				"-e", "DROP DATABASE " + cfg.LocalDBName(),
			}

			if err := local.Shell(cfg, "mysql", mysqlDropArgs...); err != nil {
				return err
			}
		}

		mysqlCreateArgs := []string{
			"-h", cfg.LocalDBHost(),
			"-u", cfg.LocalDBUser(),
			"-p" + cfg.LocalDBPassword(),
			"-e", "CREATE DATABASE " + cfg.LocalDBName(),
		}

		if err := local.Shell(cfg, "mysql", mysqlCreateArgs...); err != nil {
			return err
		}
	}

	if err := os.RemoveAll(cfg.LocalPath()); err != nil {
		return err
	}

	return nil
}

func installWithContainerWPCli(cfg *config.Config) error {
	opts := docker.WPOptions{
		Name:       cfg.DockerWPContainerName(),
		LocalPath:  cfg.LocalPath(),
		SSHKeyPath: cfg.SSHKeyPath(),
		URL:        cfg.LocalURL(),
		FQDN:       cfg.LocalFQDN(),
		CertDir:    cfg.CertDirPath(),
		SSLEnabled: cfg.DockerSSLEnabled(),
	}
	if err := docker.EnsureWP(opts); err != nil {
		return err
	}

	if err := docker.DBCreate(cfg.LocalDBName()); err != nil {
		return err
	}

	if err := docker.InstallWP(cfg); err != nil {
		return err
	}

	if err := docker.Configure(cfg); err != nil {
		return err
	}

	if err := docker.UpdateProxy(); err != nil {
		return err
	}

	return nil
}

func installWithLocalWPCli(cfg *config.Config) error {
	if cfg.DockerDBOnly() {
		_, err := docker.EnsureDB()
		if err != nil {
			return err
		}

		if err := docker.DBCreate(cfg.LocalDBName()); err != nil {
			return err
		}
	}

	if err := local.WPInstall(cfg); err != nil {
		return err
	}

	if err := local.Configure(cfg); err != nil {
		return err
	}

	return nil
}
