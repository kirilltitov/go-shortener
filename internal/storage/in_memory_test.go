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

func TestInMemory_GetStats(t *testing.T) {
	ctx := context.Background()
	storage := NewInMemoryStorage(ctx)

	stats, err1 := storage.GetStats(ctx)
	require.NoError(t, err1)
	require.Equal(
		t,
		&Stats{
			Users: 0,
			URLs:  0,
		},
		stats,
	)

	userID1, _ := uuid.NewV6()
	_, err2 := storage.Set(ctx, userID1, "https://ya.ru")
	require.NoError(t, err2)

	userID2, _ := uuid.NewV6()

	_, err3 := storage.Set(ctx, userID2, "https://ya.ru/2")
	require.NoError(t, err3)
	_, err4 := storage.Set(ctx, userID2, "https://ya.ru/3")
	require.NoError(t, err4)

	stats, err5 := storage.GetStats(ctx)
	require.NoError(t, err5)
	require.Equal(
		t,
		&Stats{
			Users: 2,
			URLs:  3,
		},
		stats,
	)
}
