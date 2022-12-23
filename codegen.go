package config

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type node struct {
	*Source

	PackageName string
	SourceFile  string
	Version     string
}

type Codegen struct {
	source *Source

	TemplatePath string
	PackageName  string
	SourceFile   string
	TargetPath   string
}

func (g *Codegen) Run() error {
	err := g.readSource()
	if err != nil {
		return err
	}

	tpl, err := g.readTemplate()
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(g.TargetPath), os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "cannot make target path")
	}

	targetFile, err := os.Create(g.TargetPath)
	if err != nil {
		return errors.Wrapf(err, "cannot create target file %s", g.TargetPath)
	}

	err = tpl.Execute(targetFile, &node{
		Source:      g.source,
		PackageName: g.PackageName,
		SourceFile:  g.SourceFile,
	})
	if err != nil {
		return errors.Wrap(err, "cannot execute config template")
	}

	return nil
}

func (g *Codegen) readTemplate() (*template.Template, error) {
	var (
		tr  io.Reader
		err error
	)

	if g.TemplatePath != "" {
		tr, err = os.Open(g.TemplatePath)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot open template file %s", g.TemplatePath)
		}

	} else {
		tr, err = TemplateFS.Open("config.tpl")
		if err != nil {
			return nil, errors.Wrap(err, "cannot open template")
		}
	}

	b, err := io.ReadAll(tr)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read file")
	}

	tpl, err := template.New(g.PackageName).Funcs(template.FuncMap{
		"toCamel":          strcase.ToCamel,
		"toSnake":          strcase.ToScreamingSnake,
		"quote":            strconv.Quote,
		"hasDateTimeFlags": g.source.Flags.HasDateTimeFlags,
	}).Parse(string(b))
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse template")
	}

	return tpl, nil
}

func (g *Codegen) readSource() error {
	src, err := os.Open(g.SourceFile)
	if err != nil {
		return errors.Wrap(err, "cannot open source config file")
	}

	source := new(Source)

	err = yaml.NewDecoder(src).Decode(source)
	if err != nil {
		return errors.Wrap(err, "cannot decode source file")
	}

	g.source = source

	return nil
}
