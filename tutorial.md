## Installation for Users

A user can download the platform-appropriate binary from
the project's [releases page](https://github.com/eBay/sbom-scorecard/releases) and save it as `sbom-scorecard`.

### Install from source

Install via git:

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

The rest of the tutorial assumes that, if you've gone this route, `$GOPATH/bin` is on your path.

## Usage:

To view helpful information about the tool and available options, run:

```bash
sbom-scorecard --help
```

To view helpful information about one particular option, specify that
option and then `--help`:

```bash
sbom-scorecard score --help
```

To run `sbom-scorecard` on an example SBOM, run:

```bash
sbom-scorecard score examples/julia.spdx.json
```

To run `sbom-scorecard` and specify the SBOM format type, run:

```bash
sbom-scorecard score --sbomtype spdx examples/julia.spdx.json
```

or

```bash
sbom-scorecard score --sbomtype cdx examples/dropwizard.cyclonedx.json
```

Note: `sbom-scorecard` will guess the type if no type is specified.

To run `sbom-scorecard` and specify the output format as JSON, run:

```bash
go run cmd/sbom-scorecard/main.go score --outputFormat json examples/dropwizard.cyclonedx.json
```
