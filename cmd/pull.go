package cmd

import (
	"croox/wpclone/cmd/common"
	"croox/wpclone/config"
	"croox/wpclone/docker"
	"croox/wpclone/local"
	"croox/wpclone/pkg/dock"
	"croox/wpclone/pkg/message"

	"github.com/urfave/cli/v2"
)

var Pull = &cli.Command{
	Name:  "pull",
	Usage: "Clone remote site to local",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "skip-rsync",
			Value: false,
		},
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Usage:   "Force pull (override local folder)",
			Value:   false,
		},
	},
	Before: common.BeforeCmds([]cli.BeforeFunc{
		common.BeforeLoadConfig,
		common.BeforeCheckGeneric,
		common.BeforeCheckLocalFolder,
		common.BeforeCheckRemoteLogin,
		common.BeforeCheckRemoteFolderExists,
		common.BeforeCheckRemoteWP,
	}),
	Action: func(ctx *cli.Context) error {
		cfg := common.ConfigFromCTX(ctx)

		if err := pull(cfg); err != nil {
			return message.ExitError(err, "pull failed")
		}

		message.Successf("Successfully pulled from %s to %s", cfg.RemoteURL(), cfg.LocalURL())
		return nil
	},
}

func pull(cfg *config.Config) error {
	if cfg.RunInDocker() {
		return pullToContainer(cfg)
	}

	return pullToLocal(cfg)
}

func pullToContainer(cfg *config.Config) error {
	spinner := message.NewInfoSpinner()

	spinner.Start("Setting up Docker environment")
	defer spinner.Stop("")

	opts := dock.WPOptions{
		Name:       cfg.DockerWPContainerName(),
		LocalPath:  cfg.LocalPath(),
		SSHKeyPath: cfg.SSHKeyPath(),
		URL:        cfg.LocalURL(),
		FQDN:       cfg.LocalFQDN(),
		CertDir:    cfg.CertDirPath(),
		SSLEnabled: cfg.DockerSSLEnabled(),
	}
	if err := dock.EnsureWP(opts); err != nil {
		return err
	}

	if err := dock.DBCreate(cfg.LocalDBName()); err != nil {
		return err
	}
	spinner.Stop("Set up Docker environment")

	if !cfg.Flags.SkipRsync {
		spinner.Start("Pulling files from remote")
		if err := docker.PullFiles(cfg); err != nil {
			return err
		}
		spinner.Stop("Pulled files from remote")
	}

	spinner.Start("Pulling database from remote")
	// no external utils used, hance no need to run in docker
	if err := local.PullDB(cfg); err != nil {
		return err
	}
	spinner.Stop("Pulled database from remote")

	spinner.Start("Configuring local environment")
	if err := docker.Configure(cfg); err != nil {
		return err
	}
	spinner.Stop("Configured local environment")

	spinner.Start("Importing database")
	if err := docker.ImportDB(cfg); err != nil {
		return err
	}
	spinner.Stop("Imported database")

	spinner.Start("Replacing URLs")
	if err := docker.SearchReplace(cfg); err != nil {
		return err
	}
	spinner.Stop("Replaced URLs")

	spinner.Start("Cleaning up")
	if err := docker.CleanUp(cfg); err != nil {
		return err
	}
	spinner.Stop("Cleaned up")

	spinner.Start("Updating proxy")
	if err := dock.UpdateProxy(); err != nil {
		return err
	}
	spinner.Stop("Updated proxy")

	return nil
}

func pullToLocal(cfg *config.Config) error {
	spinner := message.NewInfoSpinner()
	defer spinner.Stop("")

	if cfg.DockerDBOnly() {
		spinner.Start("Setting up Docker environment for DB")
		_, err := dock.EnsureDB()
		if err != nil {
			return err
		}

		if err := dock.DBCreate(cfg.LocalDBName()); err != nil {
			return err
		}
		spinner.Stop("Set up Docker environment for DB")
	}

	if !cfg.Flags.SkipRsync {
		spinner.Start("Pulling files from remote")
		if err := local.PullFiles(cfg); err != nil {
			return err
		}
		spinner.Stop("Pulled files from remote")
	}

	spinner.Start("Pulling database from remote")
	if err := local.PullDB(cfg); err != nil {
		return err
	}
	spinner.Stop("Pulled database from remote")

	spinner.Start("Configuring local environment")
	if err := local.Configure(cfg); err != nil {
		return err
	}
	spinner.Stop("Configured local environment")

	spinner.Start("Importing database")
	if err := local.ImportDB(cfg); err != nil {
		return err
	}
	spinner.Stop("Importing database")

	spinner.Start("Replacing URLs")
	if err := local.SearchReplace(cfg); err != nil {
		return err
	}
	spinner.Stop("Replacing URLs")

	spinner.Start("Cleaning up")
	if err := local.CleanUp(cfg); err != nil {
		return err
	}
	spinner.Stop("Cleaning up")

	return nil
}
