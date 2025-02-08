// Модуль main запускает веб-сервер сокращателя ссылок.
// См. [github.com/kirilltitov/go-shortener/internal/config.Config]
// на предмет конфигурационных параметров сервиса.
package main

import (
	"context"
	_ "net/http/pprof"
	"os"

	"github.com/kirilltitov/go-shortener/internal/app/http"
	"github.com/kirilltitov/go-shortener/internal/config"
	"github.com/kirilltitov/go-shortener/internal/logger"
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

	a, err := http.New(ctx, cfg)
	if err != nil {
		panic(err)
	}

	logger.Log.Infof("Starting server at %s", cfg.ServerAddress)

	a.Run()
}
