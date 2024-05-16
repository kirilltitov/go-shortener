package storage

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemory(t *testing.T) {
	ctx := context.Background()
	storage := NewInMemoryStorage(ctx)

	userID, _ := uuid.NewV6()
	shortURL, err := storage.Set(ctx, userID, "https://ya.ru")
	require.NoError(t, err)

	result, err := storage.Get(ctx, shortURL)
	require.NoError(t, err)

	assert.Equal(t, "https://ya.ru", result)
}
