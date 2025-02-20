package grpc

import (
	"context"
	"errors"
	"net"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"github.com/kirilltitov/go-shortener/internal/app"
	"github.com/kirilltitov/go-shortener/internal/app/grpc/gen"
	"github.com/kirilltitov/go-shortener/internal/app/grpc/interceptors"
	"github.com/kirilltitov/go-shortener/internal/logger"
	"github.com/kirilltitov/go-shortener/internal/shortener"
)

type Application struct {
	gen.UnimplementedShortenerServer

	Shortener shortener.Shortener
	Server    *grpc.Server

	wg *sync.WaitGroup
}

// New создает и возвращает сконфигурированный объект gRPC-приложения сервиса.
func New(s shortener.Shortener, wg *sync.WaitGroup) *Application {
	a := &Application{
		Shortener: s,
		wg:        wg,
	}

	return a
}

// Run запускает gRPC-сервер приложения.
func (a *Application) Run() {
	defer a.wg.Done()

	logger.Log.Infof("Starting a gRPC server at %s", a.Shortener.Config.GrpcAddress)

	listen, err := net.Listen("tcp", a.Shortener.Config.GrpcAddress)
	if err != nil {
		panic(err)
	}

	a.Server = grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptors.UnaryLoggerInterceptor,
		interceptors.UnaryAuthInterceptor,
	))

	gen.RegisterShortenerServer(a.Server, a)

	if err2 := a.Server.Serve(listen); err2 != nil {
		if errors.Is(err2, grpc.ErrServerStopped) {
			logger.Log.Info("gRPC server shutdown")
		} else {
			panic(err2)
		}
	}
}

func getUserID(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(app.CtxUserIDKey{}).(uuid.UUID)
	return userID, ok
}
