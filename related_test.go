package osv_schema

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRelated_Scan(t *testing.T) {
	t.Run("nil source", func(t *testing.T) {
		r := Related{}
		p := &r
		assert.Nil(t, p.Scan(nil))
	})
	t.Run("valid bytes", func(t *testing.T) {
		r := Related{}
		p := &r
		err := p.Scan([]byte(`["x","y"]`))
		assert.Nil(t, err)
		assert.Equal(t, Related{"x", "y"}, r)
	})
	t.Run("empty bytes returns nil", func(t *testing.T) {
		r := Related{}
		p := &r
		assert.Nil(t, p.Scan([]byte{}))
		assert.Equal(t, Related{}, r)
	})
	t.Run("non-byte source returns error", func(t *testing.T) {
		r := Related{}
		p := &r
		err := p.Scan("nope")
		assert.NotNil(t, err)
	})
	t.Run("invalid json returns error", func(t *testing.T) {
		r := Related{}
		p := &r
		err := p.Scan([]byte(`[bad`))
		assert.NotNil(t, err)
	})
}

func TestRelated_Value(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		r := Related{"a", "b"}
		v, err := r.Value()
		assert.Nil(t, err)
		s, ok := v.(string)
		assert.True(t, ok)
		var got Related
		assert.Nil(t, json.Unmarshal([]byte(s), &got))
		assert.Equal(t, Related{"a", "b"}, got)
	})
	t.Run("empty returns nil", func(t *testing.T) {
		r := Related{}
		v, err := r.Value()
		assert.Nil(t, err)
		assert.Nil(t, v)
	})
}
