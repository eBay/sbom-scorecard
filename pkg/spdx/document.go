// SPDX-License-Identifier: Apache-2.0

package spdx

import (
	"errors"
	"fmt"
	"os"

	spdx_json "github.com/spdx/tools-golang/json"

	"github.com/spdx/tools-golang/spdx/common"
	"github.com/spdx/tools-golang/spdx/v2_2"
	"github.com/spdx/tools-golang/spdx/v2_3"
)

const errOpenDoc = "opening SPDX %s document: %w"

var ErrUnknownFormat = fmt.Errorf("unrecognized document format")

type Document_22 v2_2.Document
type Document_23 v2_3.Document

type Package struct {
	PackageLicenseConcluded   string
	PackageExternalReferences []*PackageExternalReference
	PackageChecksums          []common.Checksum
	PackageVersion            string
}

type File struct {
	Checksums []common.Checksum
}

func LoadDocument(path string) (Document, error) {
	var err1 error
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening SPDX document: %w", err)
	}

	doc23, err1 := spdx_json.Load2_3(f)
	if err1 == nil && doc23 != nil {
		return documentFromSPDX(doc23)
	}

	// First, try to open SPDX 2.2
	doc22, err := spdx_json.Load2_2(f)
	if err != nil {
		err := fmt.Errorf(errOpenDoc, "v2.2", err)
		if err1 != nil {
			err = fmt.Errorf("%s + opening 2.3: %w ", err.Error(), err1)
		}
		return nil, err
	}

	if doc22 != nil {
		return documentFromSPDX(doc22)
	}

	return nil, ErrUnknownFormat
}

func getFiles(doc interface{}) []File {
	files := []File{}
	switch castDoc := doc.(type) {
	case *v2_2.Document:
		for _, of := range castDoc.Files {
			f := File{}
			f.Checksums = of.Checksums
			files = append(files, f)
		}
	case *v2_3.Document:
		for _, of := range castDoc.Files {
			f := File{}
			f.Checksums = of.Checksums
			files = append(files, f)
		}
	}
	return files
}

func documentFromSPDX(doc interface{}) (Document, error) {
	switch castDoc := doc.(type) {
	case *v2_2.Document:
		d := Document_22(*castDoc)
		return &d, nil
	case *v2_3.Document:
		d := Document_23(*castDoc)
		return &d, nil
	}
	return nil, errors.New("unrecognized document format")
}

func (d *Document_22) Version() string {
	return version(d)
}

func (d *Document_22) GetCreationInfo() *CreationInfo {
	return creationInfo(d)
}

func (d *Document_23) Version() string {
	return version(d)
}

func (d *Document_23) GetCreationInfo() *CreationInfo {
	return creationInfo(d)
}

func NewPackage() *Package {
	return &Package{}
}

func (p *Package) read22(sp *v2_2.Package) {
	p.PackageExternalReferences = externalReferences(sp)
	p.PackageLicenseConcluded = sp.PackageLicenseConcluded
	p.PackageChecksums = sp.PackageChecksums
	p.PackageVersion = sp.PackageVersion
}

func (p *Package) read23(sp *v2_3.Package) {
	p.PackageExternalReferences = externalReferences(sp)
	p.PackageLicenseConcluded = sp.PackageLicenseConcluded
	p.PackageChecksums = sp.PackageChecksums
	p.PackageVersion = sp.PackageVersion
}

func (d *Document_22) GetFiles() []File {
	return getFiles(d)
}

func (d *Document_23) GetFiles() []File {
	return getFiles(d)
}

func (d *Document_23) GetPackages() []Package {
	packages := []Package{}
	for _, p := range d.Packages {
		np := Package{}
		np.read23(p)
		packages = append(packages, np)
	}
	return packages
}

func (d *Document_22) GetPackages() []Package {
	packages := []Package{}
	for _, p := range d.Packages {
		np := Package{}
		np.read22(p)
		packages = append(packages, np)
	}
	return packages
}

type CreationInfo struct {
	LicenseListVersion string
	Creators           []common.Creator
	Created            string
	CreatorComment     string
}

type PackageExternalReference struct {
	Category           string
	RefType            string
	Locator            string
	ExternalRefComment string
}

type Document interface {
	Version() string
	GetCreationInfo() *CreationInfo
	GetPackages() []Package
	GetFiles() []File
}

func externalReferences(pkg interface{}) []*PackageExternalReference {
	refs := []*PackageExternalReference{}
	switch castPkg := pkg.(type) {
	case *v2_2.Package:
		if castPkg.PackageExternalReferences == nil {
			return nil
		}

		for _, p := range castPkg.PackageExternalReferences {
			refs = append(refs, &PackageExternalReference{
				Category:           p.Category,
				RefType:            p.RefType,
				Locator:            p.Locator,
				ExternalRefComment: p.ExternalRefComment,
			})
		}
	case *v2_3.Package:
		if castPkg.PackageExternalReferences == nil {
			return nil
		}

		for _, p := range castPkg.PackageExternalReferences {
			refs = append(refs, &PackageExternalReference{
				Category:           p.Category,
				RefType:            p.RefType,
				Locator:            p.Locator,
				ExternalRefComment: p.ExternalRefComment,
			})
		}
	}

	return refs
}

func creationInfo(doc Document) *CreationInfo {
	var ci *CreationInfo
	creators := []common.Creator{}

	switch castDoc := doc.(type) {
	case *Document_22:
		if castDoc.CreationInfo == nil {
			return nil
		}
		ci = &CreationInfo{
			LicenseListVersion: castDoc.CreationInfo.LicenseListVersion,
			Created:            castDoc.CreationInfo.Created,
			CreatorComment:     castDoc.CreationInfo.CreatorComment,
		}
		if castDoc.CreationInfo.Creators != nil {
			creators = castDoc.CreationInfo.Creators
		}
	case *Document_23:
		if castDoc.CreationInfo == nil {
			return nil
		}
		ci = &CreationInfo{
			LicenseListVersion: castDoc.CreationInfo.LicenseListVersion,
			Created:            castDoc.CreationInfo.Created,
			CreatorComment:     castDoc.CreationInfo.CreatorComment,
		}
		if castDoc.CreationInfo.Creators != nil {
			creators = castDoc.CreationInfo.Creators
		}
	}
	ci.Creators = creators
	return ci
}

func version(doc Document) string {
	switch doc.(type) {
	case *Document_22:
		return "2.2"
	case *Document_23:
		return "2.3"
	}
	return ""
}
