package util

import "os"

func UserHome() string {
	userHome, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return userHome
}
