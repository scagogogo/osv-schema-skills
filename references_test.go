package osv_schema

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReferences_FilterByType(t *testing.T) {
	refs := References{
		{Type: ReferenceTypeAdvisory, URL: "https://adv"},
		{Type: ReferenceTypeFix, URL: "https://fix"},
		{Type: ReferenceTypeWeb, URL: "https://web"},
		{Type: ReferenceTypeFix, URL: "https://fix2"},
	}
	t.Run("single type", func(t *testing.T) {
		got := refs.FilterByType(ReferenceTypeFix)
		assert.Equal(t, 2, len(got))
		assert.Equal(t, "https://fix", got[0].URL)
		assert.Equal(t, "https://fix2", got[1].URL)
	})
	t.Run("multiple types OR semantics", func(t *testing.T) {
		got := refs.FilterByType(ReferenceTypeAdvisory, ReferenceTypeWeb)
		assert.Equal(t, 2, len(got))
	})
	t.Run("no types returns nil", func(t *testing.T) {
		got := refs.FilterByType()
		assert.Nil(t, got)
	})
	t.Run("no match returns empty slice", func(t *testing.T) {
		got := refs.FilterByType(ReferenceTypeReport)
		assert.Equal(t, 0, len(got))
	})
	t.Run("dedupes duplicate type args", func(t *testing.T) {
		got := refs.FilterByType(ReferenceTypeFix, ReferenceTypeFix)
		assert.Equal(t, 2, len(got))
	})
}

func TestReferences_Scan(t *testing.T) {
	t.Run("nil source", func(t *testing.T) {
		r := References{}
		p := &r
		assert.Nil(t, p.Scan(nil))
	})
	t.Run("valid bytes", func(t *testing.T) {
		r := References{}
		p := &r
		err := p.Scan([]byte(`[{"type":"FIX","url":"https://x"}]`))
		assert.Nil(t, err)
		assert.Equal(t, 1, len(r))
		assert.Equal(t, ReferenceTypeFix, r[0].Type)
	})
	t.Run("empty bytes returns nil", func(t *testing.T) {
		r := References{}
		p := &r
		assert.Nil(t, p.Scan([]byte{}))
		assert.Equal(t, 0, len(r))
	})
	t.Run("non-byte source returns scan error", func(t *testing.T) {
		r := References{}
		p := &r
		err := p.Scan(3.14)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "scan error")
	})
	t.Run("invalid json returns error", func(t *testing.T) {
		r := References{}
		p := &r
		err := p.Scan([]byte(`[bad`))
		assert.NotNil(t, err)
	})
}

func TestReferences_Value(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		r := References{{Type: ReferenceTypeWeb, URL: "https://x"}}
		v, err := r.Value()
		assert.Nil(t, err)
		s, ok := v.(string)
		assert.True(t, ok)
		var got References
		assert.Nil(t, json.Unmarshal([]byte(s), &got))
		assert.Equal(t, ReferenceTypeWeb, got[0].Type)
	})
	t.Run("empty returns nil", func(t *testing.T) {
		r := References{}
		v, err := r.Value()
		assert.Nil(t, err)
		assert.Nil(t, v)
	})
}

func TestReference_Value(t *testing.T) {
	t.Run("non-nil", func(t *testing.T) {
		r := &Reference{Type: ReferenceTypeFix, URL: "https://f"}
		v, err := r.Value()
		assert.Nil(t, err)
		b, ok := v.([]byte)
		assert.True(t, ok)
		var got Reference
		assert.Nil(t, json.Unmarshal(b, &got))
		assert.Equal(t, ReferenceTypeFix, got.Type)
	})
	t.Run("nil receiver", func(t *testing.T) {
		var r *Reference
		v, err := r.Value()
		assert.Nil(t, err)
		assert.Nil(t, v)
	})
}

func TestReference_Scan(t *testing.T) {
	t.Run("nil source", func(t *testing.T) {
		r := &Reference{}
		assert.Nil(t, r.Scan(nil))
	})
	t.Run("valid bytes", func(t *testing.T) {
		r := &Reference{}
		err := r.Scan([]byte(`{"type":"WEB","url":"https://w"}`))
		assert.Nil(t, err)
		assert.Equal(t, ReferenceTypeWeb, r.Type)
		assert.Equal(t, "https://w", r.URL)
	})
	t.Run("non-byte source returns error", func(t *testing.T) {
		r := &Reference{}
		err := r.Scan(complex(1, 1))
		assert.NotNil(t, err)
	})
	t.Run("invalid json returns error", func(t *testing.T) {
		r := &Reference{}
		err := r.Scan([]byte(`{bad`))
		assert.NotNil(t, err)
	})
}
