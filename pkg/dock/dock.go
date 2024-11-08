package dock

import (
	docker "github.com/fsouza/go-dockerclient"
)

const (
	imageCLI     = "rg.nl-ams.scw.cloud/wpclone/wpclone-cli:latest"
	imageDB      = "rg.nl-ams.scw.cloud/wpclone/wpclone-db:latest"
	imageWP      = "rg.nl-ams.scw.cloud/wpclone/wpclone-wp:latest"
	imageDNSMasq = "rg.nl-ams.scw.cloud/wpclone/wpclone-dnsmasq:latest"
	imageRsync   = "ogivuk/rsync:latest"
	imageProxy   = "caddy:2.8.4-alpine"

	networkProxy = "wpclone_proxy"

	volumeDB          = "wpclone_db"
	volumeProxyCaddy  = "wpclone_proxy_caddy"
	volumeProxyData   = "wpclone_proxy_data"
	volumeProxyConfig = "wpclone_proxy_config"
)

var client *docker.Client

func getClient() (*docker.Client, error) {
	var err error

	if client != nil {
		return client, nil
	}

	client, err = docker.NewClientFromEnv()

	return client, err
}
