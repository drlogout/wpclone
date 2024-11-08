package dock

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
)

func ListContainers(flag ...string) ([]WPCloneContainerInfo, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	containers, err := client.ListContainers(docker.ListContainersOptions{
		All:     true,
		Filters: getLabelFilter(flag...),
	})
	if err != nil {
		return nil, err
	}

	wpContainers := []WPCloneContainerInfo{}

	if len(containers) == 0 {
		return wpContainers, nil
	}

	for _, container := range containers {
		wpContainers = append(wpContainers, getWPCloneContainerInfo(container))
	}

	return wpContainers, nil
}

func ListVolumes(flag ...string) ([]docker.Volume, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	options := docker.ListVolumesOptions{
		Filters: getLabelFilter(flag...),
	}

	volumeList, err := client.ListVolumes(options)
	if err != nil {
		return nil, err
	}

	return volumeList, nil
}

func ListNetworks(flag ...string) ([]docker.Network, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	filters := map[string]map[string]bool{
		"label": {
			"wpclone_type": true,
		},
	}

	options := docker.NetworkFilterOpts(filters)

	networks, err := client.FilteredListNetworks(options)
	if err != nil {
		return nil, err
	}

	return networks, nil
}

func ListDBs() ([]string, error) {
	var buffer bytes.Buffer
	var writer io.Writer = &buffer
	dbs := []string{}

	_, err := EnsureDB()
	if err != nil {
		return dbs, err
	}

	opts := ExecOptions{
		ContainerName: ContainerName("db"),
		Cmd:           []string{"dbs-list"},
		Stdout:        writer,
	}

	status, err := Exec(opts)
	if err != nil {
		return dbs, err
	}

	if status != 0 {
		return dbs, fmt.Errorf("failed to list databases")
	}

	output := strings.Split(buffer.String(), "\n")
	for _, db := range output {
		if db != "" && noSystemDB(db) {
			dbs = append(dbs, db)
		}
	}

	return dbs, nil
}

func getLabelFilter(flag ...string) map[string][]string {
	var filterValue string
	if len(flag) > 0 {
		filterValue = flag[0]
	}

	switch filterValue {
	case "wp":
		return map[string][]string{
			"label": {"wpclone_type=wp"},
		}
	case "db":
		return map[string][]string{
			"label": {"wpclone_type=db"},
		}
	case "proxy":
		return map[string][]string{
			"label": {"wpclone_type=proxy"},
		}
	case "dnsmasq":
		return map[string][]string{
			"label": {"wpclone_type=dnsmasq"},
		}
	case "ephimeral":
		return map[string][]string{
			"label": {"wpclone_ephimeral=true"},
		}
	default:
		return map[string][]string{
			"label": {"wpclone_type"},
		}
	}

}

func noSystemDB(db string) bool {
	return db != "information_schema" && db != "performance_schema" && db != "mysql" && db != "sys"
}
