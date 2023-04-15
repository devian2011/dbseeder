package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInArray(t *testing.T) {
	expectedFind := []string{"one", "two", "three", "four"}

	assert.True(t, InArray[string](expectedFind, "one"))
	assert.False(t, InArray[string](expectedFind, "five"))
}
