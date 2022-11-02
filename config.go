package config

import "embed"

//go:embed config.tpl
var TemplateFS embed.FS
