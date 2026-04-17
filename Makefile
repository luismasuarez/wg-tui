build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o wg-tui .

install: build
	mkdir -p ~/.local/bin
	cp wg-tui ~/.local/bin/wg-tui

test:
	go test ./...

lint:
	go vet ./...

.PHONY: build install test lint
