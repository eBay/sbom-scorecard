package cmd

import (
	"fmt"
	"os"

	"errors"

	"strings"

	"github.com/ebay/sbom-scorecard/pkg/cdx"
	"github.com/ebay/sbom-scorecard/pkg/scorecard"
	"github.com/ebay/sbom-scorecard/pkg/spdx"
	"github.com/spf13/cobra"
)

var flags = struct {
	sbomType     string
	outputFormat string
}{}

type options struct {
	sbomType string
	// path to file being evaluated
	path         string
	outputFormat string
}

func init() {
	scoreCmd.PersistentFlags().StringVar(&flags.sbomType, "sbomtype", "guess", "type of sbom being evaluated")
	scoreCmd.PersistentFlags().StringVar(&flags.outputFormat, "outputFormat", "text", "format to output")
	_ = scoreCmd.MarkPersistentFlagRequired("sbomType")
}

func determineSbomType(filepath string) string {
	content, err := os.ReadFile(filepath)
	if err != nil {
		panic(fmt.Sprintf("Error! %v", err))
	}

	if strings.Contains(strings.ToLower(string(content)), "spdx") {
		return "spdx"
	}
	return "cdx"
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

		if opts.sbomType == "guess" {
			opts.sbomType = determineSbomType(opts.path)
			print("Guessed: " + opts.sbomType + "\n")
		}

		switch opts.sbomType {
		case "spdx":
			r = spdx.GetSpdxReport(opts.path)
		case "cdx":
			r = cdx.GetCycloneDXReport(opts.path)
		}

		if opts.outputFormat == "json" {
			fmt.Println(scorecard.JsonGrade(r))
		} else {
			print(r.Report())
			print("==\n")
			print(scorecard.Grade(r))
		}
	},
}

func validateFlags(args []string) (options, error) {
	var opts options
	opts.sbomType = flags.sbomType
	if flags.sbomType != "spdx" &&
		flags.sbomType != "cdx" &&
		flags.sbomType != "guess" {
		return opts, errors.New(fmt.Sprintf("Unknown sbomType %s", flags.sbomType))
	}

	opts.outputFormat = flags.outputFormat
	if flags.outputFormat != "text" && flags.outputFormat != "json" {
		return opts, errors.New(fmt.Sprintf("Unknown outputFormat %s", flags.outputFormat))
	}

	if len(args) != 1 {
		return opts, fmt.Errorf("expected positional argument for file_path")
	}
	opts.path = args[0]
	return opts, nil
}
