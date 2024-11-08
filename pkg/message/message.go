package message

import (
	"fmt"
	"os"
	"time"

	"github.com/aripalo/go-delightful"
	"github.com/briandowns/spinner"
	"github.com/urfave/cli/v2"
)

var message delightful.Message

func init() {
	message = delightful.New("wpclone")
}

func SetQuiet() {
	message.SetSilentMode(true)
}

func Title(m string) {
	message.Titleln("", m)
}

func Titlef(format string, a ...any) {
	m := fmt.Sprintf(format, a...)
	Title(m)
}

func Info(m string) {
	message.Infoln("⋅", m)
}

func Infof(format string, a ...any) {
	m := fmt.Sprintf(format, a...)
	Info(m)
}

type InfoSpin struct {
	s *spinner.Spinner
}

func (i *InfoSpin) Start(m string) {
	i.s.Stop()
	i.s.Start()
	i.s.Suffix = "  " + m
}

func (i *InfoSpin) Stop(m string) {
	i.s.Stop()
	if m == "" {
		return
	}
	Info(m)
}

func NewInfoSpinner() InfoSpin {
	return InfoSpin{
		s: spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(os.Stderr)),
	}
}

func InfoSpinner(m string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Start()
	s.Suffix = "  " + m
	return s
}

func InfoSpinnerF(m string, worker func() error) error {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Start()
	Info(m)
	err := worker()
	s.Stop()
	return err
}

func Prompt(m string) {
	message.Prompt("📝", m)
}

func Promptf(format string, a ...any) {
	m := fmt.Sprintf(format, a...)
	Prompt(m)
}

func Success(m string) {
	message.Successln("✅", m)
}

func Successf(format string, a ...any) {
	m := fmt.Sprintf(format, a...)
	Success(m)
}

func Failure(m string) {
	message.Failureln("❌", m)
}

func Failuref(format string, a ...any) {
	m := fmt.Sprintf(format, a...)
	Failure(m)
}

func Exit(m string) error {
	Failure(m)
	return cli.Exit("", 1)
}

func Exitf(format string, a ...any) error {
	m := fmt.Sprintf(format, a...)
	Failure(m)
	return cli.Exit("", 1)
}

func ExitError(err error, m string) error {
	Failure(m)
	fmt.Println(err)
	return cli.Exit("", 1)
}

func ExitErrorf(err error, format string, a ...any) error {
	m := fmt.Sprintf(format, a...)
	Failure(m)
	fmt.Println(err)
	return cli.Exit("", 1)
}

func Sitef(format string, a ...any) {
	m := fmt.Sprintf(format, a...)
	message.Infoln("🌍", m)
}

func Serverf(format string, a ...any) {
	m := fmt.Sprintf(format, a...)
	message.Infoln("💻", m)
}

func DBf(format string, a ...any) {
	m := fmt.Sprintf(format, a...)
	message.Infoln("📋", m)
}
