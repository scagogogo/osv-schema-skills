package osv_schema

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvent_IsIntroduced(t *testing.T) {
	assert.True(t, (&Event{Introduced: "1.0.0"}).IsIntroduced())
	assert.False(t, (&Event{}).IsIntroduced())
}

func TestEvent_IsFixed(t *testing.T) {
	assert.True(t, (&Event{Fixed: "1.0.1"}).IsFixed())
	assert.False(t, (&Event{}).IsFixed())
}

func TestEvent_IsLastAffected(t *testing.T) {
	assert.True(t, (&Event{LastAffected: "2.0.0"}).IsLastAffected())
	assert.False(t, (&Event{}).IsLastAffected())
}

func TestEvent_IsLimit(t *testing.T) {
	assert.True(t, (&Event{Limit: "v1"}).IsLimit())
	assert.False(t, (&Event{}).IsLimit())
}

func TestEvent_Value(t *testing.T) {
	t.Run("non-nil event", func(t *testing.T) {
		e := &Event{Introduced: "1.0.0", Fixed: "1.0.1"}
		v, err := e.Value()
		assert.Nil(t, err)
		b, ok := v.([]byte)
		assert.True(t, ok)
		var got Event
		assert.Nil(t, json.Unmarshal(b, &got))
		assert.Equal(t, "1.0.0", got.Introduced)
		assert.Equal(t, "1.0.1", got.Fixed)
	})
	t.Run("nil receiver", func(t *testing.T) {
		var e *Event
		v, err := e.Value()
		assert.Nil(t, err)
		assert.Nil(t, v)
	})
}

func TestEvent_Scan(t *testing.T) {
	t.Run("nil source", func(t *testing.T) {
		e := &Event{}
		assert.Nil(t, e.Scan(nil))
	})
	t.Run("valid bytes", func(t *testing.T) {
		e := &Event{}
		err := e.Scan([]byte(`{"introduced":"0","fixed":"1.0.1"}`))
		assert.Nil(t, err)
		assert.Equal(t, "0", e.Introduced)
		assert.Equal(t, "1.0.1", e.Fixed)
	})
	t.Run("non-byte source returns error", func(t *testing.T) {
		e := &Event{}
		err := e.Scan("string-not-bytes")
		assert.NotNil(t, err)
	})
	t.Run("invalid json returns error", func(t *testing.T) {
		e := &Event{}
		err := e.Scan([]byte(`{bad`))
		assert.NotNil(t, err)
	})
}
