package osv_schema

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// ------------------------------------------------ ---------------------------------------------------------------------

// AffectedSlice 表示一个影响范围的集合
type AffectedSlice[EcosystemSpecific, DatabaseSpecific any] []*Affected[EcosystemSpecific, DatabaseSpecific]

var _ sql.Scanner = &AffectedSlice[any, any]{}
var _ driver.Valuer = &AffectedSlice[any, any]{}

func (x *AffectedSlice[EcosystemSpecific, DatabaseSpecific]) Scan(src any) error {
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

func (x AffectedSlice[EcosystemSpecific, DatabaseSpecific]) Value() (driver.Value, error) {
	if len(x) == 0 {
		return nil, nil
	}
	marshal, err := json.Marshal(x)
	if err != nil {
		return nil, err
	}
	return string(marshal), nil
}

// HasEcosystem 判断被影响到的包是否有包含给定的包管理器的，一般用于过滤
func (x AffectedSlice[EcosystemSpecific, DatabaseSpecific]) HasEcosystem(ecosystem Ecosystem) bool {
	// 这里认为这个数组不会特别大，所以就O(n)扫描了
	for _, item := range x {
		if item.Package != nil && item.Package.Ecosystem == ecosystem {
			return true
		}
	}
	return false
}

// Filter 过滤影响范围
func (x AffectedSlice[EcosystemSpecific, DatabaseSpecific]) Filter(filterFunc func(affected *Affected[EcosystemSpecific, DatabaseSpecific]) bool) AffectedSlice[EcosystemSpecific, DatabaseSpecific] {
	slice := make([]*Affected[EcosystemSpecific, DatabaseSpecific], 0)
	for _, item := range x {
		if filterFunc(item) {
			slice = append(slice, item)
		}
	}
	return slice
}

// FilterByEcosystem 根据ecosystem过滤影响范围
func (x AffectedSlice[EcosystemSpecific, DatabaseSpecific]) FilterByEcosystem(ecosystem Ecosystem) AffectedSlice[EcosystemSpecific, DatabaseSpecific] {
	if x == nil {
		return nil
	}
	return x.Filter(func(affected *Affected[EcosystemSpecific, DatabaseSpecific]) bool {
		return affected.Package.Ecosystem == ecosystem
	})
}

// ------------------------------------------------ ---------------------------------------------------------------------

// Affected 漏洞的某个影响范围，它可能会影响到很多个版本范围，这表示其中一个
// Example:
// "affected": [
//
//	{
//	  "package": {
//	    "ecosystem": "RubyGems",
//	    "name": "sprout"
//	  },
//	  "ranges": [
//	    {
//	      "type": "ECOSYSTEM",
//	      "events": [
//	        {
//	          "introduced": "0"
//	        },
//	        {
//	          "last_affected": "0.7.246"
//	        }
//	      ]
//	    }
//	  ]
//	}
//
// ],
type Affected[EcosystemSpecific, DatabaseSpecific any] struct {

	// 被此漏洞影响到的包
	Package *Package `mapstructure:"package" json:"package" yaml:"package" db:"package" bson:"package" gorm:"column:package;serializer:json"`

	// 被影响到的这个包的哪些版本，通常是版本区间
	Ranges []*Range[DatabaseSpecific] `mapstructure:"ranges" json:"ranges" yaml:"ranges" db:"ranges" bson:"ranges" gorm:"column:ranges;serializer:json"`

	// 可选的严重级别
	Severity []*Severity `mapstructure:"severity" json:"severity" yaml:"severity" db:"severity" bson:"severity" gorm:"column:severity;serializer:json"`

	// 枚举出每一个受影响的版本
	Versions []string `mapstructure:"versions" json:"versions" yaml:"versions" db:"versions" bson:"versions" gorm:"column:versions;serializer:json"`

	// 由包管理器决定
	EcosystemSpecific EcosystemSpecific `mapstructure:"ecosystem_specific" json:"ecosystem_specific" yaml:"ecosystem_specific" db:"ecosystem_specific" bson:"ecosystem_specific" gorm:"column:ecosystem_specific;serializer:json"`

	// 由具体实现的数据库决定
	DatabaseSpecific DatabaseSpecific `mapstructure:"database_specific" json:"database_specific" yaml:"database_specific" db:"database_specific" bson:"database_specific" gorm:"column:database_specific;serializer:json"`
}

var _ sql.Scanner = &Affected[any, any]{}
var _ driver.Valuer = &Affected[any, any]{}

func (x *Affected[EcosystemSpecific, DatabaseSpecific]) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}
	return json.Marshal(x)
}

func (x *Affected[EcosystemSpecific, DatabaseSpecific]) Scan(src any) error {
	if src == nil {
		return nil
	}
	bytes, ok := src.([]byte)
	if !ok {
		return wrapScanError(src, x)
	}
	return json.Unmarshal(bytes, &x)
}
