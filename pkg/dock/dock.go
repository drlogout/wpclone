package dock

import (
	docker "github.com/fsouza/go-dockerclient"
)

var client *docker.Client

func GetClient() (*docker.Client, error) {
	var err error

	if client != nil {
		return client, nil
	}

	client, err = docker.NewClientFromEnv()

	return client, err
}
