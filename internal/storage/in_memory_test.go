package storage

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInMemory(t *testing.T) {
	var storage Storage = InMemory{}

	storage.Set(1337, "foo")

	result, ok := storage.Get(1337)
	require.True(t, ok)

	assert.Equal(t, "foo", result)
}
