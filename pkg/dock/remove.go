package dock

import (
	docker "github.com/fsouza/go-dockerclient"
)

func RemoveContainers(containers []docker.APIContainers) ([]docker.APIContainers, error) {
	if len(containers) == 0 {
		return []docker.APIContainers{}, nil
	}

	for _, container := range containers {
		if err := RemoveContainer(container.ID); err != nil {
			return nil, err
		}
	}

	return containers, nil
}

func RemoveContainer(id string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	return StopAndRemoveContainer(client, id)
}

func RemoveVolumes(filters map[string][]string) error {
	volumes, err := ListVolumes(filters)
	if err != nil {
		return err
	}

	for _, volume := range volumes {
		if err := RemoveVolume(volume.Name); err != nil {
			return err
		}
	}

	return nil
}

func RemoveVolume(name string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	return client.RemoveVolumeWithOptions(docker.RemoveVolumeOptions{
		Name:  name,
		Force: true, // Force removal
	})
}

func RemoveNetworks(filters map[string]map[string]bool) error {
	networks, err := ListNetworks(filters)
	if err != nil {
		return err
	}

	for _, network := range networks {
		if err := RemoveNetwork(network.ID); err != nil {
			return err
		}
	}

	return nil
}

func RemoveNetwork(id string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	return client.RemoveNetwork(id)
}

func EnsureRemovedContainer(name string) (docker.APIContainers, error) {
	client, err := GetClient()
	if err != nil {
		return docker.APIContainers{}, err
	}

	container, err := GetContainer(client, name)
	if err != nil {
		return docker.APIContainers{}, err
	}

	if container == nil {
		return docker.APIContainers{
			State: wpContainerStateMissing,
		}, nil
	}

	if err := client.RemoveContainer(docker.RemoveContainerOptions{
		ID:    container.ID,
		Force: true,
	}); err != nil {
		return docker.APIContainers{}, err
	}

	return *container, nil
}

func ensureContainerRemoved(client *docker.Client, name string) error {
	container, err := GetContainer(client, name)
	if err != nil {
		return err
	}

	if container != nil {
		if err := StopAndRemoveContainer(client, container.ID); err != nil {
			return err
		}
	}

	return nil
}
