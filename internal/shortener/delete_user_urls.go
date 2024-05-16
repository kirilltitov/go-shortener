package shortener

import (
	"context"

	"github.com/google/uuid"

	"github.com/kirilltitov/go-shortener/internal/logger"
	"github.com/kirilltitov/go-shortener/internal/utils"
)

func (s *Shortener) DeleteUserURLs(ctx context.Context, doneCh chan struct{}, userID uuid.UUID, URLs []string) chan error {
	inputCh := utils.Generator(doneCh, URLs)

	return utils.FanIn(doneCh, utils.FanOut(10, func() chan error {
		result := make(chan error)

		go func() {
			defer close(result)

			for URL := range inputCh {
				logger.Log.Infof("About to delete URL '%s' by user %s", URL, userID)
				err := s.container.Storage.DeleteByUser(ctx, userID, URL)
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
