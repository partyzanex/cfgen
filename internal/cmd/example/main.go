package main

import (
	"github.com/partyzanex/cli-config-gen/internal/config"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := new(cli.App)
	app.Usage = "Example app"
	app.Action = action
	app.Before = before
	app.After = after
	app.Flags = config.CLIFlags()
	app.Setup()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func action(ctx *cli.Context) error {
	return nil
}

func before(ctx *cli.Context) error {
	return nil
}

func after(ctx *cli.Context) error {
	return nil
}
