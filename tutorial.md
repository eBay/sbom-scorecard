## Installation

Install via [git](https://git-scm.com/):

```bash
git clone https://github.com/eBay/sbom-scorecard
```

Navigate to the project directory:

```bash
cd sbom-scorecard
```

Assuming the user has Go [Go](https://go.dev/doc/install) installed,
the user can then install the package in `$GOPATH/bin` by running:

```bash
go build cmd/sbom-scorecard/main.go
```

The user can also execute the tool without using `go build` or `go install`
by using `go run`:

```bash
go run cmd/sbom-scorecard/main.go
```

## Usage:

To view helpful information about the tool and available options, run:

```bash
go run cmd/sbom-scorecard/main.go --help
```

To view helpful information about one particular option, specify that
option and then `--help`:

```bash
go run cmd/sbom-scorecard/main.go score --help
```

To run `sbom-scorecard` on an example SBOM, run:

```bash
go run cmd/sbom-scorecard/main.go score examples/julia.spdx.json
```

To run `sbom-scorecard` and specify the SBOM format type, run:

```bash
go run cmd/sbom-scorecard/main.go score --sbomtype spdx examples/julia.spdx.json
```

or

```bash
go run cmd/sbom-scorecard/main.go score --sbomtype cdx examples/dropwizard.cyclonedx.json
```

Note: `sbom-scorecard` will guess the type if no type is specified.

To run `sbom-scorecard` and specify the output format as JSON, run:

```bash
go run cmd/sbom-scorecard/main.go score --outputFormat json examples/dropwizard.cyclonedx.json
```