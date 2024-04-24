package shortener

import (
	"fmt"

	"github.com/kirilltitov/go-shortener/internal/config"
	"github.com/kirilltitov/go-shortener/internal/container"
)

type Shortener struct {
	config    config.Config
	container *container.Container
}

func New(cfg config.Config, cnt *container.Container) Shortener {
	return Shortener{config: cfg, container: cnt}
}

func (s *Shortener) FormatShortURL(shortURL string) string {
	return fmt.Sprintf("%s/%s", s.config.BaseURL, shortURL)
}
