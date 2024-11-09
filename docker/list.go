package docker

import "croox/wpclone/pkg/dock"

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

func ListContainers(all bool) ([]WPCloneContainer, error) {
	if all {
		return ListAllContainers()
	}

	return ListWPContainers()
}

func ListAllContainers() ([]WPCloneContainer, error) {
	apiContainers, err := dock.ListContainers(map[string][]string{})
	if err != nil {
		return nil, err
	}

	return getWPCloneContainers(apiContainers), nil
}

func ListWPContainers() ([]WPCloneContainer, error) {
	apiContainers, err := dock.ListContainers(getLabelFilter("wp"))
	if err != nil {
		return nil, err
	}

	return getWPCloneContainers(apiContainers), nil
}
