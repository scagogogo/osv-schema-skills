package osv_schema

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/golang-infrastructure/go-pointer"
	"strconv"
)

// ------------------------------------------------ ---------------------------------------------------------------------

type SeveritySlice []*Severity

var _ sql.Scanner = &SeveritySlice{}
var _ driver.Valuer = &SeveritySlice{}

func (x SeveritySlice) GetCVSS3() *Severity {
	for _, s := range x {
		if s.Type == SeverityTypeCVSS3 {
			return s
		}
	}
	return nil
}

func (x SeveritySlice) GetCVSS2() *Severity {
	for _, s := range x {
		if s.Type == SeverityTypeCVSS2 {
			return s
		}
	}
	return nil
}

func (x *SeveritySlice) Scan(src any) error {
	if src == nil {
		return nil
	}
	bytes, ok := src.([]byte)
	if !ok {
		return wrapScanError(src, x)
	}
	if len(bytes) == 0 {
		return nil
	}
	return json.Unmarshal(bytes, &x)
}

func (x SeveritySlice) Value() (driver.Value, error) {
	if len(x) == 0 {
		return nil, nil
	}
	// SeveritySlice only contains Severity pointers with string-typed exported
	// fields, so json.Marshal cannot fail.
	marshal, _ := json.Marshal(x)
	return string(marshal), nil
}

// ------------------------------------------------ ---------------------------------------------------------------------

type SeverityType string

const (

	// SeverityTypeCVSS2 e.g."AV:L/AC:M/Au:N/C:N/I:P/A:C"
	SeverityTypeCVSS2 SeverityType = "CVSS_V2"

	// SeverityTypeCVSS3 CVSS:3.1/AV:N/AC:H/PR:N/UI:N/S:C/C:H/I:N/A:N
	SeverityTypeCVSS3 SeverityType = "CVSS_V3"
)

// Severity
// Example:
//
//	{
//	  "type": "CVSS_V3",
//	  "score": "CVSS:3.1/AV:N/AC:H/PR:N/UI:N/S:U/C:N/I:N/A:H"
//	}
//
// Document: https://ossf.github.io/osv-schema/#severity-field
type Severity struct {
	Type  SeverityType `mapstructure:"type" json:"type" yaml:"type" db:"type" bson:"type" gorm:"column:type"`
	Score string       `mapstructure:"score" json:"score" yaml:"score" db:"score" bson:"score" gorm:"column:score"`

	score *float64
	err   error
}

var _ sql.Scanner = &Severity{}
var _ driver.Valuer = &Severity{}

// ------------------------------------------------- --------------------------------------------------------------------

func (x *Severity) GetScore() float64 {
	score, _ := x.GetScoreAsFloat()
	return score
}

func (x *Severity) GetScoreAsPointer() *float64 {
	score, err := x.GetScoreAsFloat()
	if err != nil {
		return nil
	} else {
		return pointer.ToPointer(score)
	}
}

func (x *Severity) GetScoreAsFloat() (float64, error) {
	if x.err != nil {
		return 0, x.err
	} else if x.score != nil {
		return pointer.FromPointer(x.score), nil
	}
	if x.Score == "" {
		x.err = fmt.Errorf("score can not be empty")
		return 0, x.err
	}
	score, err := strconv.ParseFloat(x.Score, 64)
	if err != nil {
		x.err = err
		return 0, err
	}
	x.score = pointer.ToPointer(score)
	return score, nil
}

// ------------------------------------------------- --------------------------------------------------------------------

func (x *Severity) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}
	return json.Marshal(x)
}

func (x *Severity) Scan(src any) error {
	if src == nil {
		return nil
	}
	bytes, ok := src.([]byte)
	if !ok {
		return wrapScanError(src, x)
	}
	return json.Unmarshal(bytes, &x)
}

// ------------------------------------------------- --------------------------------------------------------------------
