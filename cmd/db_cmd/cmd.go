package dbCmd

import (
	"croox/wpclone/cmd/common"
	"croox/wpclone/docker"
	"croox/wpclone/pkg/exec"
	"croox/wpclone/pkg/message"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/urfave/cli/v2"
)

var DB = &cli.Command{
	Name:  "db",
	Usage: "DB commands (list, create, remove)",
	Before: common.BeforeCmds([]cli.BeforeFunc{
		common.BeforeCheckDockerDBPorts,
	}),
	Action: func(ctx *cli.Context) error {
		cli.ShowAppHelpAndExit(ctx, 0)
		return nil
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "name",
			Usage:   "Project name",
			Aliases: []string{"n"},
		},
	},
	Subcommands: []*cli.Command{
		{
			Name:    "list",
			Usage:   "List databases",
			Aliases: []string{"ls"},
			Action: func(ctx *cli.Context) error {
				dbs, err := docker.ListDBs()
				if err != nil {
					return err
				}

				if len(dbs) == 0 {
					return message.Exit("No databases running")
				}

				t := table.NewWriter()
				t.SetOutputMirror(os.Stdout)
				t.AppendHeader(table.Row{"#", "DB name", "DB user", "DB password", "DB host"})

				for i, db := range dbs {
					t.AppendRow(table.Row{i + 1, db, db, db, "127.0.0.1"})
				}

				t.AppendSeparator()
				t.Render()

				return nil
			},
		},
		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "Create databases",
			Action: func(ctx *cli.Context) error {
				dbName := ctx.Args().First()
				if dbName == "" {
					return message.Exit("database name is required")
				}

				exists, err := docker.DBExists(dbName)
				if err != nil {
					return err
				}

				if exists {
					return message.Exitf("Database %s already exists", dbName)
				}

				if err := docker.DBCreate(dbName); err != nil {
					return err
				}

				message.Successf("Database %s created!", dbName)
				fmt.Printf("Name:     %s\nUser:     %s\nPassword: %s\nHost:     %s\nPort:     %d", dbName, dbName, dbName, "127.0.0.1", 3306)

				return nil
			},
		},
		{
			Name:    "remove",
			Aliases: []string{"rm"},
			Usage:   "Remove databases",
			Action: func(ctx *cli.Context) error {
				dbName := ctx.Args().First()
				if dbName == "" {
					return message.Exit("database name is required")
				}

				exists, err := docker.DBExists(dbName)
				if err != nil {
					return err
				}

				if !exists {
					return message.Exitf("Database %s does not exists", dbName)
				}

				if err := docker.DBRemove(dbName); err != nil {
					return err
				}

				message.Successf("Database %s removed", dbName)

				return nil
			},
		},
		{
			Name:    "login",
			Aliases: []string{"l"},
			Usage:   "Login to database (use 'current' to login to current project database)",
			Action: func(ctx *cli.Context) error {
				dbName := ctx.Args().First()
				if dbName == "" {
					return message.Exit("database name is required")
				}

				if dbName == "current" {
					if err := common.BeforeLoadConfig(ctx); err != nil {
						return err
					}
					cfg := common.ConfigFromCTX(ctx)

					dbName = cfg.LocalDBName()
				}

				exists, err := docker.DBExists(dbName)
				if err != nil {
					return err
				}

				if !exists {
					return message.Exitf("Database %s does not exists", dbName)
				}

				opts := exec.RunOpts{
					Interactive: true,
					Env: map[string]string{
						"MYSQL_PWD": dbName,
					},
				}

				return exec.RunWithOpts(opts, "mysql", "-u", dbName, "-h", "127.0.0.1", dbName)
			},
		},
	},
}
