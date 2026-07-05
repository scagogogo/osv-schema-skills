package osv_schema

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// ------------------------------------------------ ---------------------------------------------------------------------

type References []*Reference

var _ sql.Scanner = &References{}
var _ driver.Valuer = &References{}

func (x References) FilterByType(referenceTypes ...ReferenceType) References {

	if len(referenceTypes) == 0 {
		return nil
	}

	referenceTypeSet := make(map[ReferenceType]struct{}, 0)
	for _, r := range referenceTypes {
		referenceTypeSet[r] = struct{}{}
	}

	slice := make([]*Reference, 0)
	for _, r := range x {
		if _, exists := referenceTypeSet[r.Type]; exists {
			slice = append(slice, r)
		}
	}
	return slice
}

func (x *References) Scan(src any) error {
	if src == nil {
		return nil
	}
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("scan error")
	}
	if len(bytes) == 0 {
		return nil
	}
	return json.Unmarshal(bytes, &x)
}

func (x References) Value() (driver.Value, error) {
	if len(x) == 0 {
		return nil, nil
	}
	// References only contains string-typed fields (Type, URL), so json.Marshal
	// cannot fail; the error return is therefore unreachable.
	marshal, _ := json.Marshal(x)
	return string(marshal), nil
}

type ReferenceType string

const (

	// ReferenceTypeAdvisory A published security advisory for the vulnerability.
	ReferenceTypeAdvisory ReferenceType = "ADVISORY"

	// ReferenceTypeArticle An article or blog post describing the vulnerability.
	ReferenceTypeArticle ReferenceType = "ARTICLE"

	// ReferenceTypeDetection A tool, script, scanner, or other mechanism that allows for detection of the vulnerability
	// in production environments. e.g. YARA rules, hashes, virus signature, or other scanners.
	ReferenceTypeDetection ReferenceType = "DETECTION"

	// ReferenceTypeDiscussion A social media discussion regarding the vulnerability, e.g. a Twitter, Mastodon, Hacker News,
	// or Reddit thread.
	ReferenceTypeDiscussion ReferenceType = "DISCUSSION"

	// ReferenceTypeReport A report, typically on a bug or issue tracker, of the vulnerability.
	ReferenceTypeReport ReferenceType = "REPORT"

	// ReferenceTypeFix A source code browser link to the fix (e.g., a GitHub commit) Note that the fix type is meant for
	// viewing by people using web browsers. Programs interested in analyzing the exact commit range would do better to use
	// the GIT-typed affected[].ranges entries (described above).
	ReferenceTypeFix ReferenceType = "FIX"

	// ReferenceTypeIntroduced A source code browser link to the introduction of the vulnerability (e.g., a GitHub commit)
	// Note that the introduced type is meant for viewing by people using web browsers. Programs interested in analyzing the
	// exact commit range would do better to use the GIT-typed affected[].ranges entries (described above).
	ReferenceTypeIntroduced ReferenceType = "introduced"

	// ReferenceTypePackage A home web page for the package.
	ReferenceTypePackage ReferenceType = "PACKAGE"

	// ReferenceTypeEvidence A demonstration of the validity of a vulnerability claim, e.g. app.any.run replaying the
	// exploitation of the vulnerability.
	ReferenceTypeEvidence ReferenceType = "evidence"

	// ReferenceTypeWeb A web page of some unspecified kind.
	ReferenceTypeWeb ReferenceType = "WEB"
)

// ------------------------------------------------- --------------------------------------------------------------------

// Reference
// Example:
//
//	{
//	  "type": "WEB",
//	  "url": "https://github.com/tensorflow/tensorflow/security/advisories/GHSA-vxv8-r8q2-63xw"
//	}
type Reference struct {

	// 引用的类型
	Type ReferenceType `mapstructure:"type" json:"type" yaml:"type" db:"type" bson:"type" gorm:"column:type"`

	// 具体的引用链接
	URL string `mapstructure:"url" json:"url" yaml:"url" db:"url" bson:"url" gorm:"column:url"`
}

var _ sql.Scanner = &Reference{}
var _ driver.Valuer = &Reference{}

func (x *Reference) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}
	return json.Marshal(x)
}

func (x *Reference) Scan(src any) error {
	if src == nil {
		return nil
	}
	bytes, ok := src.([]byte)
	if !ok {
		return wrapScanError(src, x)
	}
	return json.Unmarshal(bytes, &x)
}
