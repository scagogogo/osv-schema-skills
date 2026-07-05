package osv_schema

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type Related []string

var _ sql.Scanner = &Related{}
var _ driver.Valuer = &Related{}

func (x *Related) Scan(src any) error {
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

func (x Related) Value() (driver.Value, error) {
	if len(x) == 0 {
		return nil, nil
	}
	// Related is []string; json.Marshal cannot fail.
	marshal, _ := json.Marshal(x)
	return string(marshal), nil
}
