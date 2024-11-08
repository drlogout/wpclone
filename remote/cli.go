package remote

import (
	"croox/wpclone/config"
	"croox/wpclone/pkg/sshexec"

	"golang.org/x/crypto/ssh"
)

func Shell(cfg *config.Config, cmd string, args ...string) error {
	opts := sshexec.RunOpts{
		Dir: cfg.RemotePath(),
	}
	return ShellWithOpts(cfg, opts, cmd, args...)
}

func ShellWithOpts(cfg *config.Config, opts sshexec.RunOpts, cmd string, args ...string) error {
	opts.SSHHost = cfg.RemoteSSHHost()
	opts.SSHPort = cfg.RemoteSSHPort()
	opts.SSHUser = cfg.RemoteSSHUser()
	opts.SSHPassword = cfg.RemoteSSHPassword()
	opts.SSHKeyPath = cfg.SSHKeyPath()

	return sshexec.RunWithOpts(opts, cmd, args...)
}

func Bash(cfg *config.Config, bashCmd string) error {
	return BashWithOpts(cfg, sshexec.RunOpts{}, bashCmd)
}

func BashWithOpts(cfg *config.Config, opts sshexec.RunOpts, bashCmd string) error {
	bashCmd = "\"" + bashCmd + "\""
	return ShellWithOpts(cfg, opts, "bash", "-c", bashCmd)
}

func WPCli(cfg *config.Config, args ...string) error {
	return Shell(cfg, cfg.RemoteWPCli(), args...)
}

func DBExists(cfg *config.Config) (bool, error) {
	err := WPCli(cfg, "db", "check")
	if err != nil {
		if _, ok := err.(*ssh.ExitError); !ok {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

func FileExists(cfg *config.Config, path string) (bool, error) {
	err := Shell(cfg, "test", "-e", path)
	if err != nil {
		if _, ok := err.(*ssh.ExitError); !ok {
			return false, err
		}

		return false, nil
	}

	return true, nil
}
