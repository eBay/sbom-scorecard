package spdx

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	spdx_json "github.com/spdx/tools-golang/json"
	spdx_common "github.com/spdx/tools-golang/spdx/common"
	"opensource.ebay.com/sbom-scorecard/pkg/scorecard"

	"bytes"

	"regexp"

	"github.com/spdx/tools-golang/spdx/v2_2"
)

var isNumeric = regexp.MustCompile(`\d`)

var missingPackages = scorecard.ReportValue{
	Ratio:     0,
	Reasoning: "No packages",
}

type SpdxReport struct {
	doc      *v2_2.Document
	docError error
	valid    bool

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
	sb.WriteString(fmt.Sprintf("%d%% have licenses.\n", scorecard.PrettyPercent(r.hasLicense, r.totalPackages)))
	sb.WriteString(fmt.Sprintf("%d%% have package digest.\n", scorecard.PrettyPercent(r.hasPackDigest, r.totalPackages)))
	sb.WriteString(fmt.Sprintf("%d%% have purls.\n", scorecard.PrettyPercent(r.hasPurl, r.totalPackages)))
	sb.WriteString(fmt.Sprintf("%d%% have CPEs.\n", scorecard.PrettyPercent(r.hasCPE, r.totalPackages)))
	sb.WriteString(fmt.Sprintf("%d%% have file digest.\n", scorecard.PrettyPercent(r.hasFileDigest, r.totalFiles)))
	sb.WriteString(fmt.Sprintf("Spec valid? %v\n", r.valid))
	sb.WriteString(fmt.Sprintf("Has creation info? %v\n", r.CreationInfo().Ratio == 1))

	return sb.String()
}

func (r *SpdxReport) IsSpecCompliant() scorecard.ReportValue {
	if r.docError != nil {
		return scorecard.ReportValue{
			Ratio:     0,
			Reasoning: r.docError.Error(),
		}
	}
	return scorecard.ReportValue{Ratio: 1}
}

func (r *SpdxReport) PackageIdentification() scorecard.ReportValue {
	if r.totalPackages == 0 {
		return missingPackages
	}
	purlPercent := scorecard.PrettyPercent(r.hasPurl, r.totalPackages)
	cpePercent := scorecard.PrettyPercent(r.hasCPE, r.totalPackages)
	return scorecard.ReportValue{
		// What percentage has both Purl & CPEs?
		Ratio:     float32(r.hasPurl+r.hasCPE) / float32(r.totalPackages*2),
		Reasoning: fmt.Sprintf("%d%% have purls and %d%% have CPEs", purlPercent, cpePercent),
	}
}

func (r *SpdxReport) PackageVersions() scorecard.ReportValue {
	if r.totalPackages == 0 {
		return scorecard.ReportValue{
			Ratio:     0,
			Reasoning: "No packages",
		}
	}
	return scorecard.ReportValue{
		// TODO
	}
}

func (r *SpdxReport) PackageLicenses() scorecard.ReportValue {
	if r.totalPackages == 0 {
		return scorecard.ReportValue{
			Ratio:     0,
			Reasoning: "No packages",
		}
	}
	return scorecard.ReportValue{
		Ratio: float32(r.hasLicense / r.totalPackages),
	}
}

func (r *SpdxReport) CreationInfo() scorecard.ReportValue {
	foundTool := false
	hasVersion := false

	if r.doc == nil || r.doc.CreationInfo == nil {
		return scorecard.ReportValue{
			Ratio:     0,
			Reasoning: "No creation info found",
		}
	}

	for _, creator := range r.doc.CreationInfo.Creators {
		if creator.CreatorType == "Tool" {
			foundTool = true
			if isNumeric.MatchString(creator.Creator) {
				hasVersion = true
			}
		}
	}

	if !foundTool {
		return scorecard.ReportValue{
			Ratio:     0,
			Reasoning: "No tool was used to create the sbom",
		}
	}

	if !hasVersion {
		return scorecard.ReportValue{
			Ratio:     .2,
			Reasoning: "The tool used to create the sbom does not have a version",
		}
	}

	return scorecard.ReportValue{
		Ratio: 1,
	}

}

func GetSpdxReport(filename string) scorecard.SbomReport {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error while opening %v for reading: %v", filename, err)
		return nil
	}
	defer f.Close()
	byteValue, _ := ioutil.ReadAll(f)
	sr := SpdxReport{}

	// try to load the SPDX file's contents as a json file, version 2.2
	spdxDoc, err := spdx_json.Load2_2(bytes.NewReader(byteValue))
	sr.doc = spdxDoc
	sr.docError = err
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
