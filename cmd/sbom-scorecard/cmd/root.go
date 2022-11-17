package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(scoreCmd)
}

var rootCmd = &cobra.Command{
	Use:   "sbom-scorecard",
	Short: "sbom-scorecard evaluates the quality of the SBOM",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
