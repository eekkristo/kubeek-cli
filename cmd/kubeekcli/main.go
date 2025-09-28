package main

import (
	"context"
	"fmt"
	"kubeekcli/cmd/kubeekcli/commands"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "kubeekcli",
		Usage: "Simple render / revert / template / generate with state tracking to speed to up rendering environments or templates.",

		Commands: commands.All(),
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
