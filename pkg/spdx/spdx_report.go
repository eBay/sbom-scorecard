package spdx

import (
	"fmt"
	"strings"

	"github.com/ebay/sbom-scorecard/pkg/scorecard"
	spdx_common "github.com/spdx/tools-golang/spdx/v2/common"

	"regexp"
)

var isNumeric = regexp.MustCompile(`\d`)

var missingPackages = scorecard.ReportValue{
	Ratio:     0,
	Reasoning: "No packages",
}

type SpdxReport struct {
	doc      Document
	docError error
	valid    bool

	totalPackages int
	totalFiles    int
	hasLicense    int
	hasPackDigest int
	hasPurl       int
	hasCPE        int
	hasPurlOrCPE  int
	hasFileDigest int
	hasPackVer    int
}

func (r *SpdxReport) Metadata() scorecard.ReportMetadata {
	return scorecard.ReportMetadata{
		TotalPackages: r.totalPackages,
	}
}

func (r *SpdxReport) Report() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d total packages\n", r.totalPackages))
	sb.WriteString(fmt.Sprintf("%d total files\n", r.totalFiles))
	sb.WriteString(fmt.Sprintf("%d%% have licenses.\n", scorecard.PrettyPercent(r.hasLicense, r.totalPackages)))
	sb.WriteString(fmt.Sprintf("%d%% have package digest.\n", scorecard.PrettyPercent(r.hasPackDigest, r.totalPackages)))
	sb.WriteString(fmt.Sprintf("%d%% have package versions.\n", scorecard.PrettyPercent(r.hasPackVer, r.totalPackages)))
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
	either := scorecard.PrettyPercent(r.hasPurlOrCPE, r.totalPackages)
	return scorecard.ReportValue{
		// What percentage has both Purl & CPEs?
		Ratio:     float32(r.hasPurlOrCPE) / float32(r.totalPackages),
		Reasoning: fmt.Sprintf("%d%% have either purls (%d%%) or CPEs (%d%%)", either, purlPercent, cpePercent),
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
		Ratio: float32(r.hasPackVer) / float32(r.totalPackages),
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
		Ratio: float32(r.hasLicense) / float32(r.totalPackages),
	}
}

func (r *SpdxReport) CreationInfo() scorecard.ReportValue {
	foundTool := false
	hasVersion := false

	if r.doc == nil || r.doc.GetCreationInfo() == nil {
		return scorecard.ReportValue{
			Ratio:     0,
			Reasoning: "No creation info found",
		}
	}

	for _, creator := range r.doc.GetCreationInfo().Creators {
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

	var score float32
	score = 1.0
	reasons := []string{}

	if !hasVersion {
		score -= .2
		reasons = append(reasons, "The tool used to create the sbom does not have a version")
	}

	if r.doc.GetCreationInfo().Created == "" {
		score -= .2
		reasons = append(reasons, "There is no timestamp for when the sbom was created")
	}

	return scorecard.ReportValue{
		Ratio:     score,
		Reasoning: strings.Join(reasons, ", "),
	}

}

func GetSpdxReport(filename string) scorecard.SbomReport {
	sr := SpdxReport{}
	doc, err := LoadDocument(filename)
	if err != nil {
		fmt.Printf("loading document: %v\n", err)
		return &sr
	}

	// try to load the SPDX file's contents as a json file, version 2.2
	sr.doc = doc
	sr.docError = err
	sr.valid = err == nil
	if sr.doc != nil {
		packages := sr.doc.GetPackages()

		for _, p := range packages {
			sr.totalPackages += 1
			if p.PackageLicenseConcluded != "NONE" &&
				p.PackageLicenseConcluded != "NOASSERTION" &&
				p.PackageLicenseConcluded != "" {
				sr.hasLicense += 1
			} else if p.PackageLicenseDeclared != "NONE" &&
				p.PackageLicenseDeclared != "NOASSERTION" &&
				p.PackageLicenseDeclared != "" {
				sr.hasLicense += 1
			}

			if len(p.PackageChecksums) > 0 {
				sr.hasPackDigest += 1
			}

			var foundCPE bool
			var foundPURL bool
			for _, ref := range p.PackageExternalReferences {
				if !foundPURL && ref.RefType == spdx_common.TypePackageManagerPURL {
					sr.hasPurl += 1
					foundPURL = true
				}

				if !foundCPE && strings.HasPrefix(ref.RefType, "cpe") {
					sr.hasCPE += 1
					foundCPE = true
				}
			}
			if foundCPE && foundPURL {
				sr.hasPurlOrCPE += 1
			}

			if p.PackageVersion != "" {
				sr.hasPackVer += 1
			}
		}

		for _, file := range sr.doc.GetFiles() {
			sr.totalFiles += 1
			if len(file.Checksums) > 0 {
				sr.hasFileDigest += 1
			}
		}
	}
	return &sr
}
