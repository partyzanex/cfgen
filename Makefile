include go.mk

go.mk:
	@tmpdir=$$(mktemp -d) && \
	git clone --depth 1 --single-branch https://github.com/partyzanex/go-makefile.git $$tmpdir && \
	cp $$tmpdir/go.mk $(CURDIR)/go.mk

CLI_CONFIG_GEN_BIN := $(LOCAL_BIN)/cli-config-gen

.PHONY: build
build: bin-default
	@go build -o $(CLI_CONFIG_GEN_BIN) ./cmd/cli-config-gen

.PHONY: example-config
example-config: build
	$(CLI_CONFIG_GEN_BIN) -s config.example.yaml -t ./internal/config/config.go
	go fmt ./internal/config/config.go
