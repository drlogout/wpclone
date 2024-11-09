package docker

import (
	"croox/wpclone/pkg/dock"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
)

type WPCloneContainer struct {
	ContainerName string
	FQDN          string
	ID            string
	SSLEnabled    bool
	State         string
	Type          string
	URL           string
}

type Site struct {
	WP WPCloneContainer
}

func RemoveAllContainers() ([]WPCloneContainer, error) {
	containers, err := ListAllContainers()
	if err != nil {
		return nil, err
	}

	apiContainers, err := dock.RemoveContainers(getAPIContainers(containers))
	if err != nil {
		return nil, err
	}

	return getWPCloneContainers(apiContainers), nil
}

func RemoveAllContainersExceptDB() ([]WPCloneContainer, error) {
	containers := []docker.APIContainers{}

	wpContainers, err := dock.ListContainers(getLabelFilter("wp"))
	if err != nil {
		return nil, err
	}
	containers = append(containers, wpContainers...)

	proxyContainer, err := dock.ListContainers(getLabelFilter("proxy"))
	if err != nil {
		return nil, err
	}
	containers = append(containers, proxyContainer...)

	dnsmasqContainer, err := dock.ListContainers(getLabelFilter("dnsmasq"))
	if err != nil {
		return nil, err
	}
	containers = append(containers, dnsmasqContainer...)

	apiContainers, err := dock.RemoveContainers(containers)
	if err != nil {
		return nil, err
	}

	return getWPCloneContainers(apiContainers), nil
}

func getWPCloneContainerInfo(container docker.APIContainers) WPCloneContainer {
	return WPCloneContainer{
		ContainerName: strings.TrimLeft(container.Names[0], "/"),
		FQDN:          container.Labels["wpclone_fqdn"],
		ID:            container.ID,
		SSLEnabled:    container.Labels["wpclone_ssl"] == "true",
		State:         container.State,
		URL:           container.Labels["wpclone_url"],
		Type:          container.Labels["wpclone_type"],
	}
}

func getWPCloneContainers(containers []docker.APIContainers) []WPCloneContainer {
	wpContainers := []WPCloneContainer{}
	for _, container := range containers {
		wpContainers = append(wpContainers, getWPCloneContainerInfo(container))
	}
	return wpContainers
}

func getAPIContainer(container WPCloneContainer) docker.APIContainers {
	return docker.APIContainers{
		ID:    container.ID,
		Names: []string{container.ContainerName},
	}
}

func getAPIContainers(containers []WPCloneContainer) []docker.APIContainers {
	apiContainers := []docker.APIContainers{}
	for _, container := range containers {
		apiContainers = append(apiContainers, getAPIContainer(container))
	}
	return apiContainers
}
