package docker

import "croox/wpclone/pkg/dock"

func RemoveAllVolumes() error {
	return dock.RemoveVolumes(getLabelFilter())
}
