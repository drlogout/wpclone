package docker

import (
	"embed"
	"os"
	"text/template"

	"croox/wpclone/config"
	"croox/wpclone/pkg/dock"

	docker "github.com/fsouza/go-dockerclient"
	log "github.com/sirupsen/logrus"
)

//go:embed *.tmpl
var tpls embed.FS

func UpdateProxy() error {
	client, err := dock.GetClient()
	if err != nil {
		return err
	}

	wpContainers, err := ListWPContainers()
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

	if err := dock.RestartContainer(client, config.ContainerName("proxy")); err != nil {
		return err
	}

	log.Debug("Proxy updated")

	return nil
}

func writeCaddyfile(wpContainers []WPCloneContainer) (string, error) {
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
	volumeCaddy, err := dock.EnsureVolume(client, volumeProxyCaddy)
	if err != nil {
		return err
	}

	err = dock.Run(dock.RunOptions{
		Name:  config.ContainerNameWithID("rsync"),
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
