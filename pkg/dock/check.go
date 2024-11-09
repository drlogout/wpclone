package dock

func Ping() error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	_, err = client.Info()
	if err != nil {
		return err
	}

	return nil
}

func GetContainerPorts(name, proto string) ([]int, error) {
	ports := make([]int, 0)

	client, err := GetClient()
	if err != nil {
		return nil, err
	}

	container, err := GetContainer(client, name)
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
