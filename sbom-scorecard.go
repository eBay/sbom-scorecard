package main

import (
	"fmt"
	"os"
	"strings"
)

type SbomReport interface {
	Report() string
}

func main() {
	filename := os.Args[1]
	fmt.Printf("Reading %s\n", filename)
	var r SbomReport

	if strings.Contains(filename, "spdx") {
		r = GetSpdxReport(filename)
	} else if strings.Contains(filename, "cyclonedx") {
		r = GetCycloneDXReport(filename)
	}

	fmt.Printf(r.Report())
}

func prettyPercent(num, denom int) int {
	if (denom == 0) {
		return 0
	}
	return 100 * (1.0 * num) / denom
}
