package dock

import (
	"fmt"
)

const (
	wpContainerStateMissing = "missing"
	wpContainerStateDeleted = "deleted"
)

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
	client, err := getClient()
	if err != nil {
		return err
	}

	_, err = ensureDNSMasq()
	if err != nil {
		return fmt.Errorf("failed to ensure dnsmasq: %w", err)
	}

	_, err = EnsureDB()
	if err != nil {
		return fmt.Errorf("failed to ensure db: %w", err)
	}

	_, err = ensureProxy(proxyOpts{
		CertDir: opts.CertDir,
	})
	if err != nil {
		return fmt.Errorf("failed to ensure proxy: %w", err)
	}

	network, err := getNetwork(client, networkProxy)
	if err != nil {
		return fmt.Errorf("failed to get network: %w", err)
	}

	_, err = ensureContainer(client, ContainerOptions{
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

func StopWP(name string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	if err := stopAndRemoveContainer(client, name); err != nil {
		return err
	}

	return nil
}

func IsWPRunning(name string) (bool, error) {
	client, err := getClient()
	if err != nil {
		return false, err
	}

	container, err := getContainer(client, name)
	if err != nil {
		return false, err
	}

	if container == nil {
		return false, nil
	}

	return true, nil
}
