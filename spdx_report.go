package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	spdx_json "github.com/spdx/tools-golang/json"
	spdx_common "github.com/spdx/tools-golang/spdx/common"
)

type SpdxReport struct {
	valid         bool
	totalPackages int
	totalFiles    int
	hasLicense    int
	hasPackDigest int
	hasPurl       int
	hasCPE        int
	hasFileDigest int
}

func (r *SpdxReport) Report() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d total packages\n", r.totalPackages))
	sb.WriteString(fmt.Sprintf("%d total files\n", r.totalFiles))
	sb.WriteString(fmt.Sprintf("%d%%%% have licenses.\n", prettyPercent(r.hasLicense, r.totalPackages)))
	sb.WriteString(fmt.Sprintf("%d%%%% have package digest.\n", prettyPercent(r.hasPackDigest, r.totalPackages)))
	sb.WriteString(fmt.Sprintf("%d%%%% have purls.\n", prettyPercent(r.hasPurl, r.totalPackages)))
	sb.WriteString(fmt.Sprintf("%d%%%% have CPEs.\n", prettyPercent(r.hasCPE, r.totalPackages)))
	sb.WriteString(fmt.Sprintf("%d%%%% have file digest.\n", prettyPercent(r.hasFileDigest, r.totalFiles)))
	sb.WriteString(fmt.Sprintf("Spec valid? %v\n", r.valid))
	return sb.String()
}

func GetSpdxReport(filename string) SbomReport {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error while opening %v for reading: %v", filename, err)
		return nil
	}
	defer f.Close()
	byteValue, _ := ioutil.ReadAll(f)
	sr := SpdxReport{}

	// // try to load the SPDX file's contents as a json file, version 2.2
	spdxDoc, err := spdx_json.Load2_2(bytes.NewReader(byteValue))
	sr.valid = err == nil
	if spdxDoc != nil {
		if len(spdxDoc.Packages) > 0 {
			for _, p := range spdxDoc.Packages {
				sr.totalPackages += 1
				if p.PackageLicenseConcluded != "NONE" &&
					p.PackageLicenseConcluded != "NOASSERTION" &&
					p.PackageLicenseConcluded != "" {
					sr.hasLicense += 1
				}

				if len(p.PackageChecksums) > 0 {
					sr.hasPackDigest += 1
				}

				for _, ref := range p.PackageExternalReferences {
					if ref.RefType == spdx_common.TypePackageManagerPURL {
						sr.hasPurl += 1
					}
				}

				for _, ref := range p.PackageExternalReferences {
					if strings.HasPrefix(ref.RefType, "cpe") {
						sr.hasCPE += 1
						break
					}
				}
			}
		}

		for _, file := range spdxDoc.Files {
			sr.totalFiles += 1
			if len(file.Checksums) > 0 {
				sr.hasFileDigest += 1
			}
		}
	}

	return &sr
}
