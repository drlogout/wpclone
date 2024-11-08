package local

import (
	"croox/wpclone/config"
	"croox/wpclone/pkg/exec"
	"croox/wpclone/pkg/util"
)

func Shell(cfg *config.Config, cmd string, args ...string) error {
	opts := exec.RunOpts{
		Dir: cfg.LocalPath(),
	}

	return exec.RunWithOpts(opts, cmd, args...)
}

func ShellWithOpts(cfg *config.Config, opts exec.RunOpts, cmd string, args ...string) error {
	return exec.RunWithOpts(opts, cmd, args...)
}

func Rsync(cfg *config.Config, src, dest string, excludes ...string) error {
	return RsyncWithOpts(exec.RunOpts{}, cfg, src, dest, excludes...)
}

func RsyncWithOpts(opts exec.RunOpts, cfg *config.Config, src, dest string, excludes ...string) error {
	args := []string{
		"-avz",
		"--delete",
	}

	for _, e := range excludes {
		args = append(args, "--exclude", e)
	}

	args = append(args, "-e", util.SSHCommand(cfg.RemoteSSHPort(), cfg.SSHKeyPath()))
	args = append(args, src, dest)

	return exec.RunWithOpts(opts, "rsync", args...)
}

func WPCli(cfg *config.Config, args ...string) error {
	opts := exec.RunOpts{
		Dir: cfg.LocalPath(),
	}

	return exec.RunWithOpts(opts, cfg.LocalWPCli(), args...)
}
