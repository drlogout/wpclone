package exec

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

var Verbose bool

type RunOpts struct {
	Dir         string
	Stdout      io.Writer
	Interactive bool
	Verbose     bool
	Env         map[string]string
}

func Run(name string, arg ...string) error {
	return RunWithOpts(RunOpts{}, name, arg...)
}

func RunWithOpts(opts RunOpts, name string, arg ...string) error {
	log.Debugf("[system.cmd]: %s %s", name, strings.Join(arg, " "))

	if Verbose {
		opts.Verbose = true
	}

	cmd := exec.Command(name, arg...)

	if opts.Dir != "" {
		cmd.Dir = opts.Dir
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		cmd.Dir = cwd
	}

	if opts.Env != nil {
		for k, v := range opts.Env {
			cmd.Env = append(cmd.Env, k+"="+v)
		}
	}

	var stdout bytes.Buffer
	var stdoutWriter io.Writer = &stdout
	var stderr bytes.Buffer
	var stderrWriter io.Writer = &stderr

	if opts.Verbose {
		stderrWriter = io.MultiWriter(stderrWriter, os.Stderr)
		stdoutWriter = io.MultiWriter(stdoutWriter, os.Stdout)
	}

	if opts.Stdout != nil {
		stdoutWriter = io.MultiWriter(stdoutWriter, opts.Stdout)
	}

	cmd.Stderr = stderrWriter
	cmd.Stdout = stdoutWriter

	if opts.Interactive {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
	}

	if err := cmd.Run(); err != nil {
		c := fmt.Sprintf("%s %s", name, strings.Join(arg, " "))
		return fmt.Errorf("[system.cmd] error: %s\n%w\n%s\n%s", c, err, stderr.String(), stdout.String())
	}

	return nil
}
