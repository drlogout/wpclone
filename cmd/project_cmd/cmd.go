package projectCmd

import (
	"croox/wpclone/cmd/common"
	"croox/wpclone/config"
	"croox/wpclone/pkg/message"
	"croox/wpclone/pkg/ternary"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"
)

var Project = &cli.Command{
	Name:   "project",
	Usage:  "Project specific commands (list)",
	Before: common.BeforeCmds([]cli.BeforeFunc{}),
	Action: func(ctx *cli.Context) error {
		cli.ShowAppHelpAndExit(ctx, 0)
		return nil
	},
	Subcommands: []*cli.Command{
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "List projects in wpclone project directory (default ~/wpclone)",
			Action: func(ctx *cli.Context) error {
				projectDir := config.ProjectDir()
				files, err := filepath.Glob(filepath.Join(projectDir, "*.wpclone.yml"))
				if err != nil {
					return fmt.Errorf("failed to list files: %w", err)
				}

				if len(files) == 0 {
					message.Info("No projects found")
					return nil
				}

				message.Title("Project directoy: " + projectDir)

				t := table.NewWriter()
				t.SetOutputMirror(os.Stdout)
				t.AppendHeader(table.Row{"#", "Project", "Local URL", "Local path", "Mode", "Remote URL"})

				for i, file := range files {
					flags := config.Flags{
						ConfigFilePath: file,
					}
					cfg, err := config.Load(flags)
					if err != nil {
						return fmt.Errorf("failed to load config: %w", err)
					}

					projectName := strings.TrimSuffix(filepath.Base(file), ".wpclone.yml")
					mode := ternary.String(cfg.RunInDocker(), "docker", "local")
					t.AppendRow(table.Row{i + 1, projectName, cfg.LocalURL(), cfg.LocalPath(), mode, cfg.RemoteURL()})
				}
				t.AppendSeparator()
				t.Render()

				return nil
			},
		},
	},
}
