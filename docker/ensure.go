package docker

import (
	"croox/wpclone/config"
	"croox/wpclone/pkg/dock"
	"fmt"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	log "github.com/sirupsen/logrus"
)

const (
	wpCloneContainerDBRootPassword = "wpclone"
)

func EnsureDB() (*docker.APIContainers, error) {
	client, err := dock.GetClient()
	if err != nil {
		return nil, err
	}

	network, err := dock.EnsureNetwork(client, networkProxy)
	if err != nil {
		return nil, err
	}

	volume, err := dock.EnsureVolume(client, volumeDB)
	if err != nil {
		return nil, err
	}

	container, err := dock.EnsureContainer(client, dock.ContainerOptions{
		Name:           config.ContainerName("db"),
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

	log.Debugf("Ensure %s", config.ContainerName("db"))

	if err := dock.WaitForContainerHealthy(client, container.ID, time.Second*60); err != nil {
		return nil, err
	}

	log.Debugf("Container %s is healthy", config.ContainerName("db"))

	return container, nil
}

type WPOptions struct {
	Name       string
	URL        string
	FQDN       string
	LocalPath  string
	SSHKeyPath string
	CertDir    string
	SSLEnabled bool
}

func EnsureWP(opts WPOptions) error {
	client, err := dock.GetClient()
	if err != nil {
		return err
	}

	_, err = EnsureDNSMasq()
	if err != nil {
		return fmt.Errorf("failed to ensure dnsmasq: %w", err)
	}

	_, err = EnsureDB()
	if err != nil {
		return fmt.Errorf("failed to ensure db: %w", err)
	}

	_, err = EnsureProxy(proxyOpts{
		CertDir: opts.CertDir,
	})
	if err != nil {
		return fmt.Errorf("failed to ensure proxy: %w", err)
	}

	network, err := dock.GetNetwork(client, networkProxy)
	if err != nil {
		return fmt.Errorf("failed to get network: %w", err)
	}

	_, err = dock.EnsureContainer(client, dock.ContainerOptions{
		Name:           opts.Name,
		Image:          imageWP,
		PrimaryNetwork: network,
		Binds: []string{
			fmt.Sprintf("%s:/var/www/html", opts.LocalPath),
			fmt.Sprintf("%s:/wpclone/sshkey", opts.SSHKeyPath),
		},
		Labels: map[string]string{
			"wpclone_type": "wp",
			"wpclone_url":  opts.URL,
			"wpclone_fqdn": opts.FQDN,
			"wpclone_ssl":  fmt.Sprintf("%t", opts.SSLEnabled),
		},
		RestartPolicy: "unless-stopped",
	})
	if err != nil {
		return fmt.Errorf("failed to ensure container: %w", err)
	}

	return nil
}

func EnsureDNSMasq() (*docker.APIContainers, error) {
	client, err := dock.GetClient()
	if err != nil {
		return nil, err
	}

	network, err := dock.EnsureNetwork(client, networkProxy)
	if err != nil {
		return nil, err
	}

	container, err := dock.EnsureContainer(client, dock.ContainerOptions{
		Name:           config.ContainerName("dnsmasq"),
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

	log.Debugf("Container %s running", config.ContainerName("dnsmasq"))

	return container, nil
}

type proxyOpts struct {
	CertDir string
}

func EnsureProxy(opts proxyOpts) (*docker.APIContainers, error) {
	client, err := dock.GetClient()
	if err != nil {
		return nil, err
	}

	network, err := dock.EnsureNetwork(client, networkProxy)
	if err != nil {
		return nil, err
	}

	volumeCaddy, err := dock.EnsureVolume(client, volumeProxyCaddy)
	if err != nil {
		return nil, err
	}

	volumeData, err := dock.EnsureVolume(client, volumeProxyData)
	if err != nil {
		return nil, err
	}

	volumeConfig, err := dock.EnsureVolume(client, volumeProxyConfig)
	if err != nil {
		return nil, err
	}

	container, err := dock.EnsureContainer(client, dock.ContainerOptions{
		Name:           config.ContainerName("proxy"),
		Image:          imageProxy,
		PrimaryNetwork: network,
		Binds: []string{
			volumeData.Name + ":/data",
			volumeConfig.Name + ":/config",
			volumeCaddy.Name + ":/etc/caddy",
			opts.CertDir + ":/data/caddy/pki/authorities/local",
		},
		Ports: map[docker.Port][]docker.PortBinding{
			"80/tcp": {
				{
					HostPort: "80",
				},
			},
			"443/tcp": {
				{
					HostPort: "443",
				},
			},
		},
		Labels: map[string]string{
			"wpclone_type": "proxy",
		},
		RestartPolicy: "unless-stopped",
	})
	if err != nil {
		return nil, err
	}

	log.Debugf("Container %s running", config.ContainerName("proxy"))

	return container, nil
}
