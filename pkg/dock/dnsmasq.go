package dock

import (
	docker "github.com/fsouza/go-dockerclient"
	log "github.com/sirupsen/logrus"
)

func ensureDNSMasq() (*docker.APIContainers, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	network, err := ensureNetwork(client, networkProxy)
	if err != nil {
		return nil, err
	}

	container, err := ensureContainer(client, ContainerOptions{
		Name:           ContainerName("dnsmasq"),
		Image:          imageDNSMasq,
		PrimaryNetwork: network,
		Ports: map[docker.Port][]docker.PortBinding{
			"53/udp": {
				{
					HostPort: "53",
				},
			},
		},
		Labels: map[string]string{
			"wpclone_type": "dnsmasq",
		},
		RestartPolicy: "unless-stopped",
	})
	if err != nil {
		return nil, err
	}

	log.Debugf("Container %s running", ContainerName("dnsmasq"))

	return container, nil
}
