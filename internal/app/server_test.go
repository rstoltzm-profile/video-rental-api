package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	got := 1
	want := 1
	assert.Equal(t, want, got)
}
