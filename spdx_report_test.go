package main

import (
	"strings"
	"testing"
)

func TestSpdxE2e(t *testing.T) {
	r := GetSpdxReport("./examples/julia.spdx.json")
	report_text := r.Report()

	if strings.Trim(report_text, " \n") != `34 total packages
0 total files
100% have licenses.
0% have package digest.
0% have purls.
0% have CPEs.
0% have file digest.
Spec valid? true` {
		t.Log("Incorrect report results generated.\n" +
			"Got this: \n" + report_text)
		t.Fail()
	}
}
