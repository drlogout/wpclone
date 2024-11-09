package dock

import (
	"fmt"
	"io"
	"os"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
	log "github.com/sirupsen/logrus"
)

type RunOptions struct {
	Name       string
	Cmd        []string
	Image      string
	Network    string
	Binds      []string
	Env        []string
	WorkingDir string
	Labels     map[string]string
	Writer     io.Writer
}

func Run(opts RunOptions) error {
	log.Debugf("[docker.run]: %s", strings.Join(opts.Cmd, " "))

	client, err := GetClient()
	if err != nil {
		return err
	}

	containerOpts := ContainerOptions{
		Name:       opts.Name,
		Image:      opts.Image,
		Binds:      opts.Binds,
		Cmd:        opts.Cmd,
		WorkingDir: opts.WorkingDir,
		Labels:     opts.Labels,
		AutoRemove: true,
	}

	if opts.Writer != nil {
		containerOpts.Stdout = opts.Writer
	}

	if opts.Network != "" {
		n, err := EnsureNetwork(client, opts.Network)
		if err != nil {
			return err
		}

		containerOpts.PrimaryNetwork = n
	}

	container, _, err := runContainer(client, containerOpts)
	if err != nil {
		return err
	}

	if err := StopAndRemoveContainer(client, container.ID); err != nil {
		return err
	}

	return err
}

func runContainer(client *docker.Client, opts ContainerOptions) (*docker.APIContainers, int, error) {
	var status = 1

	if Verbose {
		opts.Verbose = true
	}

	if err := ensureContainerRemoved(client, opts.Name); err != nil {
		return nil, status, err
	}

	dockerContainer, err := CreateContainer(client, opts)
	if err != nil {
		return nil, status, err
	}

	attachOptions := docker.AttachToContainerOptions{
		Container: dockerContainer.ID,
		Stdout:    true,
		Stderr:    true,
		Stream:    true,
	}

	if opts.Verbose {
		attachOptions.ErrorStream = os.Stderr
		attachOptions.OutputStream = os.Stdout
	}

	if opts.Stdout != nil {
		attachOptions.OutputStream = opts.Stdout
	}

	// Start the container and attach to its output
	go func() {
		err = client.AttachToContainer(attachOptions)
		if err != nil {
			log.Fatalf("Failed to attach to container: %s", err)
		}
	}()

	if err := StartContainer(client, dockerContainer.ID); err != nil {
		return nil, status, err
	}

	// Wait for the container to finish
	status, err = client.WaitContainer(dockerContainer.ID)
	if err != nil {
		log.Fatalf("Failed to wait for container: %s", err)
	}

	if status != 0 {
		return nil, status, fmt.Errorf("container exited with status %d", status)
	}

	container, err := GetContainer(client, opts.Name)
	if err != nil {
		return nil, status, err
	}

	return container, status, nil
}
