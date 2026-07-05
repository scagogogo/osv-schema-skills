package osv_schema

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testFixture = "test_data/GHSA-vxv8-r8q2-63xw.json"

func TestUnmarshalFromJson(t *testing.T) {
	t.Run("valid json", func(t *testing.T) {
		bytes, err := os.ReadFile(testFixture)
		assert.Nil(t, err)
		v, err := UnmarshalFromJson[any, any](bytes)
		assert.Nil(t, err)
		assert.NotNil(t, v)
		assert.Equal(t, "GHSA-vxv8-r8q2-63xw", v.ID)
		assert.Equal(t, "1.4.0", v.SchemaVersion)
		assert.True(t, len(v.Aliases) > 0)
	})
	t.Run("invalid json returns error and nil pointer", func(t *testing.T) {
		v, err := UnmarshalFromJson[any, any]([]byte(`{invalid`))
		assert.NotNil(t, err)
		assert.Nil(t, v)
	})
	t.Run("valid but empty object returns pointer with empty fields", func(t *testing.T) {
		v, err := UnmarshalFromJson[any, any]([]byte(`{}`))
		assert.Nil(t, err)
		assert.NotNil(t, v)
		assert.Equal(t, "", v.ID)
	})
}

func TestUnmarshalFromJsonFile(t *testing.T) {
	t.Run("valid file", func(t *testing.T) {
		v, err := UnmarshalFromJsonFile[any, any](testFixture)
		assert.Nil(t, err)
		assert.NotNil(t, v)
		assert.Equal(t, "GHSA-vxv8-r8q2-63xw", v.ID)
	})
	t.Run("non-existent file returns error and nil pointer", func(t *testing.T) {
		v, err := UnmarshalFromJsonFile[any, any]("does-not-exist.json")
		assert.NotNil(t, err)
		assert.Nil(t, v)
	})
	t.Run("file with invalid json returns error", func(t *testing.T) {
		dir := t.TempDir()
		p := filepath.Join(dir, "bad.json")
		assert.Nil(t, os.WriteFile(p, []byte(`{not json`), 0o644))
		v, err := UnmarshalFromJsonFile[any, any](p)
		assert.NotNil(t, err)
		assert.Nil(t, v)
	})
}
