package spdx

import (
	"strings"
	"testing"

	"opensource.ebay.com/sbom-scorecard/pkg/scorecard"
)

func TestSpdxE2eReport(t *testing.T) {
	r := GetSpdxReport("../../examples/julia.spdx.json")
	report_text := r.Report()

	if strings.Trim(report_text, " \n") != `34 total packages
0 total files
100% have licenses.
0% have package digest.
2% have package versions.
0% have purls.
0% have CPEs.
0% have file digest.
Spec valid? true` {
		t.Log("Incorrect report results generated.\n" +
			"Got this: \n" + report_text)
		t.Fail()
	}
}

func TestSpdxE2eGrade(t *testing.T) {
	r := GetSpdxReport("../../examples/julia.spdx.json")
	report_text := scorecard.Grade(r)

	if strings.Trim(report_text, " \n") != `Spec Compliance: 25/25
Package ID: 0/20 (0% have purls and 0% have CPEs)
Package Versions: 0/20
Package Licenses: 20/20
Total points: 45/85 or 52%` {
		t.Log("Incorrect report results generated.\n" +
			"Got this: \n" + report_text)
		t.Fail()
	}
}
