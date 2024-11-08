package wp

import (
	"net/url"
	"strings"
)

type URLVariant struct {
}

func URLVariants(u string) ([]string, error) {
	uObj, err := url.Parse(u)
	if err != nil {
		return []string{}, err
	}

	variants := []string{}

	variants = append(variants, "https://"+uObj.Host)
	variants = append(variants, "http://"+uObj.Host)

	if strings.HasPrefix(uObj.Host, "www.") {
		variants = append(variants, "https://"+strings.TrimPrefix(uObj.Host, "www."))
		variants = append(variants, "http://"+strings.TrimPrefix(uObj.Host, "www."))
	} else {
		variants = append(variants, "https://www."+uObj.Host)
		variants = append(variants, "http://www."+uObj.Host)
	}

	return variants, nil
}

func AppendSSLURLVariants(variants []string, u string, sslEnabled bool) ([]string, error) {
	uObj, err := url.Parse(u)
	if err != nil {
		return []string{}, err
	}

	if sslEnabled {
		variants = append(variants, "http://"+uObj.Host)
	} else {
		variants = append(variants, "https://"+uObj.Host)
	}

	return variants, nil
}
