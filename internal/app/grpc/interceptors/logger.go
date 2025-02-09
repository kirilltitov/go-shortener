package interceptors

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var Log = logrus.New()

// UnaryLoggerInterceptor является gRPC-интерцептором для логирования информации о запросе.
func UnaryLoggerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	resp, err := handler(ctx, req)

	duration := time.Since(start)

	inf := status.Convert(err)
	var statusCode, errMsg string

	if err != nil {
		statusCode = inf.Code().String()
		errMsg = inf.Message()
	} else {
		statusCode = "OK"
	}

	Log.WithFields(logrus.Fields{
		"method":      info.FullMethod,
		"status":      statusCode,
		"errMsg":      errMsg,
		"duration_μs": duration.Microseconds(),
	}).Info("Served gRPC request")

	return resp, err
}
