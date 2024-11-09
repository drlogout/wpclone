package docker_cmd

import (
	"croox/wpclone/docker"
)

func removeWPClone() ([]docker.WPCloneContainer, error) {
	containers, err := docker.RemoveAllContainers()
	if err != nil {
		return nil, err
	}

	if err := docker.RemoveAllNetworks(); err != nil {
		return nil, err
	}

	return containers, nil
}
