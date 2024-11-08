package cmd

import (
	dbCmd "croox/wpclone/cmd/db_cmd"
	dockerCmd "croox/wpclone/cmd/docker_cmd"
	projectCmd "croox/wpclone/cmd/project_cmd"
	"embed"

	"github.com/urfave/cli/v2"
)

//go:embed *.tmpl
var tpls embed.FS

var Commands = []*cli.Command{
	Init,
	Pull,
	Push,
	dockerCmd.Docker,
	dbCmd.DB,
	projectCmd.Project,
	Setup,
	Install,
}
