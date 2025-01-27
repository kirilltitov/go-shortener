package app

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/kirilltitov/go-shortener/internal/logger"
)

// Claims является объектом для декодинга переданного JWT.
type Claims struct {
	jwt.RegisteredClaims
}

const (
	// JWTCookieName является ключом для названия куки, в которой будет храниться авторизационный JWT.
	JWTCookieName = "access_token"

	// JWTSecret является секретом для подписи авторизационного JWT.
	JWTSecret = "hesoyam"
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

			cookie, err2 := newCookie(userID)
			if err2 != nil {
				return nil, err2
			}

			http.SetCookie(w, cookie)

			logger.Log.Infof("Auth cookie set for user %s", userID.String())

			return &userID, nil
		} else {
			return nil, err
		}
	}

	tokenString := cookie.Value
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})

	if err != nil || !token.Valid {
		if force {
			return nil, nil
		}
		logger.Log.Infof("Could not parse auth cookie or JWT not valid, will issue new")

		userID, err := uuid.NewV6()
		if err != nil {
			return nil, err
		}

		cookie, err := newCookie(userID)
		if err != nil {
			return nil, err
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

func newCookie(userID uuid.UUID) (*http.Cookie, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: userID.String(),
			},
		},
	)

	tokenString, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		return nil, err
	}

	cookie := http.Cookie{
		Name:    JWTCookieName,
		Value:   tokenString,
		Expires: time.Now().AddDate(1, 0, 0),
	}

	return &cookie, nil
}
