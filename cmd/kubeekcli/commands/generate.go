package commands

import (
	"context"
	"maps"
	"os"

	cli "github.com/urfave/cli/v3"

	"kubeekcli/internal/config"
	"kubeekcli/internal/runner/generate"
)

func GenerateCmd() *cli.Command {
	return &cli.Command{
		Name:  "generate",
		Usage: "Interactively fill placeholders and generate a new folder.",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "config", Value: "config.json", Usage: "optional config file with defaults"},
			&cli.BoolFlag{Name: "no-config", Usage: "do not load a config file; use only --defaults and discovered placeholders"},
			&cli.BoolFlag{Name: "force", Aliases: []string{"f"}, Usage: "overwrite destination folder without prompting"},
			&cli.StringFlag{Name: "template", Usage: "Name of the folder that is taken as an input", Value: "./templates"},
			&cli.StringFlag{Name: "name", Usage: "Name of the folder to be created (Rquired)", Value: "generated", Required: true},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			ac := config.DefaultAppConfig()
			var err error

			if c.Bool("no-config") {
				ac = config.AppConfig{Placeholders: config.Config{}}
			} else {
				ac, _, err = config.LoadAppConfig(c.String("config"))
				if err != nil {
					if os.IsNotExist(err) {
						ac = config.AppConfig{Placeholders: config.Config{}}
					} else {
						return err
					}
				}
			}
			// Ensure the map is non-nil before merging (maps.Copy panics on nil dst)
			if ac.Placeholders == nil {
				ac.Placeholders = make(config.Config)
			}
			maps.Copy(ac.Placeholders, parseDefaults(c.String("defaults")))

			dst, err := resolveGenerateDest(c)
			if err != nil {
				return err
			}

			return generate.RunInteractive(c.String("template"), dst, ac, c.Bool("force"))
		},
	}
}
