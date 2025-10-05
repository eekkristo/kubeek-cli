package commands

import (
	"context"
	"maps"
	"os"
	"strings"

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
			&cli.BoolFlag{Name: "interactive", Usage: "prompt for placeholder values (default true)", Value: true},
			&cli.BoolFlag{Name: "force", Aliases: []string{"f"}, Usage: "overwrite destination folder without prompting"},
			&cli.StringSliceFlag{Name: "defaults", Aliases: []string{"d"}, Usage: "pass argumants directly via cli"},
			&cli.StringFlag{Name: "template", Usage: "name of the folder that is taken as an input", Value: "./templates"},
			&cli.StringFlag{Name: "name", Usage: "name of the folder to be created (Rquired)", Value: "generated", Required: true},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			ac := config.DefaultAppConfig()
			var err error
			var interactive = true

			if c.Bool("no-config") {
				ac = config.AppConfig{Placeholders: config.Config{}}
			} else if len(c.StringSlice("defaults")) > 0 {
				maps.Copy(ac.Placeholders, parseDefaults(strings.Join(c.StringSlice("defaults"), ",")))
				interactive = false
			} else if c.String("config") != "" {
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

			dst, err := resolveGenerateDest(c)
			if err != nil {
				return err
			}

			return generate.Render(c.String("template"), dst, ac, interactive, c.Bool("force"))
		},
	}
}
