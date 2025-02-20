// Модуль main запускает веб-сервер сокращателя ссылок.
// См. [github.com/kirilltitov/go-shortener/internal/config.Config]
// на предмет конфигурационных параметров сервиса.
package main

import (
	"context"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	grpcServer "github.com/kirilltitov/go-shortener/internal/app/grpc"
	httpServer "github.com/kirilltitov/go-shortener/internal/app/http"
	"github.com/kirilltitov/go-shortener/internal/config"
	"github.com/kirilltitov/go-shortener/internal/container"
	"github.com/kirilltitov/go-shortener/internal/logger"
	"github.com/kirilltitov/go-shortener/internal/shortener"
	"github.com/kirilltitov/go-shortener/internal/version"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	v := version.Version{
		BuildVersion: buildVersion,
		BuildDate:    buildDate,
		BuildCommit:  buildCommit,
	}
	v.Print(os.Stdout)

	cfg := config.New()
	ctx := context.Background()
	cnt, err := container.New(ctx, cfg)
	if err != nil {
		panic(err)
	}

	service := shortener.New(cfg, cnt)

	run(service)
}

func run(service shortener.Shortener) {
	wg := &sync.WaitGroup{}

	httpApplication := httpServer.New(service, wg)
	grpcApplication := grpcServer.New(service, wg)

	wg.Add(1)
	go httpApplication.Run()
	wg.Add(1)
	go grpcApplication.Run()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	sig := <-signalChan
	logger.Log.Infof("Received signal: %v", sig)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Log.Info("Shutting down HTTP server")
		if err := httpApplication.Server.Shutdown(shutdownCtx); err != nil {
			logger.Log.WithError(err).Error("Could not shutdown HTTP server properly")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Log.Info("Shutting down gRPC server")
		grpcApplication.Server.GracefulStop()
	}()

	wg.Wait()

	service.Container.Storage.Close()

	logger.Log.Info("Goodbye")
}
