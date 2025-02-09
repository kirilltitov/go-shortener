package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/kirilltitov/go-shortener/internal/config"
)

// JWTSecret является секретом для подписи авторизационного JWT.
const JWTSecret = "hesoyam"

// Claims является объектом для декодинга переданного JWT.
type Claims struct {
	jwt.RegisteredClaims
}

type Auth struct {
	config config.Config
}

func New(cfg config.Config) *Auth {
	return &Auth{
		config: cfg,
	}
}

func IssueNewToken(userID uuid.UUID) (*string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   userID.String(),
				ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(1, 0, 0)),
			},
		},
	)

	tokenString, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func ParseTokenString(input string, claims *Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(input, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})
}
