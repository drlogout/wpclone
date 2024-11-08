package common

import (
	"croox/wpclone/config"
	"croox/wpclone/pkg/message"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/urfave/cli/v2"
)

var expectedCommands = []string{
	"mysql",
	"mysqldump",
	"rsync",
	"ssh",
	"wp",
}

var expectedDockerCommands = []string{
	"docker",
}

var BeforeCheckSystem = func(ctx *cli.Context) error {
	switch runtime.GOOS {
	case "linux", "darwin", "freebsd", "openbsd", "netbsd", "solaris", "aix":
		return nil
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

var BeforeCheckGeneric = func(ctx *cli.Context) error {
	cfg := ConfigFromCTX(ctx)

	if cfg.RunInDocker() {
		beforeCmds := BeforeCmds([]cli.BeforeFunc{
			BeforeCheckDocker,
			BeforeCheckDockerWebPorts,
			BeforeCheckDockerDBPorts,
		})

		if err := beforeCmds(ctx); err != nil {
			return err
		}
	}

	if cfg.DockerDBOnly() {
		beforeCmds := BeforeCmds([]cli.BeforeFunc{
			BeforeCheckDocker,
			BeforeCheckDockerDBPorts,
		})

		if err := beforeCmds(ctx); err != nil {
			return err
		}
	}

	if err := checkCommands(cfg); err != nil {
		return err
	}

	message.Infof("Using config file: %s, (mode: %s)", cfg.WPCloneYAMLPath(), cfg.Mode())

	return nil
}

func checkCommands(cfg *config.Config) error {
	if cfg.RunAllInDocker() {
		return runCommandsCheck(expectedDockerCommands)
	}

	checkCmds := expectedCommands

	if cfg.RunInDocker() {
		checkCmds = append(checkCmds, expectedDockerCommands...)
	}

	return runCommandsCheck(checkCmds)
}

func runCommandsCheck(expectedCommands []string) error {
	missingCommands := []string{}

	for _, command := range expectedCommands {
		_, err := exec.LookPath(command)
		if err != nil {
			missingCommands = append(missingCommands, command)
		}
	}

	if len(missingCommands) > 0 {
		return fmt.Errorf("the following commands are not installed: %s", strings.Join(missingCommands, ", "))
	}

	return nil
}
