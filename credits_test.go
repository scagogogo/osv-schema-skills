package osv_schema

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCredits_Value(t *testing.T) {
	t.Run("non-nil credits", func(t *testing.T) {
		c := &Credits{Name: "Alice", Type: "FINDER", Contact: []string{"a@x"}}
		v, err := c.Value()
		assert.Nil(t, err)
		b, ok := v.([]byte)
		assert.True(t, ok)
		var got Credits
		assert.Nil(t, json.Unmarshal(b, &got))
		assert.Equal(t, "Alice", got.Name)
		assert.Equal(t, "FINDER", got.Type)
		assert.Equal(t, []string{"a@x"}, got.Contact)
	})
	t.Run("nil receiver", func(t *testing.T) {
		var c *Credits
		v, err := c.Value()
		assert.Nil(t, err)
		assert.Nil(t, v)
	})
}

func TestCredits_Scan(t *testing.T) {
	t.Run("nil source", func(t *testing.T) {
		c := &Credits{}
		assert.Nil(t, c.Scan(nil))
	})
	t.Run("valid bytes", func(t *testing.T) {
		c := &Credits{}
		err := c.Scan([]byte(`{"name":"Bob","type":"REPORTER","contact":["b@y"]}`))
		assert.Nil(t, err)
		assert.Equal(t, "Bob", c.Name)
		assert.Equal(t, "REPORTER", c.Type)
		assert.Equal(t, []string{"b@y"}, c.Contact)
	})
	t.Run("non-byte source returns error", func(t *testing.T) {
		c := &Credits{}
		err := c.Scan(int64(1))
		assert.NotNil(t, err)
	})
	t.Run("invalid json returns error", func(t *testing.T) {
		c := &Credits{}
		err := c.Scan([]byte(`{bad`))
		assert.NotNil(t, err)
	})
}
