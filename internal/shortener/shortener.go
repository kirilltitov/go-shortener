package shortener

import (
	"fmt"

	"github.com/kirilltitov/go-shortener/internal/config"
	"github.com/kirilltitov/go-shortener/internal/container"
)

// Shortener является объектом, инкапсулирующим в себе бизнес-логику сервиса по сокращению ссылок.
type Shortener struct {
	Config    config.Config
	Container *container.Container
}

// New создает, конфигурирует и возвращает экземпляр объекта сервиса.
func New(cfg config.Config, cnt *container.Container) Shortener {
	return Shortener{Config: cfg, Container: cnt}
}

// FormatShortURL возвращает полный URL для данного сокращенного идентификатора ссылки.
func (s *Shortener) FormatShortURL(shortURL string) string {
	return fmt.Sprintf("%s/%s", s.Config.BaseURL, shortURL)
}
