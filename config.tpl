package {{.PackageName}}

import (
"fmt"
"time"

"github.com/urfave/cli/v2"
. "github.com/partyzanex/cli-config-gen"
)

// Description
const (
Name = "{{.App.Name}}"
Desc = "{{.App.Desc}}"
)

// Env
const (
{{range .App.Env}}Env{{toCamel .String}} EnvName = "{{.String}}"
{{end}}
)

// Flag names
const (
EnvFlagName = "env"
{{range .Flags}}{{toCamel .Name}}FlagName = "{{.Name}}"
{{end}}
)

{{range .Flags}}
    {{if eq .Type "enum"}}
        {{$flagName := toCamel .Name}}
        // {{$flagName}} enums
        const (
        {{range .Enum}}{{$flagName}}{{toCamel .}} = "{{.}}"
        {{end}}
        )
    {{end}}
{{end}}
// Env should be setup the default environment name.
var Env EnvName

func init() {
Env = GetEnvName("{{toSnake $.App.Name}}_ENV", {{range $.App.Env}}Env{{toCamel .String}},{{end}})
}

// Flag values
var ({{range .Flags}}
  // {{toCamel .Name}} contains default environments values.
  {{$flag := .}}{{toCamel .Name}} = NewValue(Env){{range $.App.Env}}.
  {{$flag.ValueSetMethodName}}(Env{{toCamel .String}}, {{$flag.Args .String}}){{end}}
{{end}}
)

// EnvFlag returns *cli.StringFlag for --env flag.
func EnvFlag() *cli.StringFlag {
return &cli.StringFlag{
Name:        EnvFlagName,
Category:    "",
DefaultText: "",
FilePath:    "",
Usage:       "Environment name",
Required:    false,
Hidden:      false,
HasBeenSet:  false,
Value:       Env.String(),
Destination: nil,
Aliases:     nil,
EnvVars:     []string{"{{toSnake $.App.Name}}_ENV"},
TakesFile:   false,
Action: func(_ *cli.Context, s string) error {
env := EnvName(s)

switch env {
{{range $.App.Env}}case Env{{toCamel .String}}:
{{end}}default:
return fmt.Errorf("invalid environment name %q", s)
}

Env = env

return nil
},
}
}

{{range .Flags}}
  // {{toCamel .Name}}Flag returns a *cli.{{.CLIFlagType}} for --{{.Name}} flag.
  func {{toCamel .Name}}Flag() *cli.{{.CLIFlagType}} {
  return &cli.{{.CLIFlagType}}{
  Name:        {{toCamel .Name}}FlagName,
  Usage:       {{quote .Desc}},
  Required:    {{.RequiredField}},
  Value:       {{toCamel .Name}}.{{.ValueType}}(),
  EnvVars:     {{.EnvVarsField $.App.Name}},
  {{ if eq .Type.String "timestamp"}}
    Action: func(_ *cli.Context, v *time.Time) error {
    if v != nil {
    Start.SetTimestamp(Env, v.Format(time.RFC3339))
    }

    return nil
    },
  {{else}}
    Action: func(_ *cli.Context, v {{.GoType}}) error {
    {{toCamel .Name}}.{{.ValueSetMethodName}}(Env, v{{if .IsSlice}}...{{end}})

    return nil
    },
  {{end}}
  }
  }
{{end}}

func CLIFlags() []cli.Flag {
return []cli.Flag{
{{range .Flags}}{{toCamel .Name}}Flag(),
{{end}}
}
}
