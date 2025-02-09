package shortener

import (
	"context"
	"net"

	"github.com/kirilltitov/go-shortener/internal/storage"
)

// GetStats возвращает статистику использования сервиса, см. структуру [storage.Stats].
func (s *Shortener) GetStats(ctx context.Context, clientIP string) (*storage.Stats, error) {
	if !s.isTrustedClientIP(clientIP) {
		return nil, ErrorUnauthorized
	}

	return s.Container.Storage.GetStats(ctx)
}

func (s *Shortener) isTrustedClientIP(clientIP string) bool {
	if s.Config.TrustedSubnet == nil {
		return false
	}

	ip := net.ParseIP(clientIP)
	if ip == nil {
		return false
	}

	return s.Config.TrustedSubnet.Contains(ip)
}
