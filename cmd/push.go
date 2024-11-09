package cmd

import (
	"croox/wpclone/cmd/common"
	"croox/wpclone/config"
	"croox/wpclone/docker"
	"croox/wpclone/local"
	"croox/wpclone/pkg/message"
	"croox/wpclone/remote"

	"github.com/urfave/cli/v2"
)

var Push = &cli.Command{
	Name:  "push",
	Usage: "Clone local site to remote",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "skip-rsync",
			Value: false,
		},
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Usage:   "Force push (override remote folder)",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "push-wp-config",
			Aliases: []string{"wp-config"},
			Usage:   "Force push local wp-config.php (DB settings required in wpclone.yml)",
			Value:   false,
		},
	},
	Before: common.BeforeCmds([]cli.BeforeFunc{
		common.BeforeLoadConfig,
		common.BeforeCheckGeneric,
		common.BeforeCheckRemoteLogin,
		common.BeforeCheckRemoteFolder,
		common.BeforeCheckRemoteFolderParentExists,
		common.BeforeCheckPush,
	}),
	Action: func(ctx *cli.Context) error {
		cfg := common.ConfigFromCTX(ctx)

		err := push(cfg)
		if err != nil {
			return message.ExitError(err, "push failed")
		}

		message.Successf("Successfully pushed from %s to %s", cfg.LocalURL(), cfg.RemoteURL())
		return nil
	},
}

func push(cfg *config.Config) error {
	if cfg.RunInDocker() {
		return pushFromContainer(cfg)
	}

	return pushFromLocal(cfg)
}

func pushFromContainer(cfg *config.Config) error {
	spinner := message.NewInfoSpinner()
	defer spinner.Stop("")

	spinner.Start("Setting up Docker environment")

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

	spinner.Stop("Set up Docker environment")

	if !cfg.Flags.SkipRsync {
		spinner.Start("Pushing files to remote")
		if err := docker.PushFiles(cfg); err != nil {
			return err
		}

		if err := ensureRemoteWPConfig(cfg); err != nil {
			return err
		}

		spinner.Stop("Pushed files to remote")
	}

	// TODO push in docker if DockerDBOnly
	spinner.Start("Pushing database to remote")
	if err := docker.PushDB(cfg); err != nil {
		return err
	}
	spinner.Stop("Pushed database to remote")

	spinner.Start("Importing database")
	if err := remote.ImportDB(cfg); err != nil {
		return err
	}
	spinner.Stop("Imported database")

	spinner.Start("Replacing URLs")
	if err := remote.SearchReplace(cfg); err != nil {
		return err
	}
	spinner.Stop("Replaced URLs")

	spinner.Start("Cleaning up")
	if err := remote.CleanUp(cfg); err != nil {
		return err
	}
	spinner.Stop("Cleaned up")

	return nil
}

func pushFromLocal(cfg *config.Config) error {
	spinner := message.NewInfoSpinner()
	defer spinner.Stop("")

	if cfg.DockerDBOnly() {
		spinner.Start("Setting up Docker environment for DB")
		_, err := docker.EnsureDB()
		if err != nil {
			return err
		}

		if err := docker.DBCreate(cfg.LocalDBName()); err != nil {
			return err
		}
		spinner.Stop("Set up Docker environment for DB")
	}

	if !cfg.Flags.SkipRsync {
		spinner.Start("Pushing files")
		if err := local.PushFiles(cfg); err != nil {
			return err
		}

		if err := ensureRemoteWPConfig(cfg); err != nil {
			return err
		}

		spinner.Stop("Pushed files")
	}

	spinner.Start("Pushing database")
	if err := local.PushDB(cfg); err != nil {
		return err
	}
	spinner.Stop("Pushed database")

	spinner.Start("Importing database")
	if err := remote.ImportDB(cfg); err != nil {
		return err
	}
	spinner.Stop("Imported database")

	spinner.Start("Replacing URLs")
	if err := remote.SearchReplace(cfg); err != nil {
		return err
	}
	spinner.Stop("Replaced URLs")

	spinner.Start("Cleaning up")
	if err := remote.CleanUp(cfg); err != nil {
		return err
	}
	spinner.Stop("Cleaned up")

	return nil
}

// TODO if PushWPConfig read remote config and use those values
func ensureRemoteWPConfig(cfg *config.Config) error {
	if cfg.Remote.WPConfigExists && !cfg.Flags.PushWPConfig {
		return nil
	}

	if cfg.RunInDocker() {
		if err := docker.PushWPConfig(cfg); err != nil {
			return err
		}
	} else {
		if err := local.PushWPConfig(cfg); err != nil {
			return err
		}
	}

	if err := remote.Configure(cfg); err != nil {
		return err
	}

	return nil
}
