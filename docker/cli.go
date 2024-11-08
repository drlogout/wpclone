package docker

import (
	"croox/wpclone/config"
	"croox/wpclone/pkg/dock"
	"croox/wpclone/pkg/util"
	"fmt"
)

func Shell(cfg *config.Config, name string, arg ...string) error {
	return runShell(cfg, name, arg...)
}

func Rsync(cfg *config.Config, src, dest string, excludes ...string) error {
	args := []string{
		"-avz",
		"--delete",
	}

	for _, e := range excludes {
		args = append(args, "--exclude", e)
	}

	sshCommand := util.SSHCommand(cfg.RemoteSSHPort(), "/root/.ssh/id_rsa")
	sshCommand = fmt.Sprintf("'%s'", sshCommand)

	args = append(args, "-e", sshCommand)
	args = append(args, src, dest)

	opts := dock.RunOpts{
		Binds: []string{
			cfg.LocalPath() + ":/var/www/html",
			cfg.SSHKeyPath() + ":/root/.ssh/id_rsa",
		},
	}

	if cfg.RunInDocker() {
		opts.Env = []string{
			"RSYNC_UID=33",
			"RSYNC_GID=33",
		}
	}

	if err := dock.RunRsync(opts, args...); err != nil {
		return err
	}

	return nil
}

func DBConfig(cfg *config.Config) error {
	return runShell(cfg, "wpclone_get_dbconf", "wpclone_remote_wp-config.php")
}

func WPCli(cfg *config.Config, arg ...string) error {
	return runShell(cfg, "wp", arg...)
}

func runShell(cfg *config.Config, name string, arg ...string) error {
	opts := dock.RunOpts{
		Binds: []string{
			cfg.LocalPath() + ":/var/www/html",
		},
	}

	return dock.RunShell(opts, name, arg...)
}
