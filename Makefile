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

examples: examples/julia.spdx.json examples/dropwizard.cyclonedx.json examples/openfeature-javasdk.cyclonedx.xml
examples/julia.spdx.json:
	curl -Lo examples/julia.spdx.json https://github.com/JuliaLang/julia/raw/master/julia.spdx.json
examples/dropwizard.cyclonedx.json:
	curl -Lo examples/dropwizard.cyclonedx.json https://github.com/CycloneDX/bom-examples/raw/master/SBOM/dropwizard-1.3.15/bom.json
examples/openfeature-javasdk.cyclonedx.xml:
	curl -Lo examples/openfeature-javasdk.cyclonedx.xml https://s01.oss.sonatype.org/content/repositories/snapshots/dev/openfeature/sdk/0.3.1-SNAPSHOT/sdk-0.3.1-20221014.132148-1-cyclonedx.xml

slsa: slsa/goreleaser-linux-amd64.yml slsa/goreleaser-linux-arm64.yml slsa/goreleaser-darwin-amd64.yml slsa/goreleaser-darwin-arm64.yml slsa/goreleaser-windows-amd64.yml slsa/goreleaser-windows-arm64.yml

slsa/goreleaser-linux-amd64.yml:
	GOOS=linux GOARCH=amd64 make TEMPLATE
slsa/goreleaser-linux-arm64.yml:
	GOOS=linux GOARCH=arm64 make TEMPLATE

slsa/goreleaser-darwin-amd64.yml:
	GOOS=darwin GOARCH=amd64 make TEMPLATE
slsa/goreleaser-darwin-arm64.yml:
	GOOS=darwin GOARCH=arm64 make TEMPLATE

slsa/goreleaser-windows-amd64.yml:
	GOOS=windows GOARCH=amd64 make TEMPLATE
slsa/goreleaser-windows-arm64.yml:
	GOOS=windows GOARCH=arm64 make TEMPLATE

TEMPLATE:
	cat slsa/template.yml | sed -e s/OS_HERE/${GOOS}/g | sed -e s/ARCH_HERE/${GOARCH}/g > slsa/goreleaser-${GOOS}-${GOARCH}.yml
