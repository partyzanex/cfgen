# cli-config-gen

## Install

```shell
go install github.com/partyzanex/cli-config-gen@latest
```

## Usage

```shell
cli-config-gen -h
```

```
NAME:
   cli-config-gen - cli tool for generates config package from YAML

USAGE:
   cli-config-gen [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --source value, -s value, --src value   Path to source config.yaml file (default: "./config.yaml")
   --target value, -t value                Path to target directory (default: "./internal/config/config.go")
   --package value, -p value, --pkg value  Target go package name (default: "config")
   --template value, --tpl value           Path to template file
   --help, -h                              show help (default: false)

```

## Development

Build CLI app:
```shell
make build
```

Generate example config:
```shell
make example-config
```

Run:
```shell
go run ./internal/cmd/example -h
```

Run tests
```shell
make test
```
