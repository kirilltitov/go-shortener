package shortener

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/kirilltitov/go-shortener/internal/logger"
	"github.com/kirilltitov/go-shortener/internal/utils"
)

// DeleteUserURLs асинхронно удаляет произвольное количество коротких ссылок указанного пользователя.
// Возвращает канал, который может возвращать ошибки удаления коротких ссылок, либо закроется после успешного удаления
// всего переданного списка ссылок.
func (s *Shortener) DeleteUserURLs(
	ctx context.Context,
	doneCh chan struct{},
	userID uuid.UUID,
	URLs []string,
	wg *sync.WaitGroup,
) chan error {
	inputCh := utils.Generator(doneCh, URLs)

	return utils.FanIn(doneCh, utils.FanOut(10, func() chan error {
		result := make(chan error)

		wg.Add(1)
		go func() {
			defer wg.Done()

			defer close(result)

			for URL := range inputCh {
				logger.Log.Infof("About to delete URL '%s' by user %s", URL, userID)
				err := s.Container.Storage.DeleteByUser(ctx, userID, URL)
				if err != nil {
					logger.Log.Warnf("Could not delete URL '%s' for user '%s': %s", URL, userID, err.Error())
				}
				logger.Log.Infof("URL '%s' by user %s deleted", URL, userID)
				select {
				case <-doneCh:
					return
				case result <- err:
				}
			}
		}()

		return result
	}))
}
