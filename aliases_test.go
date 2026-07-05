package osv_schema

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAliases_GetCVE(t *testing.T) {
	t.Run("has CVE uppercase", func(t *testing.T) {
		a := Aliases{"GHSA-abc", "CVE-2024-1234"}
		assert.Equal(t, "CVE-2024-1234", a.GetCVE())
	})
	t.Run("lowercase cve is upper-cased before matching and return", func(t *testing.T) {
		a := Aliases{"cve-2024-5678"}
		assert.Equal(t, "CVE-2024-5678", a.GetCVE())
	})
	t.Run("mixed case", func(t *testing.T) {
		a := Aliases{"CvE-2024-0001"}
		assert.Equal(t, "CVE-2024-0001", a.GetCVE())
	})
	t.Run("no CVE returns empty", func(t *testing.T) {
		a := Aliases{"GHSA-abc", "RHSA-xyz"}
		assert.Equal(t, "", a.GetCVE())
	})
	t.Run("empty aliases", func(t *testing.T) {
		a := Aliases{}
		assert.Equal(t, "", a.GetCVE())
	})
	t.Run("returns first CVE only", func(t *testing.T) {
		a := Aliases{"CVE-2024-1", "CVE-2024-2"}
		assert.Equal(t, "CVE-2024-1", a.GetCVE())
	})
}

func TestAliases_Filter(t *testing.T) {
	t.Run("keeps matching", func(t *testing.T) {
		a := Aliases{"CVE-1", "GHSA-2", "CVE-3"}
		got := a.Filter(func(s string) bool {
			return len(s) > 0 && s[0] == 'C'
		})
		assert.Equal(t, Aliases{"CVE-1", "CVE-3"}, got)
	})
	t.Run("none match returns empty", func(t *testing.T) {
		a := Aliases{"GHSA-1"}
		got := a.Filter(func(s string) bool { return false })
		assert.Equal(t, Aliases{}, got)
	})
	t.Run("empty input", func(t *testing.T) {
		a := Aliases{}
		got := a.Filter(func(s string) bool { return true })
		assert.Equal(t, Aliases{}, got)
	})
}

func TestAliases_Scan(t *testing.T) {
	t.Run("nil source", func(t *testing.T) {
		a := Aliases{}
		p := &a
		assert.Nil(t, p.Scan(nil))
	})
	t.Run("valid bytes", func(t *testing.T) {
		a := Aliases{}
		p := &a
		err := p.Scan([]byte(`["CVE-1","GHSA-2"]`))
		assert.Nil(t, err)
		assert.Equal(t, Aliases{"CVE-1", "GHSA-2"}, a)
	})
	t.Run("empty bytes returns nil", func(t *testing.T) {
		a := Aliases{}
		p := &a
		assert.Nil(t, p.Scan([]byte{}))
		assert.Equal(t, Aliases{}, a)
	})
	t.Run("non-byte source returns error", func(t *testing.T) {
		a := Aliases{}
		p := &a
		err := p.Scan(42)
		assert.NotNil(t, err)
	})
	t.Run("invalid json returns error", func(t *testing.T) {
		a := Aliases{}
		p := &a
		err := p.Scan([]byte(`[bad`))
		assert.NotNil(t, err)
	})
}

func TestAliases_Value(t *testing.T) {
	t.Run("non-empty aliases", func(t *testing.T) {
		a := Aliases{"CVE-1", "GHSA-2"}
		v, err := a.Value()
		assert.Nil(t, err)
		s, ok := v.(string)
		assert.True(t, ok)
		var got []string
		assert.Nil(t, json.Unmarshal([]byte(s), &got))
		assert.Equal(t, []string{"CVE-1", "GHSA-2"}, got)
	})
	t.Run("empty aliases returns nil", func(t *testing.T) {
		a := Aliases{}
		v, err := a.Value()
		assert.Nil(t, err)
		assert.Nil(t, v)
	})
}
