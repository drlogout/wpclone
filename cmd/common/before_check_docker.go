package common

import (
	"bufio"
	"croox/wpclone/pkg/dock"
	"croox/wpclone/pkg/exec"
	"croox/wpclone/pkg/message"
	"croox/wpclone/pkg/util"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/miekg/dns"
	"github.com/urfave/cli/v2"
)

var requiredTCPPorts = []int{80, 443}
var requiredUDPPorts = []int{53}
var requiredDBPorts = []int{3306}

var BeforeCheckDockerWebPorts = func(ctx *cli.Context) error {
	var portsInUse bool

	resolverOK, err := isResolverInstalled()
	if err != nil {
		return message.Exit(err.Error())
	}

	if !resolverOK {
		message.Info("Resolver not installed, trying to install...")
		if err := installResolver(); err != nil {
			return err
		}

		resolverOK, err := isResolverInstalled()
		if err != nil {
			return message.Exit(err.Error())
		}

		if !resolverOK {
			return message.Exit("Failed to install resolver")
		}
	}

	ok, err := areWebPortsInUseByWpclone()
	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	if checkLocalDNSServer() {
		portsInUse = true
		message.Info("Port 53 is already in use. (Find process listening on 53 with 'lsof -i -P -n | grep 53')")
	}

	for _, port := range requiredTCPPorts {
		if isPortInUse(port) {
			portsInUse = true
			message.Infof("Port %d is already in use. (Find process listening on %d with 'lsof -i -P -n | grep %d')", port, port, port)
		}
	}

	if portsInUse {
		return message.Exit("Please make sure the ports are available before running docker commands.")
	}

	return nil
}

var BeforeCheckDockerDBPorts = func(ctx *cli.Context) error {
	var portsInUse bool

	ok, err := areDBPortsInUseByWpclone()
	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	for _, port := range requiredDBPorts {
		if isPortInUse(port) {
			portsInUse = true
			message.Infof("Port %d is already in use. (Find process listening on %d with 'lsof -i -P -n | grep %d')", port, port, port)
		}
	}

	if portsInUse {
		return message.Exit("Please make sure the ports are available before running db commands.")
	}

	return nil

}

var BeforeCheckDocker = func(ctx *cli.Context) error {
	if err := dock.Ping(); err != nil {
		return message.Exit("Docker is not running. Start docker and try again.")
	}

	return nil
}

var BeforeCheckConfig = func(ctx *cli.Context) error {
	cfg := ConfigFromCTX(ctx)

	if !cfg.RunInDocker() {
		return message.Exit("Docker not configured in wpclone.yaml")
	}

	return nil
}

func isPortInUse(port int, opts ...string) bool {
	protocol := "tcp"
	if len(opts) > 0 {
		protocol = opts[0]
	}

	// Try to listen on the given port.
	address := fmt.Sprintf("0.0.0.0:%d", port)
	listener, err := net.Listen(protocol, address)

	if err != nil {
		return true // Port is in use
	}

	// Close the listener if no error, meaning the port is not in use.
	listener.Close()
	return false
}

func checkLocalDNSServer() bool {
	// Create a new DNS client
	client := new(dns.Client)
	client.Timeout = 2 * time.Second // Set a timeout for the DNS query

	// Create a new DNS message
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn("croox.com"), dns.TypeA)

	// Specify the DNS server as localhost
	server := "127.0.0.1:53"

	// Send the DNS query to the localhost server
	_, _, err := client.Exchange(m, server)

	// If there's no error, the DNS server is running on localhost
	return err == nil
}

func areWebPortsInUseByWpclone() (bool, error) {
	ports, err := dock.Ports()
	if err != nil {
		return false, err
	}

	if !hasAllRequiredPorts(ports, requiredTCPPorts) {
		return false, nil
	}

	if !hasAllRequiredPorts(ports, requiredUDPPorts) {
		return false, nil
	}

	return true, nil
}

func areDBPortsInUseByWpclone() (bool, error) {
	ports, err := dock.Ports()
	if err != nil {
		return false, err
	}

	if !hasAllRequiredPorts(ports, requiredDBPorts) {
		return false, nil
	}

	return true, nil
}

func hasAllRequiredPorts(ports, requiredPorts []int) bool {
	for _, port := range requiredPorts {
		if !hasPort(ports, port) {
			return false
		}
	}

	return true
}

func hasPort(ports []int, port int) bool {
	for _, p := range ports {
		if p == port {
			return true
		}
	}

	return false
}

func isResolverInstalled() (bool, error) {
	resolverFile := "/etc/resolver/test"

	if !util.FileExists(resolverFile) {
		return false, nil
	}

	file, err := os.Open(resolverFile)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		lineCount++
		// Trim any whitespace and check the line content
		if strings.TrimSpace(scanner.Text()) != "nameserver 127.0.0.1" {
			return false, nil
		}
	}

	// Check that there was exactly one line
	if lineCount != 1 {
		return false, fmt.Errorf("unexpected number of lines in %s", resolverFile)
	}

	return true, scanner.Err()
}

func installResolver() error {
	return exec.Run("sudo", "bash", "-c", "echo nameserver 127.0.0.1 > /etc/resolver/test")
}
