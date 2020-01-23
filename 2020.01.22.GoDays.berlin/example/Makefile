.PHONY: all deps build check clean

GO ?= go
GOLINT ?= golint

deps:
	$(GO) mod download

build_no_generate: deps
	$(GO) build -o bin/example

build: deps generate build_no_generate

check:
	$(GO) vet $$($(GO) list ./...)
	$(GOLINT) $$($(GO) list ./...)

generate:
	$(GO) generate models/secret.go

clean:
	rm -rf models/gen_* bin