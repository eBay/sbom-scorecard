package spdx

import (
	"strings"
	"testing"

	"github.com/ebay/sbom-scorecard/pkg/scorecard"
)

var report_tests = []struct {
	path     string
	expected string
}{
	{"../../examples/julia.spdx.json", `34 total packages
0 total files
100% have licenses.
0% have package digest.
2% have package versions.
0% have purls.
0% have CPEs.
0% have file digest.
Spec valid? true
Has creation info? false`},
	{"../../examples/photon.spdx.json", `38 total packages
0 total files
94% have licenses.
2% have package digest.
97% have package versions.
0% have purls.
0% have CPEs.
0% have file digest.
Spec valid? true
Has creation info? true`},
}

func TestSpdxE2eReport(t *testing.T) {
	for _, e := range report_tests {
		res := GetSpdxReport(e.path)
		report_text := res.Report()
		if strings.Trim(report_text, " \n") != e.expected {
			t.Errorf("GetSpdxReport(%v) = %v, expected %v",
				e.path, strings.Trim(report_text, " \n"), e.expected)
		}
	}
}

var grade_tests = []struct {
	path     string
	expected string
}{
	{"../../examples/julia.spdx.json", `Spec Compliance: 25/25
Package ID: 0/20 (0% have purls and 0% have CPEs)
Package Versions: 0/20
Package Licenses: 20/20
Creation Info: 0/15 (No tool was used to create the sbom)
Total points: 45/100 or 45%`},
	{"../../examples/photon.spdx.json", `Spec Compliance: 25/25
Package ID: 0/20 (0% have purls and 0% have CPEs)
Package Versions: 19/20
Package Licenses: 18/20
Creation Info: 15/15
Total points: 78/100 or 78%`},
}

func TestSpdxE2eGrade(t *testing.T) {
	for _, e := range grade_tests {
		res := GetSpdxReport(e.path)
		report_text := scorecard.Grade(res)
		if strings.Trim(report_text, " \n") != e.expected {
			t.Errorf("GetSpdxReport(%v) = %v, expected %v",
				e.path, strings.Trim(report_text, " \n"), e.expected)
		}
	}
}
