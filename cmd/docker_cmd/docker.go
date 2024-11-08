package docker_cmd

import (
	"croox/wpclone/pkg/dock"
)

func removeWPClone() ([]dock.WPCloneContainerInfo, error) {
	containers, err := dock.RemoveAllContainers()
	if err != nil {
		return nil, err
	}

	if err := dock.RemoveAllNetworks(); err != nil {
		return nil, err
	}

	return containers, nil
}

func listContainers(all bool) ([]dock.WPCloneContainerInfo, error) {
	if all {
		return dock.ListContainers()
	}

	return dock.ListContainers("wp")
}
