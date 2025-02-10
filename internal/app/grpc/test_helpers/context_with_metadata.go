package test_helpers

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func NewContextWithValue(key string, value string) context.Context {
	return metadata.NewIncomingContext(
		grpc.NewContextWithServerTransportStream(context.Background(), &MockServerTransportStream{}),
		metadata.New(map[string]string{key: value}),
	)
}
