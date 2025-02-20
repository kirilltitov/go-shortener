package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kirilltitov/go-shortener/internal/app/auth"

	"github.com/kirilltitov/go-shortener/internal/logger"
)

// Конфигурация JWT.
const (
	// JWTCookieName является ключом для названия куки, в которой будет храниться авторизационный JWT.
	JWTCookieName = "access_token"
)

func (a *Application) authenticate(r *http.Request, w http.ResponseWriter, force bool) (*uuid.UUID, error) {
	logger.Log.Infof("Will try to authenticate user")

	cookie, err := r.Cookie(JWTCookieName)
	if err != nil {
		if force {
			return nil, nil
		}
		if errors.Is(err, http.ErrNoCookie) {
			logger.Log.Infof("Auth cookie not found, will issue new")

			userID, err2 := uuid.NewV6()
			if err2 != nil {
				return nil, err2
			}

			newCookie, err2 := a.newCookie(userID)
			if err2 != nil {
				return nil, err2
			}

			http.SetCookie(w, newCookie)

			logger.Log.Infof("Auth cookie set for user %s", userID.String())

			return &userID, nil
		} else {
			return nil, err
		}
	}

	claims := &auth.Claims{}
	token, err := auth.ParseTokenString(cookie.Value, claims)

	if err != nil || !token.Valid {
		if force {
			return nil, nil
		}
		logger.Log.WithError(err).Info("Could not parse auth cookie or JWT not valid, will issue new")

		userID, err2 := uuid.NewV6()
		if err2 != nil {
			return nil, err2
		}

		cookie, err2 := a.newCookie(userID)
		if err2 != nil {
			return nil, err2
		}

		http.SetCookie(w, cookie)

		logger.Log.Infof("Auth cookie set for user %s", userID.String())

		return &userID, nil
	}
	if claims.Subject == "" {
		return nil, nil
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, err
	}

	logger.Log.Infof("Authenticated user %s by JWT cookie", userID.String())

	return &userID, nil
}

func (a *Application) newCookie(userID uuid.UUID) (*http.Cookie, error) {
	tokenString, err := auth.IssueNewToken(userID)
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:    JWTCookieName,
		Value:   *tokenString,
		Expires: time.Now().AddDate(1, 0, 0),
	}, nil
}
