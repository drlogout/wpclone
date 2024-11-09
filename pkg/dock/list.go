package dock

import (
	docker "github.com/fsouza/go-dockerclient"
)

func ListContainers(filters map[string][]string) ([]docker.APIContainers, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}

	containers, err := client.ListContainers(docker.ListContainersOptions{
		All:     true,
		Filters: filters,
	})
	if err != nil {
		return nil, err
	}

	wpContainers := []docker.APIContainers{}

	if len(containers) == 0 {
		return wpContainers, nil
	}

	for _, container := range containers {
		wpContainers = append(wpContainers, container)
	}

	return wpContainers, nil
}

func ListVolumes(filters map[string][]string) ([]docker.Volume, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}

	options := docker.ListVolumesOptions{
		Filters: filters,
	}

	volumeList, err := client.ListVolumes(options)
	if err != nil {
		return nil, err
	}

	return volumeList, nil
}

func ListNetworks(filters map[string]map[string]bool) ([]docker.Network, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}

	options := docker.NetworkFilterOpts(filters)

	networks, err := client.FilteredListNetworks(options)
	if err != nil {
		return nil, err
	}

	return networks, nil
}
