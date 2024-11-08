package docker_cmd

import (
	"croox/wpclone/cmd/common"
	"croox/wpclone/config"
	"croox/wpclone/docker"
	"croox/wpclone/pkg/dock"
	"croox/wpclone/pkg/exec"
	"croox/wpclone/pkg/message"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/urfave/cli/v2"
)

var Docker = &cli.Command{
	Name:  "docker",
	Usage: "Docker commands (ps, start, stop, remove, down, cert, db-only)",
	Before: common.BeforeCmds([]cli.BeforeFunc{
		common.BeforeCheckDocker,
	}),
	Action: func(ctx *cli.Context) error {
		cli.ShowAppHelpAndExit(ctx, 0)
		return nil
	},
	Subcommands: []*cli.Command{
		{
			Name:  "ps",
			Usage: "List running wp containers",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "all",
					Usage:   "List all containers",
					Aliases: []string{"a"},
				},
			},
			Action: func(ctx *cli.Context) error {
				containers, err := listContainers(ctx.Bool("all"))
				if err != nil {
					return err
				}
				if len(containers) == 0 {
					if ctx.Bool("all") {
						return message.Exitf("No containers running")
					}
					return message.Exit("No wordpress containers running")
				}

				t := table.NewWriter()
				t.SetOutputMirror(os.Stdout)
				t.AppendHeader(table.Row{"#", "Local URL", "Container Name", "Type"})

				for i, container := range containers {
					containerName := strings.TrimLeft(container.ContainerName, "/")
					t.AppendRow(table.Row{i + 1, container.URL, containerName, container.Type})
				}

				t.AppendSeparator()
				t.Render()

				return nil
			},
		},
		{
			Name:  "start",
			Usage: "Start wordpress container",
			Before: common.BeforeCmds([]cli.BeforeFunc{
				common.BeforeLoadConfig,
				common.BeforeCheckDockerWebPorts,
				common.BeforeCheckDockerDBPorts,
			}),
			Action: func(ctx *cli.Context) error {
				cfg := common.ConfigFromCTX(ctx)

				if !cfg.RunInDocker() {
					return message.Exit("Docker not enabled in wpclone.yaml")
				}

				running, err := dock.IsWPRunning(cfg.DockerWPContainerName())
				if err != nil {
					return err
				}

				if running {
					return message.Exitf("Container %s is already running", cfg.LocalURL())
				}

				opts := dock.WPOptions{
					Name:       cfg.DockerWPContainerName(),
					LocalPath:  cfg.LocalPath(),
					SSHKeyPath: cfg.SSHKeyPath(),
					URL:        cfg.LocalURL(),
					FQDN:       cfg.LocalFQDN(),
					CertDir:    cfg.CertDirPath(),
					SSLEnabled: cfg.DockerSSLEnabled(),
				}
				if err := dock.EnsureWP(opts); err != nil {
					return err
				}

				if err := dock.UpdateProxy(); err != nil {
					return err
				}

				message.Successf("Successfully started %s", cfg.LocalURL())
				return nil
			},
		},
		{
			Name:  "stop",
			Usage: "Stop wordpress container",
			Before: common.BeforeCmds([]cli.BeforeFunc{
				common.BeforeLoadConfig,
				common.BeforeCheckDockerWebPorts,
				common.BeforeCheckDockerDBPorts,
			}),
			Action: func(ctx *cli.Context) error {
				cfg := common.ConfigFromCTX(ctx)

				running, err := dock.IsWPRunning(cfg.DockerWPContainerName())
				if err != nil {
					return err
				}

				if !running {
					return message.Exitf("Container %s is not running", cfg.LocalURL())
				}

				if err := docker.StopWP(cfg); err != nil {
					return err
				}

				message.Successf("Successfully stopped %s", cfg.LocalURL())
				return nil
			},
		},
		{
			Name:    "remove",
			Aliases: []string{"rm"},
			Usage:   "Remove wp container and delete database",
			Before: common.BeforeCmds([]cli.BeforeFunc{
				common.BeforeLoadConfig,
			}),
			Action: func(ctx *cli.Context) error {
				cfg := common.ConfigFromCTX(ctx)

				site, err := docker.RemoveWP(cfg)
				if err != nil {
					return err
				}

				if site == nil {
					return message.Exitf("Docker app %s does not exist", cfg.LocalFQDN())
				}
				message.Successf("Successfully removed %s | URL: %s", cfg.LocalFQDN(), cfg.LocalURL())
				return nil
			},
		},
		{
			Name:  "down",
			Usage: "Shutdown all containers (all WP containers, proxy, database, dnsmasq)",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "volumes",
					Usage:   "Remove volumes",
					Aliases: []string{"v"},
				},
			},
			Action: func(ctx *cli.Context) error {
				containers, err := removeWPClone()
				if err != nil {
					return err
				}

				if ctx.Bool("volumes") {
					if err := dock.RemoveAllVolumes(); err != nil {
						return err
					}
				}

				if len(containers) == 0 {
					return message.Exit("No containers to remove")
				}

				for _, container := range containers {
					message.Successf("Removed %s", container.ContainerName)
				}

				message.Successf("Successfully removed all containers")
				return nil
			},
		},
		{
			Name:  "cert",
			Usage: "Print path to root certificate",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "open",
					Usage: "Open certificate",
				},
			},
			Action: func(ctx *cli.Context) error {
				certPath := filepath.Join(config.CertDirPath(), "root.crt")

				message.Titlef("Certificate path is: %s", certPath)
				message.Info("run `wpclone docker cert --open` to open the certificate")
				message.Info("After that make sure you have trusted the certificate")

				if ctx.Bool("open") {
					if runtime.GOOS == "darwin" {
						return exec.Run("open", certPath)
					}
				}

				return nil
			},
		},
		{
			Name:  "db-only",
			Usage: "Run only the database container, stop all other containers",
			Before: common.BeforeCmds([]cli.BeforeFunc{
				common.BeforeCheckDockerDBPorts,
			}),
			Action: func(ctx *cli.Context) error {
				_, err := dock.RemoveAllContainersExceptDB()
				if err != nil {
					return err
				}

				_, err = dock.EnsureDB()
				if err != nil {
					return err
				}

				return nil
			},
		},
		{
			Name:  "login",
			Usage: "Login to the wordpress container",
			Before: common.BeforeCmds([]cli.BeforeFunc{
				common.BeforeLoadConfig,
				common.BeforeCheckDockerWebPorts,
				common.BeforeCheckDockerDBPorts,
			}),
			Action: func(ctx *cli.Context) error {
				cfg := common.ConfigFromCTX(ctx)

				opts := exec.RunOpts{
					Interactive: true,
				}

				if err := exec.RunWithOpts(opts, "docker", "exec", "-it", "-u", "www-data", cfg.DockerWPContainerName(), "/bin/bash"); err != nil {
					return err
				}

				return nil
			},
		},
	},
}
