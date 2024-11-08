package dock

func Ping() error {
	client, err := getClient()
	if err != nil {
		return err
	}

	_, err = client.Info()
	if err != nil {
		return err
	}

	return nil
}

func Ports() ([]int, error) {
	ports := make([]int, 0)
	proxyPorts, err := getPorts(ContainerName("proxy"), "tcp")
	if err != nil {
		return nil, err
	}
	ports = append(ports, proxyPorts...)

	dnsmasqPorts, err := getPorts(ContainerName("dnsmasq"), "udp")
	if err != nil {
		return nil, err
	}
	ports = append(ports, dnsmasqPorts...)

	dbPorts, err := getPorts(ContainerName("db"), "tcp")
	if err != nil {
		return nil, err
	}
	ports = append(ports, dbPorts...)

	return ports, nil
}

func getPorts(name, proto string) ([]int, error) {
	ports := make([]int, 0)

	client, err := getClient()
	if err != nil {
		return nil, err
	}

	container, err := getContainer(client, name)
	if err != nil {
		return nil, err
	}

	if container == nil {
		return ports, nil
	}

	for _, port := range container.Ports {
		if port.Type == proto && port.PublicPort != 0 {
			ports = append(ports, int(port.PublicPort))
		}
	}

	return ports, nil
}
