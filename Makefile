VERSION=$(shell git describe --tags --always)
COMMIT=$(shell git rev-parse HEAD)
BUILD=$(shell date +%FT%T%z)
PKG=github.com/eBay/sbom-scorecard


LDFLAGS="-X $(PKG).version=$(VERSION) -X $(PKG).commit=$(COMMIT) -X $(PKG).date=$(BUILD)"

.DEFAULT_GOAL := build

.PHONY: build
build: ## Build a version
	go build -ldflags ${LDFLAGS} -o bin/sbom-scorecard cmd/sbom-scorecard/main.go

.PHONY: clean
clean: ## Remove temporary files
	go clean

.PHONY: test
test: ## Run the unit tests
	echo 'mode: atomic' > coverage.txt && go test -covermode=atomic -coverprofile=coverage.txt -v -race -timeout=30s ./...

phony:
	@echo Use specific targets to download individual needed files.

examples/julia.spdx.json:
	curl -Lo examples/julia.spdx.json https://github.com/JuliaLang/julia/raw/master/julia.spdx.json

examples/dropwizard.cyclonedx.json:
	curl -Lo examples/dropwizard.cyclonedx.json https://github.com/CycloneDX/bom-examples/raw/master/SBOM/dropwizard-1.3.15/bom.json
