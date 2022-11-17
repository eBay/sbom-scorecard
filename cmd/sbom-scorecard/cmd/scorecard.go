//
// Copyright 2022 The GUAC Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	exampleCmd.PersistentFlags().StringVar(&flags.sbomType, "sbomtype", "spdx", "type of sbom being evaluated")
	_ = exampleCmd.MarkPersistentFlagRequired("sbomType")
}

var exampleCmd = &cobra.Command{
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
