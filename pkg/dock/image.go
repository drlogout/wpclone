package dock

import (
	"fmt"

	docker "github.com/fsouza/go-dockerclient"
	log "github.com/sirupsen/logrus"
)

func ensureImage(client *docker.Client, name string) error {
	_, err := client.InspectImage(name)
	if err == nil {
		return nil
	}

	if err != docker.ErrNoSuchImage {
		return fmt.Errorf("failed to ensure image %s: %w", name, err)
	}

	log.Debugf("Pulling image %s", name)

	err = client.PullImage(docker.PullImageOptions{
		Repository: name,
	}, docker.AuthConfiguration{})
	if err != nil {
		return fmt.Errorf("failed to ensure image %s: %w", name, err)
	}

	return nil
}
