package commands

import (
	"context"
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
			&cli.StringFlag{Name: "defaults", Value: "", Usage: `comma-separated defaults, e.g. "_clustername_=dev,_role_=tooling"`},

			&cli.StringFlag{Name: "template", Usage: "Name of the folder that is taken as an input", Value: "./templates"},
			&cli.StringFlag{Name: "name", Usage: "Name of the folder to be created (Rquired)", Value: "generated", Required: true},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			var ac config.AppConfig
			var err error

			if c.Bool("no-config") {
				ac = config.AppConfig{Placeholders: config.Config{}}
			} else {
				ac, _, err = config.LoadAppConfig(c.String("config"))
				// if file doesn't exist, fall back to empty config silently
				if err != nil {
					if os.IsNotExist(err) {
						ac = config.AppConfig{Placeholders: config.Config{}}
					} else {
						return err
					}
				}
			}
			// merge CLI defaults into placeholders (CLI has priority)
			if ac.Placeholders == nil {
				ac.Placeholders = config.Config{}
			}
			for k, v := range parseDefaults(c.String("defaults")) {
				ac.Placeholders[k] = v
			}

			opts := config.ResolveOpts(map[string]string{
				"exts":          c.String("exts"),
				"exclude-dirs":  c.String("exclude-dirs"),
				"exclude-files": c.String("exclude-files"),
			}, ac)

			dst, err := resolveGenerateDest(c)
			if err != nil {
				return err
			}

			var meta *bool
			if c.IsSet("meta") {
				v := c.Bool("meta")
				meta = &v
			}

			return generate.RunInteractive(c.String("template"), dst, ac, opts, c.Bool("force"), meta)
		},
	}
}
