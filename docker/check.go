package docker

import (
	"croox/wpclone/config"
	"croox/wpclone/pkg/dock"
)

func UsedPorts() ([]int, error) {
	ports := make([]int, 0)
	proxyPorts, err := dock.GetContainerPorts(config.ContainerName("proxy"), "tcp")
	if err != nil {
		return nil, err
	}
	ports = append(ports, proxyPorts...)

	dnsmasqPorts, err := dock.GetContainerPorts(config.ContainerName("dnsmasq"), "udp")
	if err != nil {
		return nil, err
	}
	ports = append(ports, dnsmasqPorts...)

	dbPorts, err := dock.GetContainerPorts(config.ContainerName("db"), "tcp")
	if err != nil {
		return nil, err
	}
	ports = append(ports, dbPorts...)

	return ports, nil
}
