package dock

import (
	docker "github.com/fsouza/go-dockerclient"
)

func RemoveAllContainers() ([]WPCloneContainerInfo, error) {
	containers, err := ListContainers()
	if err != nil {
		return nil, err
	}

	if len(containers) == 0 {
		return []WPCloneContainerInfo{}, nil
	}

	for _, container := range containers {
		if err := RemoveContainer(container.ID); err != nil {
			return nil, err
		}
	}

	return containers, nil
}

func RemoveAllContainersExceptDB() ([]WPCloneContainerInfo, error) {
	containers := []WPCloneContainerInfo{}

	wpContainers, err := ListContainers("wp")
	if err != nil {
		return nil, err
	}
	containers = append(containers, wpContainers...)

	proxyContainer, err := ListContainers("proxy")
	if err != nil {
		return nil, err
	}
	containers = append(containers, proxyContainer...)

	dnsmasqContainer, err := ListContainers("dnsmasq")
	if err != nil {
		return nil, err
	}
	containers = append(containers, dnsmasqContainer...)

	if len(containers) == 0 {
		return []WPCloneContainerInfo{}, nil
	}

	for _, container := range containers {
		if err := RemoveContainer(container.ID); err != nil {
			return nil, err
		}
	}

	return containers, nil
}

func RemoveContainer(id string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	return stopAndRemoveContainer(client, id)
}

func RemoveAllVolumes() error {
	volumes, err := ListVolumes()
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
	client, err := getClient()
	if err != nil {
		return err
	}

	return client.RemoveVolumeWithOptions(docker.RemoveVolumeOptions{
		Name:  name,
		Force: true, // Force removal
	})
}

func RemoveAllNetworks() error {
	networks, err := ListNetworks()
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
	client, err := getClient()
	if err != nil {
		return err
	}

	return client.RemoveNetwork(id)
}

func EnsureRemovedContainer(name string) (WPCloneContainerInfo, error) {
	client, err := getClient()
	if err != nil {
		return WPCloneContainerInfo{}, err
	}

	container, err := getContainer(client, name)
	if err != nil {
		return WPCloneContainerInfo{}, err
	}

	if container == nil {
		return WPCloneContainerInfo{
			State: wpContainerStateMissing,
		}, nil
	}

	if err := client.RemoveContainer(docker.RemoveContainerOptions{
		ID:    container.ID,
		Force: true,
	}); err != nil {
		return WPCloneContainerInfo{}, err
	}

	info := getWPCloneContainerInfo(*container)
	info.State = wpContainerStateDeleted

	return info, nil
}

func ensureContainerRemoved(client *docker.Client, name string) error {
	container, err := getContainer(client, name)
	if err != nil {
		return err
	}

	if container != nil {
		if err := stopAndRemoveContainer(client, container.ID); err != nil {
			return err
		}
	}

	return nil
}
