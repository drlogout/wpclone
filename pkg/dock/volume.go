package dock

import (
	docker "github.com/fsouza/go-dockerclient"
	log "github.com/sirupsen/logrus"
)

func ensureVolume(client *docker.Client, name string) (*docker.Volume, error) {
	existingVolume, err := getVolume(client, name)
	if err != nil && err != docker.ErrNoSuchVolume {
		return nil, err
	}

	if existingVolume != nil {
		return existingVolume, nil
	}

	volumeOptions := docker.CreateVolumeOptions{
		Name:   name,
		Driver: "local",
		Labels: map[string]string{
			"wpclone_type": containerNameSuffix(name),
		},
	}

	volume, err := client.CreateVolume(volumeOptions)
	if err != nil {
		return nil, err
	}

	log.Debugf("Successfully created volume: %s", volume.Name)
	return volume, nil
}

func getVolume(client *docker.Client, name string) (*docker.Volume, error) {
	volume, err := client.InspectVolume(name)
	if err != nil {
		return nil, err
	}

	return volume, nil
}
