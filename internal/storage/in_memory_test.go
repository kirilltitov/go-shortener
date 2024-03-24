package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemory(t *testing.T) {
	storage := InMemory{}

	storage.Set(1337, "foo")

	result, ok := storage.Get(1337)
	require.True(t, ok)

	assert.Equal(t, "foo", result)
}
