package scorecard

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ReportValue struct {
	Ratio     float32
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

type ScoreValue struct {
	ReportValue
	MaxPoints float32
}

func (sv *ScoreValue) Score() float32 {
	return sv.Ratio * sv.MaxPoints
}

type ReportResult struct {
	Compliance            ScoreValue
	PackageIdentification ScoreValue
	PackageVersions       ScoreValue
	PackageLicenses       ScoreValue
	Total                 ScoreValue
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

func reportValueToScore(rv ReportValue, maxPoints float32) ScoreValue {
	sv := ScoreValue{
		MaxPoints: maxPoints,
	}
	sv.Ratio = rv.Ratio
	sv.Reasoning = rv.Reasoning
	return sv
}

func getScore(sr SbomReport) ReportResult {
	rr := ReportResult{
		Compliance:            reportValueToScore(sr.IsSpecCompliant(), 25),
		PackageIdentification: reportValueToScore(sr.PackageIdentification(), 20),
		PackageVersions:       reportValueToScore(sr.PackageVersions(), 20),
		PackageLicenses:       reportValueToScore(sr.PackageLicenses(), 20),
	}
	var totalPoints float32
	var maxPoints float32

	maxPoints += rr.Compliance.MaxPoints
	totalPoints += rr.Compliance.Score()

	maxPoints += rr.PackageIdentification.MaxPoints
	totalPoints += rr.PackageIdentification.Score()

	maxPoints += rr.PackageVersions.MaxPoints
	totalPoints += rr.PackageVersions.Score()

	maxPoints += rr.PackageLicenses.MaxPoints
	totalPoints += rr.PackageLicenses.Score()

	rr.Total = ScoreValue{
		MaxPoints: maxPoints,
	}
	rr.Total.Ratio = totalPoints / maxPoints

	return rr
}

func JsonGrade(sr SbomReport) string {
	out, _ := json.Marshal(getScore(sr))
	return string(out)
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
