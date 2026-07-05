package osv_schema

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAffectedSlice_HasEcosystem(t *testing.T) {
	t.Run("has matching ecosystem", func(t *testing.T) {
		s := AffectedSlice[any, any]{
			{Package: &Package{Ecosystem: EcosystemPyPI, Name: "requests"}},
			{Package: &Package{Ecosystem: EcosystemNpm, Name: "lodash"}},
		}
		assert.True(t, s.HasEcosystem(EcosystemPyPI))
		assert.True(t, s.HasEcosystem(EcosystemNpm))
	})
	t.Run("no matching ecosystem", func(t *testing.T) {
		s := AffectedSlice[any, any]{
			{Package: &Package{Ecosystem: EcosystemPyPI, Name: "requests"}},
		}
		assert.False(t, s.HasEcosystem(EcosystemMaven))
	})
	t.Run("nil package entries are skipped", func(t *testing.T) {
		s := AffectedSlice[any, any]{
			{Package: nil},
			{Package: &Package{Ecosystem: EcosystemGo, Name: "mod"}},
		}
		assert.True(t, s.HasEcosystem(EcosystemGo))
		assert.False(t, s.HasEcosystem(EcosystemPyPI))
	})
	t.Run("empty slice", func(t *testing.T) {
		s := AffectedSlice[any, any]{}
		assert.False(t, s.HasEcosystem(EcosystemPyPI))
	})
}

func TestAffectedSlice_Filter(t *testing.T) {
	s := AffectedSlice[any, any]{
		{Package: &Package{Ecosystem: EcosystemPyPI, Name: "a"}},
		{Package: &Package{Ecosystem: EcosystemMaven, Name: "b"}},
		{Package: &Package{Ecosystem: EcosystemPyPI, Name: "c"}},
	}
	got := s.Filter(func(a *Affected[any, any]) bool {
		return a.Package != nil && a.Package.Ecosystem == EcosystemPyPI
	})
	assert.Equal(t, 2, len(got))
	assert.Equal(t, "a", got[0].Package.Name)
	assert.Equal(t, "c", got[1].Package.Name)
}

func TestAffectedSlice_FilterByEcosystem(t *testing.T) {
	t.Run("nil receiver returns nil", func(t *testing.T) {
		var s AffectedSlice[any, any] = nil
		assert.Nil(t, s.FilterByEcosystem(EcosystemPyPI))
	})
	t.Run("filters to matching ecosystem", func(t *testing.T) {
		s := AffectedSlice[any, any]{
			{Package: &Package{Ecosystem: EcosystemPyPI, Name: "a"}},
			{Package: &Package{Ecosystem: EcosystemMaven, Name: "b"}},
		}
		got := s.FilterByEcosystem(EcosystemPyPI)
		assert.Equal(t, 1, len(got))
		assert.Equal(t, "a", got[0].Package.Name)
	})
	t.Run("no match returns empty slice", func(t *testing.T) {
		s := AffectedSlice[any, any]{
			{Package: &Package{Ecosystem: EcosystemPyPI, Name: "a"}},
		}
		got := s.FilterByEcosystem(EcosystemMaven)
		assert.Equal(t, 0, len(got))
	})
}

func TestAffectedSlice_Scan(t *testing.T) {
	t.Run("nil source", func(t *testing.T) {
		s := AffectedSlice[any, any]{}
		p := &s
		assert.Nil(t, p.Scan(nil))
	})
	t.Run("valid bytes", func(t *testing.T) {
		s := AffectedSlice[any, any]{}
		p := &s
		err := p.Scan([]byte(`[{"package":{"ecosystem":"PyPI","name":"requests"}}]`))
		assert.Nil(t, err)
		assert.Equal(t, 1, len(s))
		assert.Equal(t, EcosystemPyPI, s[0].Package.Ecosystem)
	})
	t.Run("empty bytes returns nil", func(t *testing.T) {
		s := AffectedSlice[any, any]{}
		p := &s
		assert.Nil(t, p.Scan([]byte{}))
		assert.Equal(t, 0, len(s))
	})
	t.Run("non-byte source returns error", func(t *testing.T) {
		s := AffectedSlice[any, any]{}
		p := &s
		err := p.Scan(false)
		assert.NotNil(t, err)
	})
	t.Run("invalid json returns error", func(t *testing.T) {
		s := AffectedSlice[any, any]{}
		p := &s
		err := p.Scan([]byte(`[{bad}`))
		assert.NotNil(t, err)
	})
}

func TestAffectedSlice_Value(t *testing.T) {
	t.Run("non-empty slice", func(t *testing.T) {
		s := AffectedSlice[any, any]{
			{Package: &Package{Ecosystem: EcosystemPyPI, Name: "requests"}},
		}
		v, err := s.Value()
		assert.Nil(t, err)
		str, ok := v.(string)
		assert.True(t, ok)
		var got AffectedSlice[any, any]
		assert.Nil(t, json.Unmarshal([]byte(str), &got))
		assert.Equal(t, 1, len(got))
		assert.Equal(t, "requests", got[0].Package.Name)
	})
	t.Run("empty slice returns nil", func(t *testing.T) {
		s := AffectedSlice[any, any]{}
		v, err := s.Value()
		assert.Nil(t, err)
		assert.Nil(t, v)
	})
	t.Run("marshal error returns err", func(t *testing.T) {
		// a channel value cannot be JSON-marshalled, forcing json.Marshal to error
		s := AffectedSlice[any, any]{
			{EcosystemSpecific: make(chan int)},
		}
		v, err := s.Value()
		assert.NotNil(t, err)
		assert.Nil(t, v)
	})
}

func TestAffected_Value(t *testing.T) {
	t.Run("non-nil affected", func(t *testing.T) {
		a := &Affected[any, any]{Package: &Package{Ecosystem: EcosystemGo, Name: "mod"}}
		v, err := a.Value()
		assert.Nil(t, err)
		b, ok := v.([]byte)
		assert.True(t, ok)
		var got Affected[any, any]
		assert.Nil(t, json.Unmarshal(b, &got))
		assert.Equal(t, EcosystemGo, got.Package.Ecosystem)
	})
	t.Run("nil receiver", func(t *testing.T) {
		var a *Affected[any, any]
		v, err := a.Value()
		assert.Nil(t, err)
		assert.Nil(t, v)
	})
}

func TestAffected_Scan(t *testing.T) {
	t.Run("nil source", func(t *testing.T) {
		a := &Affected[any, any]{}
		assert.Nil(t, a.Scan(nil))
	})
	t.Run("valid bytes", func(t *testing.T) {
		a := &Affected[any, any]{}
		err := a.Scan([]byte(`{"package":{"ecosystem":"Maven","name":"g:a"}}`))
		assert.Nil(t, err)
		assert.Equal(t, EcosystemMaven, a.Package.Ecosystem)
		assert.Equal(t, "g:a", a.Package.Name)
	})
	t.Run("non-byte source returns error", func(t *testing.T) {
		a := &Affected[any, any]{}
		err := a.Scan(7)
		assert.NotNil(t, err)
	})
	t.Run("invalid json returns error", func(t *testing.T) {
		a := &Affected[any, any]{}
		err := a.Scan([]byte(`{bad`))
		assert.NotNil(t, err)
	})
}
