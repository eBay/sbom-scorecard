package scorecard

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alexeyco/simpletable"
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

type ReportMetadata struct {
	TotalPackages int
}

type SbomReport interface {
	IsSpecCompliant() ReportValue
	PackageIdentification() ReportValue
	PackageVersions() ReportValue
	PackageLicenses() ReportValue
	CreationInfo() ReportValue
	Metadata() ReportMetadata
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
	CreationInfo          ScoreValue
	Total                 ScoreValue
	Metadata              ReportMetadata
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
		Compliance:            reportValueToScore(sr.IsSpecCompliant(), validPoints),
		PackageIdentification: reportValueToScore(sr.PackageIdentification(), packageSectionWeight),
		PackageVersions:       reportValueToScore(sr.PackageVersions(), packageSectionWeight),
		PackageLicenses:       reportValueToScore(sr.PackageLicenses(), packageSectionWeight),
		CreationInfo:          reportValueToScore(sr.CreationInfo(), generationPoints),
		Metadata: ReportMetadata{
			TotalPackages: sr.Metadata().TotalPackages,
		},
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

	maxPoints += rr.CreationInfo.MaxPoints
	totalPoints += rr.CreationInfo.Score()

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

func getReasoningText(sv ScoreValue) string {
	if sv.Reasoning != "" {
		fmt.Sprintf("(%s)", sv.Reasoning)
	}
	return ""
}
func getReportValueInfo(title string, sv ScoreValue) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s: %d/%d", title, int(sv.Score()), int(sv.MaxPoints)))
	if sv.Reasoning != "" {
		sb.WriteString(fmt.Sprintf(" (%s)", sv.Reasoning))
	}
	sb.WriteString("\n")
	return sb.String()
}
func Grade(sr SbomReport) string {
	sv := getScore(sr)
	var sb strings.Builder

	sb.WriteString(getReportValueInfo("Spec Compliance", sv.Compliance))
	sb.WriteString(getReportValueInfo("Package ID", sv.PackageIdentification))
	sb.WriteString(getReportValueInfo("Package Versions", sv.PackageVersions))
	sb.WriteString(getReportValueInfo("Package Licenses", sv.PackageLicenses))
	sb.WriteString(getReportValueInfo("Creation Info", sv.CreationInfo))

	sb.WriteString(fmt.Sprintf("Total points: %d/%d or %d%%\n", int(sv.Total.Score()), int(sv.Total.MaxPoints), PrettyPercent(int(sv.Total.Score()), int(sv.Total.MaxPoints))))

	return sb.String()
}
func GradeTableFormat(sr SbomReport) {
	sv := getScore(sr)
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Criteria"},
			{Align: simpletable.AlignCenter, Text: "Points"},
			{Align: simpletable.AlignCenter, Text: "Reasoning"},
		},
	}

	var n int

	var cells [][]*simpletable.Cell
	// ("Spec Compliance", sv.Compliance)
	n++
	cells = append(cells, *&[]*simpletable.Cell{
		{Text: fmt.Sprintf("%d", n)},
		{Text: "Spec Compliance"},
		{Text: fmt.Sprintf("%d/%d", int(sv.Compliance.Score()), int(sv.Compliance.MaxPoints))},
		{Text: red(fmt.Sprintf("%v", sv.Compliance.Reasoning))},
	})
	// ("Package ID", sv.PackageIdentification)
	n++
	cells = append(cells, *&[]*simpletable.Cell{
		{Text: fmt.Sprintf("%d", n)},
		{Text: "Package ID"},
		{Text: fmt.Sprintf("%d/%d", int(sv.PackageIdentification.Score()), int(sv.PackageIdentification.MaxPoints))},
		{Text: red(fmt.Sprintf("%v", sv.PackageIdentification.Reasoning))},
	})
	// ("Package Versions", sv.PackageVersions)
	n++
	cells = append(cells, *&[]*simpletable.Cell{
		{Text: fmt.Sprintf("%d", n)},
		{Text: "Package Versions"},
		{Text: fmt.Sprintf("%d/%d", int(sv.PackageVersions.Score()), int(sv.PackageVersions.MaxPoints))},
		{Text: red(fmt.Sprintf("%v", sv.PackageVersions.Reasoning))},
	})
	// ("Package Licenses", sv.PackageLicenses)
	n++
	cells = append(cells, *&[]*simpletable.Cell{
		{Text: fmt.Sprintf("%d", n)},
		{Text: "Package Licenses"},
		{Text: fmt.Sprintf("%d/%d", int(sv.PackageLicenses.Score()), int(sv.PackageLicenses.MaxPoints))},
		{Text: red(fmt.Sprintf("%v", sv.PackageLicenses.Reasoning))},
	})
	// ("Creation Info", sv.CreationInfo)
	n++
	cells = append(cells, *&[]*simpletable.Cell{
		{Text: fmt.Sprintf("%d", n)},
		{Text: "Creation Info"},
		{Text: fmt.Sprintf("%d/%d", int(sv.CreationInfo.Score()), int(sv.CreationInfo.MaxPoints))},
		{Text: red(fmt.Sprintf("%v", sv.CreationInfo.Reasoning))},
	})

	table.Body = &simpletable.Body{Cells: cells}
	total := fmt.Sprintf("Total points: %d/%d or %d%%\n", int(sv.Total.Score()), int(sv.Total.MaxPoints), PrettyPercent(int(sv.Total.Score()), int(sv.Total.MaxPoints)))
	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Span: 4, Text: yellow(total)},
	}}

	table.SetStyle(simpletable.StyleUnicode)
	table.Println()
}

const (
	ColorDefault = "\x1b[39m"
	ColorRed     = "\x1b[91m"
	ColorYellow  = "\x1b[93m"
	ColorGreen   = "\x1b[32m"
	ColorBlue    = "\x1b[94m"
	ColorGray    = "\x1b[90m"
)

func red(s string) string {
	return fmt.Sprintf("%s%s%s", ColorRed, s, ColorDefault)
}
func yellow(s string) string {
	return fmt.Sprintf("%s%s%s", ColorYellow, s, ColorDefault)
}
func PrettyPercent(num, denom int) int {
	if denom == 0 {
		return 0
	}
	return 100 * (1.0 * num) / denom
}
