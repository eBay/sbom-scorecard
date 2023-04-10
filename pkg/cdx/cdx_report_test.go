package cdx

import (
	"strings"
	"testing"

	"github.com/ebay/sbom-scorecard/pkg/scorecard"
)

func assertTextEqual(t *testing.T, actual string, expected string) {
	if strings.Trim(actual, " \n") != expected {
		t.Log("Incorrect report text generated.\n" +
			"Got this:\n" + actual + "\n\n but expected:\n" + expected)
		t.Fail()
	}
}

func TestCycloneE2eReport(t *testing.T) {
	r := GetCycloneDXReport("../../examples/dropwizard.cyclonedx.json")

	report_text := r.Report()
	assertTextEqual(t,
		report_text,
		`167 total packages
100% have versions.
79% have licenses.
100% have package digest.
100% have purls.
0% have CPEs.
Has creation info? true
Spec valid? true`)
}

func TestCycloneE2eGrade(t *testing.T) {
	r := GetCycloneDXReport("../../examples/dropwizard.cyclonedx.json")

	report_text := scorecard.Grade(r)
	assertTextEqual(t,
		report_text,
		`Spec Compliance: 25/25
Package ID: 20/20 (100% have either a purl (100%) or CPE (0%))
Package Versions: 20/20
Package Licenses: 15/20
Creation Info: 15/15
Total points: 95/100 or 95%`)
}

func TestCycloneXML(t *testing.T) {
	r := GetCycloneDXReport("../../examples/openfeature-javasdk.cyclonedx.xml")

	report_text := scorecard.Grade(r)
	assertTextEqual(t,
		report_text,
		`Spec Compliance: 25/25
Package ID: 20/20 (100% have either a purl (100%) or CPE (0%))
Package Versions: 20/20
Package Licenses: 18/20
Creation Info: 15/15
Total points: 98/100 or 98%`)
}

func TestCycloneInvalid(t *testing.T) {
	r := GetCycloneDXReport("../../examples/invalid.json")

	report_text := scorecard.Grade(r)
	assertTextEqual(t,
		report_text,
		`Spec Compliance: 0/25 (Couldn't parse the SBOM)
Package ID: 0/20 (No packages)
Package Versions: 0/20
Package Licenses: 0/20
Creation Info: 0/15 (Missing creation info)
Total points: 0/100 or 0%`)
}
