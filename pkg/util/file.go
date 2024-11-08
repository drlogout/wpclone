package util

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}

		panic(err)
	}

	return true
}

func AbsPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	if strings.HasPrefix(path, "~/") {

		return filepath.Join(UserHome(), strings.TrimLeft(path, "~/"))
	}

	return filepath.Join(Getwd(), path)
}

func JoinPath(elem ...string) string {
	elems := normalizePath(elem...)
	return filepath.Join(elems...)
}

func normalizePath(elem ...string) []string {
	elems := []string{}
	for _, path := range elem {
		if strings.HasPrefix(path, "~/") {
			path = filepath.Join(UserHome(), strings.TrimLeft(path, "~/"))
		}
		elems = append(elems, path)
	}

	return elems
}

func Getwd() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return wd
}

func EnsureDir(path string) error {
	if !FileExists(path) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}

	return nil
}

func FolderEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// read in ONLY one file
	_, err = f.Readdir(1)

	// and if the file is EOF... well, the dir is empty.
	if err == io.EOF {
		return true, nil
	}

	return false, nil
}
