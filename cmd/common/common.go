package common

import (
	"croox/wpclone/config"

	"github.com/gosimple/slug"
	"github.com/urfave/cli/v2"
)

func ConfigFromCTX(ctx *cli.Context) *config.Config {
	return ctx.App.Metadata["config"].(*config.Config)
}

func SaveConfigToCTX(ctx *cli.Context, cfg *config.Config) {
	ctx.App.Metadata["config"] = cfg
}

func ConfigFlags(ctx *cli.Context) config.Flags {
	return config.Flags{
		Verbose:           ctx.Bool("verbose"),
		Debug:             ctx.Bool("debug"),
		ConfigFilePath:    ctx.String("config"),
		Dir:               ctx.String("dir"),
		Project:           slug.Make(ctx.String("project")),
		SkipRsync:         ctx.Bool("skip-rsync"),
		Locale:            ctx.String("locale"),
		ListAllContainers: ctx.Command.Name == "list" && ctx.Bool("all"),
		AdminEmail:        ctx.String("admin-email"),
		AdminPassword:     ctx.String("admin-password"),
		PushWPConfig:      ctx.Bool("push-wp-config"),
	}
}

func BeforeCmds(funcs []cli.BeforeFunc) cli.BeforeFunc {
	return func(ctx *cli.Context) error {
		for _, f := range funcs {
			if err := f(ctx); err != nil {
				return err
			}
		}
		return nil
	}
}
