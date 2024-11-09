package dock

import (
	"fmt"
	"io"
	"reflect"
	"runtime"
	"strings"
	"time"

	docker "github.com/fsouza/go-dockerclient"
)

const (
	exitError = 1
)

var Verbose bool

type ContainerOptions struct {
	AutoRemove     bool
	Binds          []string
	Cmd            []string
	Env            []string
	HealthCheck    *docker.HealthConfig
	Image          string
	Labels         map[string]string
	Name           string
	Ports          map[docker.Port][]docker.PortBinding
	PrimaryNetwork *docker.Network
	Stdout         io.Writer
	Stderr         io.Writer
	Verbose        bool
	WorkingDir     string
	RestartPolicy  string
}

func CreateContainer(client *docker.Client, opts ContainerOptions) (*docker.Container, error) {
	if err := ensureImage(client, opts.Image); err != nil {
		return nil, err
	}

	healtCheck := &docker.HealthConfig{}
	if opts.HealthCheck != nil {
		healtCheck = opts.HealthCheck
	}

	labels := map[string]string{}
	if opts.Labels != nil {
		for key, value := range opts.Labels {
			labels[key] = value
		}
	}

	primaryNetworkConfig := &docker.NetworkingConfig{}
	if opts.PrimaryNetwork != nil {
		primaryNetworkConfig = &docker.NetworkingConfig{
			EndpointsConfig: map[string]*docker.EndpointConfig{
				opts.PrimaryNetwork.Name: {
					NetworkID: opts.PrimaryNetwork.ID,
				},
			},
		}
	}

	exposePorts := map[docker.Port]struct{}{}
	portBindings := map[docker.Port][]docker.PortBinding{}
	if opts.Ports != nil {
		for port := range opts.Ports {
			exposePorts[port] = struct{}{}
			portBindings[port] = opts.Ports[port]
		}
	}

	containerOptions := docker.CreateContainerOptions{
		Name: opts.Name,
		Config: &docker.Config{
			Image:        opts.Image,
			Env:          opts.Env,
			Healthcheck:  healtCheck,
			ExposedPorts: exposePorts,
			Labels:       labels,
			Cmd:          opts.Cmd,
		},
		HostConfig: &docker.HostConfig{
			Binds:        opts.Binds,
			PortBindings: portBindings,
			AutoRemove:   opts.AutoRemove,
		},
		NetworkingConfig: primaryNetworkConfig,
	}
	if opts.WorkingDir != "" {
		containerOptions.Config.WorkingDir = opts.WorkingDir
	}

	if opts.RestartPolicy != "" {
		containerOptions.HostConfig.RestartPolicy = docker.RestartPolicy{
			Name: opts.RestartPolicy,
		}
	}

	container, err := client.CreateContainer(containerOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}

	return container, nil
}

func StartContainer(client *docker.Client, id string) error {
	if err := client.StartContainer(id, nil); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	return nil
}

func CreateAndStartContainer(client *docker.Client, opts ContainerOptions) (*docker.APIContainers, error) {
	dockerContainer, err := CreateContainer(client, opts)
	if err != nil {
		return nil, err
	}

	if err := StartContainer(client, dockerContainer.ID); err != nil {
		return nil, err
	}

	return GetContainer(client, opts.Name)
}

func EnsureContainer(client *docker.Client, opts ContainerOptions) (*docker.APIContainers, error) {
	container, err := GetContainer(client, opts.Name)
	if err != nil {
		return nil, err
	}

	if container != nil {
		if !isRunning(container) {
			if err := StartContainer(client, container.ID); err != nil {
				return nil, err
			}
		}

		if !hasExpectedLabels(container.Labels, opts.Labels) {
			if err := StopAndRemoveContainer(client, container.ID); err != nil {
				return nil, err
			}

			return CreateAndStartContainer(client, opts)
		}
		return container, nil
	}

	return CreateAndStartContainer(client, opts)
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func WaitForContainerHealthy(client *docker.Client, containerID string, timeout time.Duration) error {
	startTime := time.Now()

	for {
		if time.Since(startTime) > timeout {
			return fmt.Errorf("timeout reached while waiting for container to become healthy")
		}

		container, err := client.InspectContainerWithOptions(docker.InspectContainerOptions{
			ID: containerID,
		})
		if err != nil {
			return fmt.Errorf("%s | failed to inspect container: %w", GetFunctionName(WaitForContainerHealthy), err)
		}

		switch container.State.Health.Status {
		case "healthy":
			return nil
		}

		time.Sleep(5 * time.Second)
	}
}

func StopAndRemoveContainer(client *docker.Client, containerID string) error {
	container, err := client.InspectContainerWithOptions(docker.InspectContainerOptions{
		ID: containerID,
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "No such container") {
			return nil
		}

		return fmt.Errorf("%s | failed to inspect container: %w", GetFunctionName(StopAndRemoveContainer), err)
	}

	if container.State.Running {
		if err := client.StopContainer(containerID, 10); err != nil {
			return fmt.Errorf("failed to stop container: %w", err)
		}
	}

	if container.State.Dead {
		return nil
	}

	err = client.RemoveContainer(docker.RemoveContainerOptions{
		ID:    containerID,
		Force: true, // Force removal
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "No such container") {
			return nil
		}

		return fmt.Errorf("%s | failed to remove container: %w", GetFunctionName(StopAndRemoveContainer), err)
	}

	return nil
}

func RestartContainer(client *docker.Client, containerName string) error {
	container, err := GetContainer(client, containerName)
	if err != nil {
		return err
	}

	if container == nil {
		return fmt.Errorf("container %s not found", containerName)
	}

	if err := client.RestartContainer(container.ID, 10); err != nil {
		return err
	}

	return nil
}

func GetContainer(client *docker.Client, name string) (*docker.APIContainers, error) {
	containers, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		return nil, err
	}
	for _, container := range containers {
		for _, containerName := range container.Names {
			if containerName == "/"+name {
				return &container, nil
			}
		}
	}
	return nil, nil
}

func isRunning(container *docker.APIContainers) bool {
	return container.State == "running"
}

func hasExpectedLabels(containerLabels, expectedLabels map[string]string) bool {
	for key, value := range expectedLabels {
		if containerLabels[key] != value {
			return false
		}
	}

	return true
}
