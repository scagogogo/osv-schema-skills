package osv_schema

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strings"
)

// Aliases 一般可能会放漏洞编号啥的
type Aliases []string

var _ sql.Scanner = &Aliases{}
var _ driver.Valuer = &Aliases{}

// GetCVE 获取别名中的CVE编号
func (x Aliases) GetCVE() string {
	for _, s := range x {
		s = strings.ToUpper(s)
		if strings.HasPrefix(s, "CVE-") {
			return s
		}
	}
	return ""
}

// Filter 过滤出需要的编号
func (x Aliases) Filter(filterFunc func(alias string) bool) Aliases {
	slice := make([]string, 0)
	for _, alias := range x {
		if filterFunc(alias) {
			slice = append(slice, alias)
		}
	}
	return slice
}

func (x *Aliases) Scan(src any) error {
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

func (x Aliases) Value() (driver.Value, error) {
	if len(x) == 0 {
		return nil, nil
	}
	// Aliases is []string; json.Marshal cannot fail.
	marshal, _ := json.Marshal(x)
	return string(marshal), nil
}
