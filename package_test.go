package osv_schema

import (
	"database/sql/driver"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackage_IsMaven(t *testing.T) {
	t.Run("maven package", func(t *testing.T) {
		p := &Package{Ecosystem: EcosystemMaven, Name: "org.apache:lib"}
		assert.True(t, p.IsMaven())
	})
	t.Run("non-maven package", func(t *testing.T) {
		p := &Package{Ecosystem: EcosystemPyPI, Name: "requests"}
		assert.False(t, p.IsMaven())
	})
	t.Run("empty ecosystem", func(t *testing.T) {
		p := &Package{Name: "requests"}
		assert.False(t, p.IsMaven())
	})
}

func TestPackage_GetGroupID(t *testing.T) {
	t.Run("maven name with colon", func(t *testing.T) {
		p := &Package{Name: "org.apache.logging.log4j:log4j-core"}
		assert.Equal(t, "org.apache.logging.log4j", p.GetGroupID())
	})
	t.Run("name without colon", func(t *testing.T) {
		p := &Package{Name: "requests"}
		assert.Equal(t, "", p.GetGroupID())
	})
	t.Run("nil receiver", func(t *testing.T) {
		var p *Package
		assert.Equal(t, "", p.GetGroupID())
	})
	t.Run("name with multiple colons splits on first", func(t *testing.T) {
		p := &Package{Name: "g:a:b"}
		assert.Equal(t, "g", p.GetGroupID())
	})
}

func TestPackage_GetArtifactID(t *testing.T) {
	t.Run("maven name with colon", func(t *testing.T) {
		p := &Package{Name: "org.apache.logging.log4j:log4j-core"}
		assert.Equal(t, "log4j-core", p.GetArtifactID())
	})
	t.Run("name without colon", func(t *testing.T) {
		p := &Package{Name: "requests"}
		assert.Equal(t, "", p.GetArtifactID())
	})
	t.Run("nil receiver", func(t *testing.T) {
		var p *Package
		assert.Equal(t, "", p.GetArtifactID())
	})
	t.Run("name with multiple colons keeps remainder", func(t *testing.T) {
		p := &Package{Name: "g:a:b"}
		assert.Equal(t, "a:b", p.GetArtifactID())
	})
}

func TestPackage_Value(t *testing.T) {
	t.Run("non-nil package", func(t *testing.T) {
		p := &Package{Ecosystem: EcosystemNpm, Name: "lodash"}
		v, err := p.Value()
		assert.Nil(t, err)
		assert.NotNil(t, v)
		// should be valid JSON bytes
		s, ok := v.([]byte)
		assert.True(t, ok)
		var got Package
		assert.Nil(t, json.Unmarshal(s, &got))
		assert.Equal(t, EcosystemNpm, got.Ecosystem)
		assert.Equal(t, "lodash", got.Name)
	})
	t.Run("nil receiver", func(t *testing.T) {
		var p *Package
		v, err := p.Value()
		assert.Nil(t, err)
		assert.Nil(t, v)
	})
}

func TestPackage_Scan(t *testing.T) {
	t.Run("nil source", func(t *testing.T) {
		p := &Package{}
		assert.Nil(t, p.Scan(nil))
	})
	t.Run("valid bytes", func(t *testing.T) {
		p := &Package{}
		err := p.Scan([]byte(`{"ecosystem":"Maven","name":"g:a"}`))
		assert.Nil(t, err)
		assert.Equal(t, EcosystemMaven, p.Ecosystem)
		assert.Equal(t, "g:a", p.Name)
	})
	t.Run("non-byte source returns error", func(t *testing.T) {
		p := &Package{}
		err := p.Scan(123)
		assert.NotNil(t, err)
	})
	t.Run("invalid json returns error", func(t *testing.T) {
		p := &Package{}
		err := p.Scan([]byte(`{invalid`))
		assert.NotNil(t, err)
	})
}

// compile-time assertion that Package implements sql.Scanner and driver.Valuer
var _ driver.Valuer = (*Package)(nil)
