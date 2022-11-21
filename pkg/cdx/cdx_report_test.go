package cdx

import (
	"strings"
	"testing"

	"opensource.ebay.com/sbom-scorecard/pkg/scorecard"
)

func TestCycloneE2eReport(t *testing.T) {
	r := GetCycloneDXReport("../../examples/dropwizard.cyclonedx.json")

	report_text := r.Report()

	if strings.Trim(report_text, " \n") != `167 total packages
79% have licenses.
100% have package digest.
100% have purls.
0% have CPEs.
Has creation info? true
Spec valid? true` {
		t.Log("Incorrect report text generated.\n" +
			"Got this:\n" + report_text)
		t.Fail()
	}
}

func TestCycloneE2eGrade(t *testing.T) {
	r := GetCycloneDXReport("../../examples/dropwizard.cyclonedx.json")

	report_text := scorecard.Grade(r)

	if strings.Trim(report_text, " \n") != `Spec Compliance: 25/25
Package ID: 10/20 (100% have purls and 0% have CPEs)
Package Versions: 20/20
Package Licenses: 15/20
Creation Info: 15/15
Total points: 85/100 or 85%` {
		t.Log("Incorrect report text generated.\n" +
			"Got this:\n" + report_text)
		t.Fail()
	}
}
