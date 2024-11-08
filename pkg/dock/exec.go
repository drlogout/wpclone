package dock

import (
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"

	docker "github.com/fsouza/go-dockerclient"
)

type ExecOptions struct {
	ContainerName string
	Cmd           []string
	WorkingDir    string
	User          string
	Verbose       bool
	Stdout        io.Writer
	Interactive   bool
}

func Exec(opts ExecOptions) (int, error) {
	log.Debugf("[docker.exec@%s] %s", opts.ContainerName, opts.Cmd)

	client, err := getClient()
	if err != nil {
		return exitError, err
	}

	if Verbose {
		opts.Verbose = true
	}

	exec, err := createExec(client, opts)
	if err != nil {
		return exitError, err
	}

	startExecConfig := docker.StartExecOptions{}

	if opts.Verbose {
		startExecConfig.ErrorStream = os.Stderr
		startExecConfig.OutputStream = os.Stdout
	}

	if opts.Stdout != nil {
		startExecConfig.OutputStream = opts.Stdout
	}

	if opts.Interactive {
		startExecConfig.Tty = true
		startExecConfig.InputStream = os.Stdin
		startExecConfig.OutputStream = os.Stdout
		startExecConfig.ErrorStream = os.Stderr
		startExecConfig.RawTerminal = true
	}

	if err := client.StartExec(exec.ID, startExecConfig); err != nil {
		return exitError, err
	}

	inspect, err := client.InspectExec(exec.ID)

	return inspect.ExitCode, err
}

func createExec(client *docker.Client, opts ExecOptions) (*docker.Exec, error) {
	container, err := getContainer(client, opts.ContainerName)
	if err != nil {
		return nil, err
	}

	if container == nil {
		return nil, fmt.Errorf("container %s not found", opts.ContainerName)
	}

	execConfig := docker.CreateExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Container:    container.ID,
		Cmd:          opts.Cmd,
	}

	if opts.WorkingDir != "" {
		execConfig.WorkingDir = opts.WorkingDir
	}

	if opts.User != "" {
		execConfig.User = opts.User
	}

	if opts.Interactive {
		execConfig.Tty = true
		execConfig.AttachStdin = true
		execConfig.Env = []string{"TERM=xterm"}
	}

	return client.CreateExec(execConfig)
}
