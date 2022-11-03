package main

import (
	"fmt"
	"bytes"
	"strings"
	"io/ioutil"
	"encoding/json"
	"github.com/spdx/tools-golang/jsonloader"
	"os"
)

type CreationInfo struct {
	Created string `json:"created"`
	Creators []string
}

type ExternalRef struct {
	ReferenceCategory string
	ReferenceLocator string
	ReferenceType string
}

type Package struct {
	Name string
	LicenseConcluded string
	ExternalRefs []ExternalRef
}

type SpdxDocument struct {
	Name string `json:name`
	SpdxVersion string `json:"spdxVersion"`
	License string `json:"dataLicense"`
	CreationInfo CreationInfo `json:"creationInfo"`
	Packages []Package
}

type SpdxReport struct {
	valid bool
	totalPackages int
	hasLicense int
	hasPurl int
	hasCPE int
}

func (r *SpdxReport) Report() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d total packages\n", r.totalPackages))
	sb.WriteString(fmt.Sprintf("%d%%%% have licenses.\n", prettyPercent(r.hasLicense, r.totalPackages)))
	sb.WriteString(fmt.Sprintf("%d%%%% have purls.\n", prettyPercent(r.hasPurl, r.totalPackages)))
	sb.WriteString(fmt.Sprintf("%d%%%% have CPEs.\n", prettyPercent(r.hasCPE, r.totalPackages)))
	sb.WriteString(fmt.Sprintf("Spec valid? %v\n", r.valid))
	return sb.String()	
}


func GetSpdxReport(filename string) SbomReport {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error while opening %v for reading: %v", filename, err)
		return nil
	}
	defer f.Close()
	byteValue, _ := ioutil.ReadAll(f)
	sr := SpdxReport{}
	
	// // try to load the SPDX file's contents as a json file, version 2.2
	_, err = jsonloader.Load2_2(bytes.NewReader(byteValue))
	sr.valid = err == nil

	var doc2 SpdxDocument
	err = json.Unmarshal(byteValue, &doc2)
	if err != nil {
		panic(err)
	}



	for _, p := range doc2.Packages {
		sr.totalPackages += 1
		if p.LicenseConcluded != "NONE" &&
			p.LicenseConcluded != "NOASSERTION" &&
			p.LicenseConcluded != "" {
			sr.hasLicense += 1
		}

		for _, ref := range p.ExternalRefs {
			if ref.ReferenceType == "purl" {
				sr.hasPurl += 1
				break
			}
		}

		for _, ref := range p.ExternalRefs {
			if ref.ReferenceType == "cpe23Type" {
				sr.hasCPE += 1
				break
			}
		}
	}

	return &sr
}
