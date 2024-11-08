package common

import (
	"croox/wpclone/config"
	"croox/wpclone/pkg/message"
	"croox/wpclone/pkg/util"
	"croox/wpclone/remote"
	"fmt"
	"path"

	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh"
)

var BeforeCheckRemoteFolder = func(ctx *cli.Context) error {
	if ctx.Bool("force") {
		return nil
	}

	cfg := ConfigFromCTX(ctx)

	empty, err := isRemoteFolderEmpty(cfg, cfg.RemotePath())
	if err != nil {
		return err
	}

	if !empty {
		label := fmt.Sprintf("❓ Remote folder %s@%s:%s is not empty, do you want to proceed?", cfg.RemoteSSHUser(), cfg.RemoteSSHHost(), cfg.RemotePath())
		ok := util.YesNoPrompt(label, false)
		if !ok {
			return fmt.Errorf("Aborted")
		}
	}

	return nil
}

var BeforeCheckRemoteLogin = func(ctx *cli.Context) error {
	cfg := ConfigFromCTX(ctx)

	if err := remote.Shell(cfg, "exit"); err != nil {
		return message.Exitf("Failed to login to remote %s@%s and cd to %s", cfg.RemoteSSHUser(), cfg.RemoteSSHHost(), cfg.RemotePath())
	}

	return nil
}

var BeforeCheckRemoteWP = func(ctx *cli.Context) error {
	cfg := ConfigFromCTX(ctx)

	installed, err := remote.WPInstalled(cfg)
	if err != nil {
		return err
	}

	if !installed {
		return message.Exitf("WordPress is not installed on remote %s@%s:%s", cfg.RemoteSSHUser(), cfg.RemoteSSHHost(), cfg.RemotePath())
	}

	return nil
}

var BeforeCheckRemoteFolderExists = func(ctx *cli.Context) error {
	cfg := ConfigFromCTX(ctx)

	parentExists, err := remoteFolderExists(cfg, cfg.RemotePath())
	if err != nil {
		return err
	}

	if !parentExists {
		return message.Exitf("Remote directory %s does not exist on remote", cfg.RemotePath())
	}

	return nil
}

var BeforeCheckRemoteFolderParentExists = func(ctx *cli.Context) error {
	cfg := ConfigFromCTX(ctx)

	parentDir := path.Dir(cfg.RemotePath())

	parentExists, err := remoteFolderExists(cfg, parentDir)
	if err != nil {
		return err
	}

	if !parentExists {
		return message.Exitf("Parent directory %s does not exist on remote", parentDir)
	}

	return nil
}

func isRemoteFolderEmpty(cfg *config.Config, remotePath string) (bool, error) {
	bashCmd := fmt.Sprintf("test -z \"$(ls -A %s)\"", remotePath)

	err := remote.Bash(cfg, bashCmd)
	if err != nil {
		if _, ok := err.(*ssh.ExitError); !ok {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

func remoteFolderExists(cfg *config.Config, remotePath string) (bool, error) {
	bashCmd := fmt.Sprintf("test -d $(dirname %s)", remotePath)

	err := remote.Bash(cfg, bashCmd)
	if err != nil {
		if _, ok := err.(*ssh.ExitError); !ok {
			return false, err
		}

		return false, nil
	}

	return true, nil
}
