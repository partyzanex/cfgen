package main

import (
	"log"
	"os"

	config "github.com/partyzanex/cli-config-gen"
	"github.com/urfave/cli/v2"
)

const (
	sourceFileFlag   = "source"
	targetPathFlag   = "target"
	packageFlag      = "package"
	templatePathFlag = "template"
)

func main() {
	app := new(cli.App)
	app.Usage = "cli tool for generates config package from YAML"
	app.Action = action
	app.Flags = []cli.Flag{
		&cli.PathFlag{
			Name:       sourceFileFlag,
			Aliases:    []string{"s", "src"},
			Usage:      "Path to source config.yaml file",
			Required:   true,
			Value:      "./config.yaml",
			HasBeenSet: true,
		},
		&cli.PathFlag{
			Name:       targetPathFlag,
			Aliases:    []string{"t"},
			Usage:      "Path to target directory",
			Required:   true,
			Value:      "./internal/config/config.go",
			HasBeenSet: true,
		},
		&cli.StringFlag{
			Name:       packageFlag,
			Aliases:    []string{"p", "pkg"},
			Usage:      "Target go package name",
			Required:   true,
			Value:      "config",
			HasBeenSet: true,
		},
		&cli.PathFlag{
			Name:    templatePathFlag,
			Aliases: []string{"tpl"},
			Usage:   "Path to template file",
			Value:   "",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func action(ctx *cli.Context) error {
	gen := &config.Codegen{
		TemplatePath: ctx.Path(templatePathFlag),
		SourceFile:   ctx.Path(sourceFileFlag),
		TargetPath:   ctx.Path(targetPathFlag),
		PackageName:  ctx.String(packageFlag),
	}

	return gen.Run()
}
