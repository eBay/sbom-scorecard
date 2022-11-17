package scorecard

import (
	"strings"
	"fmt"
)

type ReportValue struct {
	Ratio float32
	Reasoning string
}

const validPoints = 25
const generationPoints = 15
const packageSectionWeight = 20
// TODO capture generation points
const maxPoints = validPoints + packageSectionWeight*3

type SbomReport interface {
	IsSpecCompliant() ReportValue
	PackageIdentification() ReportValue
	PackageVersions() ReportValue
	PackageLicenses() ReportValue
	Report() string
}

func getReportValueInfo(title string, rv ReportValue, maxPoints float32) (string, int) {
	sc := int(rv.Ratio * maxPoints)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s: %d/%d", title, sc, int(maxPoints)))
	if rv.Reasoning != "" {
		sb.WriteString(fmt.Sprintf(" (%s)", rv.Reasoning))
	}
	sb.WriteString("\n")
	return sb.String(), sc
}

func Grade(sr SbomReport) string {
	points := 0
	var sb strings.Builder

	output, score := getReportValueInfo("Spec Compliance", sr.IsSpecCompliant(), validPoints)
	points += score
	sb.WriteString(output)

	output, score = getReportValueInfo("Package ID", sr.PackageIdentification(), packageSectionWeight)
	points += score
	sb.WriteString(output)

	output, score = getReportValueInfo("Package Versions", sr.PackageVersions(), packageSectionWeight)
	points += score
	sb.WriteString(output)

	output, score = getReportValueInfo("Package Licenses", sr.PackageLicenses(), packageSectionWeight)
	points += score
	sb.WriteString(output)

	sb.WriteString(fmt.Sprintf("Total points: %d/%d or %d%%\n", points, maxPoints, PrettyPercent(points, maxPoints)))
	
	return sb.String()
}

func PrettyPercent(num, denom int) int {
	if denom == 0 {
		return 0
	}
	return 100 * (1.0 * num) / denom
}
