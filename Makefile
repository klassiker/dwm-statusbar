default: all

all: bin

.PHONY: bin
bin: fmt
	@CGO_ENABLED=0 go build -ldflags '-s -w' -o statusbar

.PHONY: dev
dev:
	@CGO_ENABLED=0 go build

.PHONY: fmt
fmt:
	find . -name '*.go' | xargs gofmt -w

.PHONY: test
test:
	gotestsum ./... -- -timeout=10m

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: install
install: test lint
	go install
	systemctl --user restart statusbar