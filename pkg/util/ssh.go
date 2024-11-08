package util

import "fmt"

func SSHCommand(port int, args ...string) string {
	var keyPath string
	if len(args) > 0 {
		keyPath = args[0]
	}

	if keyPath != "" {
		return fmt.Sprintf("ssh -o StrictHostKeyChecking=no -i %s -p %d", keyPath, port)
	}

	return fmt.Sprintf("ssh -o StrictHostKeyChecking=no -p %d", port)
}
