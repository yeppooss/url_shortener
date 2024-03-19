package random

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestNewRandomString(t *testing.T) {
	newString := NewRandomString(6)

	assert.Equal(t, len(newString), 6)
	assert.NotEqual(t, len(newString), 5)
}
