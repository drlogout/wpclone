package docker

import "croox/wpclone/pkg/dock"

func RemoveAllNetworks() error {
	filters := map[string]map[string]bool{
		"label": {
			networkProxy: true,
		},
	}

	return dock.RemoveNetworks(filters)
}
