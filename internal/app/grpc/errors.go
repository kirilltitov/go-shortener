package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrUnauthorized является ошибкой отсутствия авторизации в ходе обработки gRPC-запроса.
var ErrUnauthorized = status.Error(codes.Unauthenticated, "unauthorized")

// ErrInternal является внутренней ошибкой выполнения gRPC-запроса.
var ErrInternal = status.Error(codes.Internal, "internal error")
