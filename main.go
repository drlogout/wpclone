package main

import (
	"croox/wpclone/cmd"
	"croox/wpclone/cmd/common"
	"croox/wpclone/pkg/dock"
	"croox/wpclone/pkg/exec"
	"croox/wpclone/pkg/message"
	sshexec "croox/wpclone/pkg/sshexec"
	"fmt"

	_ "embed"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})

	cli.AppHelpTemplate = additionalHelp()
}

func main() {
	app := &cli.App{
		Name:  "wpclone",
		Usage: "Clone WordPress sites (configure in wpclone.yml)",
		Action: func(ctx *cli.Context) error {
			cli.ShowAppHelpAndExit(ctx, 0)
			return nil
		},
		Before: func(ctx *cli.Context) error {
			if err := common.BeforeCheckSystem(ctx); err != nil {
				return err
			}

			if ctx.Bool("verbose") {
				sshexec.Verbose = true
				exec.Verbose = true
				dock.Verbose = true
			}

			if ctx.Bool("debug") {
				log.SetLevel(log.DebugLevel)
			} else {
				log.SetLevel(log.InfoLevel)
			}

			if ctx.Bool("quiet") {
				log.SetLevel(log.FatalLevel)
				message.SetQuiet()
			}

			return nil
		},
		Commands: cmd.Commands,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Usage:   "Specify wpclone.yml file",
				Aliases: []string{"c"},
			},
			&cli.StringFlag{
				Name:    "project",
				Usage:   "Specify project name (see wpclone project ls)",
				Aliases: []string{"p"},
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Usage:   "Verbose output",
				Value:   false,
				Aliases: []string{"v"},
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Debug output",
				Value: false,
			},
			&cli.BoolFlag{
				Name:    "quiet",
				Aliases: []string{"q"},
				Usage:   "no output",
				Value:   false,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func additionalHelp() string {
	additionalHelp := fmt.Sprintf(`%s
ENV VARS: 
	 WPCLONE_PROJECT_DIR: 			 		 Specify project directory (default: $HOME/wpclone)
	`, cli.AppHelpTemplate)

	return additionalHelp
}
