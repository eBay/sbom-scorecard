package main

import (
	"os"
	"fmt"
	cdx "github.com/CycloneDX/cyclonedx-go"
	"strings"
)


type CycloneDXReport struct {
	valid bool
	totalPackages int
	hasLicense int
	hasPackDigest int
	hasPurl int
	hasCPE int
}

func (r *CycloneDXReport) Report() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d total packages\n", r.totalPackages))
	sb.WriteString(fmt.Sprintf("%d%% have licenses.\n", prettyPercent(r.hasLicense, r.totalPackages)))
	sb.WriteString(fmt.Sprintf("%d%% have package digest.\n", prettyPercent(r.hasPackDigest, r.totalPackages)))
	sb.WriteString(fmt.Sprintf("%d%% have purls.\n", prettyPercent(r.hasPurl, r.totalPackages)))
	sb.WriteString(fmt.Sprintf("%d%% have CPEs.\n", prettyPercent(r.hasCPE, r.totalPackages)))
	sb.WriteString(fmt.Sprintf("Spec valid? %v\n", r.valid))
	return sb.String()

}
 
func GetCycloneDXReport(filename string) SbomReport {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error while opening %v for reading: %v", filename, err)
		return nil
	}
	defer f.Close()

	r := CycloneDXReport{}

	bom := new(cdx.BOM)
	decoder := cdx.NewBOMDecoder(f, cdx.BOMFileFormatJSON)
	if err = decoder.Decode(bom); err != nil {
		r.valid = false
		return &r;
	}
	r.valid = true

	if bom.Components != nil {
		for _, p := range *bom.Components {
			r.totalPackages += 1
			if len(*p.Licenses) > 0 {
				r.hasLicense += 1
			}
			if len(*p.Hashes) > 0 {
				r.hasPackDigest += 1
			}
			if p.PackageURL != "" {
				r.hasPurl += 1
			}
			if p.CPE != "" {
				r.hasCPE += 1
			}
		}
	}

	return &r
}
