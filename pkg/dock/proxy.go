package dock

import (
	"embed"
	"html/template"
	"os"

	docker "github.com/fsouza/go-dockerclient"
	log "github.com/sirupsen/logrus"
)

//go:embed *.tmpl
var tpls embed.FS

type proxyOpts struct {
	CertDir string
}

func ensureProxy(opts proxyOpts) (*docker.APIContainers, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	network, err := ensureNetwork(client, networkProxy)
	if err != nil {
		return nil, err
	}

	volumeCaddy, err := ensureVolume(client, volumeProxyCaddy)
	if err != nil {
		return nil, err
	}

	volumeData, err := ensureVolume(client, volumeProxyData)
	if err != nil {
		return nil, err
	}

	volumeConfig, err := ensureVolume(client, volumeProxyConfig)
	if err != nil {
		return nil, err
	}

	container, err := ensureContainer(client, ContainerOptions{
		Name:           ContainerName("proxy"),
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

	log.Debugf("Container %s running", ContainerName("proxy"))

	return container, nil
}

func UpdateProxy() error {
	client, err := getClient()
	if err != nil {
		return err
	}

	wpContainers, err := ListContainers("wp")
	if err != nil {
		return err
	}

	tempFile, err := writeCaddyfile(wpContainers)
	if err != nil {
		return err
	}

	if err := installCaddyfile(client, tempFile); err != nil {
		return err
	}

	os.RemoveAll(tempFile)

	if err := restartContainer(client, ContainerName("proxy")); err != nil {
		return err
	}

	log.Debug("Proxy updated")

	return nil
}

func writeCaddyfile(wpContainers []WPCloneContainerInfo) (string, error) {
	t, err := template.ParseFS(tpls, "*")
	if err != nil {
		return "", err
	}

	tempFile, err := os.CreateTemp("", "Caddyfile")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	if err := t.ExecuteTemplate(tempFile, "Caddyfile.tmpl", wpContainers); err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}

func installCaddyfile(client *docker.Client, tmpCaddyFile string) error {
	volumeCaddy, err := ensureVolume(client, volumeProxyCaddy)
	if err != nil {
		return err
	}

	err = run(runOptions{
		Name:  ContainerNameWithID("rsync"),
		Image: imageRsync,
		Binds: []string{
			tmpCaddyFile + ":/Caddyfile",
			volumeCaddy.Name + ":/etc/caddy",
		},
		Cmd: []string{"/Caddyfile", "/etc/caddy/Caddyfile"},
		Labels: map[string]string{
			"wpclone_type":      "rsync",
			"wpclone_ephimeral": "true",
		},
	})
	if err != nil {
		return err
	}

	return nil
}
