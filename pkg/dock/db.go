package dock

import (
	"bytes"
	"fmt"
	"io"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	log "github.com/sirupsen/logrus"
)

const (
	wpCloneContainerDBRootPassword = "wpclone"
)

func EnsureDB() (*docker.APIContainers, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	network, err := ensureNetwork(client, networkProxy)
	if err != nil {
		return nil, err
	}

	volume, err := ensureVolume(client, volumeDB)
	if err != nil {
		return nil, err
	}

	container, err := ensureContainer(client, ContainerOptions{
		Name:           ContainerName("db"),
		Image:          imageDB,
		PrimaryNetwork: network,
		Binds: []string{
			volume.Name + ":/var/lib/mysql",
		},
		Env: []string{
			fmt.Sprintf("MARIADB_ROOT_PASSWORD=%s", wpCloneContainerDBRootPassword),
		},
		HealthCheck: &docker.HealthConfig{
			Test:     []string{"CMD", "mysqladmin", "ping", "-h", "127.0.0.1", "-u", "root", fmt.Sprintf("-p%s", wpCloneContainerDBRootPassword)},
			Interval: 10 * time.Second,
			Timeout:  5 * time.Second,
			Retries:  3,
		},
		Labels: map[string]string{
			"wpclone_type": "db",
		},
		Ports: map[docker.Port][]docker.PortBinding{
			"3306/tcp": {
				{
					HostIP:   "127.0.0.1",
					HostPort: "3306",
				},
			},
		},
		RestartPolicy: "unless-stopped",
	})
	if err != nil {
		return nil, err
	}

	log.Debugf("Ensure %s", ContainerName("db"))

	if err := waitForContainerHealthy(client, container.ID, time.Second*60); err != nil {
		return nil, err
	}

	log.Debugf("Container %s is healthy", ContainerName("db"))

	return container, nil
}

func DBCreate(dbName string) error {
	opts := ExecOptions{
		ContainerName: ContainerName("db"),
		Cmd:           []string{"db-create", dbName},
	}

	_, err := Exec(opts)
	if err != nil {
		return err
	}

	log.Debugf("Database %s created", dbName)

	return nil
}

func DBExists(dbname string) (bool, error) {
	var buffer bytes.Buffer
	var writer io.Writer = &buffer

	_, err := Exec(ExecOptions{
		ContainerName: ContainerName("db"),
		Cmd:           []string{"db-exists", dbname},
		Stdout:        writer,
	})

	return buffer.String() != "", err
}

func DBRemove(dbName string) error {
	opts := ExecOptions{
		ContainerName: ContainerName("db"),
		Cmd:           []string{"db-remove", dbName},
	}

	_, err := Exec(opts)
	if err != nil {
		return err
	}

	log.Debugf("Database %s removed", dbName)

	return nil
}

func DBWait(dbName string) error {
	opts := ExecOptions{
		ContainerName: ContainerName("db"),
		Cmd:           []string{"db-wait", dbName},
	}

	_, err := Exec(opts)
	if err != nil {
		return fmt.Errorf("database %s is not ready", dbName)
	}

	log.Debugf("Database %s is ready", dbName)

	return nil
}
