package cdx

import (
	"strings"
	"testing"
)

func TestCycloneE2e(t *testing.T) {
	r := GetCycloneDXReport("../../examples/dropwizard.cyclonedx.json")

	report_text := r.Report()

	if strings.Trim(report_text, " \n") != `167 total packages
79% have licenses.
100% have package digest.
100% have purls.
0% have CPEs.
Spec valid? true` {
		t.Log("Incorrect report text generated.\n" +
			"Got this:\n" + report_text)
		t.Fail()
	}
}
