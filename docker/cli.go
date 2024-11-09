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

	opts := RunOpts{
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

	if err := RunRsync(opts, args...); err != nil {
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
	opts := RunOpts{
		Binds: []string{
			cfg.LocalPath() + ":/var/www/html",
		},
	}

	return RunShell(opts, name, arg...)
}

type RunOpts struct {
	Binds []string
	Env   []string
}

func RunRsync(runOpts RunOpts, arg ...string) error {
	opts := dock.RunOptions{
		Name:  config.ContainerNameWithID("rsync"),
		Image: imageRsync,
		Binds: runOpts.Binds,
		Cmd:   arg,
		Labels: map[string]string{
			"wpclone_type":      "rsync",
			"wpclone_ephimeral": "true",
		},
		Env: runOpts.Env,
	}

	if err := dock.Run(opts); err != nil {
		return err
	}

	return nil
}

func RunShell(runOpts RunOpts, name string, arg ...string) error {
	return dock.Run(dock.RunOptions{
		Cmd:        append([]string{name}, arg...),
		Name:       config.ContainerNameWithID("cli"),
		Image:      imageCLI,
		Network:    networkProxy,
		Binds:      runOpts.Binds,
		WorkingDir: "/var/www/html",
		Labels: map[string]string{
			"wpclone_type":      "cli",
			"wpclone_ephimeral": "true",
		},
	})
}
