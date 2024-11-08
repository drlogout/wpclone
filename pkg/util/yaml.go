package util

import (
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

var mutex sync.Mutex

func LoadYAML(path string, out interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, out); err != nil {
		return err
	}

	return nil
}
