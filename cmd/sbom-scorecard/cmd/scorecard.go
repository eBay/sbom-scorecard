package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"opensource.ebay.com/sbom-scorecard/pkg/cdx"
	"opensource.ebay.com/sbom-scorecard/pkg/scorecard"
	"opensource.ebay.com/sbom-scorecard/pkg/spdx"
)

var flags = struct {
	sbomType string
}{}

type options struct {
	sbomType string
	// path to file being evaluated
	path string
}

func init() {
	scoreCmd.PersistentFlags().StringVar(&flags.sbomType, "sbomtype", "spdx", "type of sbom being evaluated")
	_ = scoreCmd.MarkPersistentFlagRequired("sbomType")
}

var scoreCmd = &cobra.Command{
	Use:   "score",
	Short: "score evaluates the SBOM being passed and outputs based on the composition and completion",
	Run: func(cmd *cobra.Command, args []string) {

		opts, err := validateFlags(args)
		if err != nil {
			fmt.Printf("unable to validate flags: %v\n", err)
			os.Exit(1)
		}

		var r scorecard.SbomReport
		switch opts.sbomType {
		case "spdx":
			fmt.Printf("Reading SPDX SBOM %s\n", opts.path)
			r = spdx.GetSpdxReport(opts.path)
		case "cdx":
			fmt.Printf("Reading CDX SBOM %s\n", opts.path)
			r = cdx.GetCycloneDXReport(opts.path)
		}

		fmt.Print(r.Report())
	},
}

func validateFlags(args []string) (options, error) {
	var opts options
	opts.sbomType = flags.sbomType
	if len(args) != 1 {
		return opts, fmt.Errorf("expected positional argument for file_path")
	}
	opts.path = args[0]
	return opts, nil
}
