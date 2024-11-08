package sshexec

import (
	"fmt"
	"io"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

var Verbose bool

type RunOpts struct {
	Stdout      io.Writer
	Verbose     bool
	Dir         string
	SSHHost     string
	SSHPort     int
	SSHUser     string
	SSHPassword string
	SSHKeyPath  string
}

func RunWithOpts(opts RunOpts, cmd string, args ...string) error {
	args = append([]string{cmd}, args...)
	cmd = strings.Join(args, " ")

	if Verbose {
		opts.Verbose = true
	}

	auth, err := getAuthMethod(opts)
	if err != nil {
		return err
	}

	clientConfig := &ssh.ClientConfig{
		User: opts.SSHUser,
		Auth: []ssh.AuthMethod{
			auth,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// dial SSH connection
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", opts.SSHHost, opts.SSHPort), clientConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	// create new SSH session
	session, err := conn.NewSession()
	if err != nil {
		return err
	}

	if opts.Verbose {
		session.Stderr = os.Stderr
		session.Stdout = os.Stdout
	}

	if opts.Stdout != nil {
		session.Stdout = opts.Stdout
	}

	if opts.Dir != "" {
		cmd = fmt.Sprintf("cd %s && %s", opts.Dir, cmd)
	}

	log.Debugf("[system.ssh.cmd]: %v", cmd)

	return session.Run(cmd)
}

func getAuthMethod(opts RunOpts) (ssh.AuthMethod, error) {
	if opts.SSHPassword != "" {
		return ssh.Password(opts.SSHPassword), nil
	}

	identityFile := opts.SSHKeyPath

	// private key for auth
	key, err := os.ReadFile(identityFile)
	if err != nil {
		return nil, err
	}

	// create signer for auth
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	return ssh.PublicKeys(signer), nil
}
