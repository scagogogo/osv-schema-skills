package osv_schema

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRange_Value(t *testing.T) {
	t.Run("non-nil range", func(t *testing.T) {
		r := &Range[any]{Type: RangeTypeEcosystem, Repo: "https://git"}
		v, err := r.Value()
		assert.Nil(t, err)
		b, ok := v.([]byte)
		assert.True(t, ok)
		var got Range[any]
		assert.Nil(t, json.Unmarshal(b, &got))
		assert.Equal(t, RangeType("ECOSYSTEM"), got.Type)
		assert.Equal(t, "https://git", got.Repo)
	})
	t.Run("nil receiver", func(t *testing.T) {
		var r *Range[any]
		v, err := r.Value()
		assert.Nil(t, err)
		assert.Nil(t, v)
	})
}

func TestRange_Scan(t *testing.T) {
	t.Run("nil source", func(t *testing.T) {
		r := &Range[any]{}
		assert.Nil(t, r.Scan(nil))
	})
	t.Run("valid bytes", func(t *testing.T) {
		r := &Range[any]{}
		err := r.Scan([]byte(`{"type":"SEMVER","repo":"r"}`))
		assert.Nil(t, err)
		assert.Equal(t, RangeType("SEMVER"), r.Type)
		assert.Equal(t, "r", r.Repo)
	})
	t.Run("non-byte source returns error", func(t *testing.T) {
		r := &Range[any]{}
		err := r.Scan(uint(1))
		assert.NotNil(t, err)
	})
	t.Run("invalid json returns error", func(t *testing.T) {
		r := &Range[any]{}
		err := r.Scan([]byte(`{bad`))
		assert.NotNil(t, err)
	})
}
