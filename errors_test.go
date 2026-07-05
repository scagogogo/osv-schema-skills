package osv_schema

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapScanError(t *testing.T) {
	err := wrapScanError(123, &Package{})
	assert.NotNil(t, err)
	// message mentions both source and dest type names
	msg := err.Error()
	assert.Contains(t, msg, "can not scan from")
	// int type's reflect name is "int"
	assert.Contains(t, msg, "int")
}

// ensure the returned error is a plain error (not a sentinel) and can be unwrapped/errors.Is
func TestWrapScanError_IsPlainError(t *testing.T) {
	err := wrapScanError("str", &Aliases{})
	assert.False(t, errors.Is(err, assertAnError))
}

var assertAnError = errors.New("sentinel")
