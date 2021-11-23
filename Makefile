PREFIX ?= /usr/local

.PHONY: build
build:
	go build -o dist/local/mod cmd/main.go

.PHONY: install
install: build
	@cp dist/local/mod $(PREFIX)/bin/mod
	@chmod 755 $(PREFIX)/bin/mod
	@mod --version
