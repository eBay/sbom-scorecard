package main

import (
	"fmt"
	"os"
)

type SbomReport interface {
	Report() string
}

func main() {
	filename := os.Args[1]
	fmt.Printf("Reading %s\n", filename)

	// TODO: Logic to dispatch across file names & types
	r := GetSpdxReport(filename)

	fmt.Printf(r.Report())
}

func prettyPercent(num, denom int) int {
	if (denom == 0) {
		return 0
	}
	return 100 * (1.0 * num) / denom
}
