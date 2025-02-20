package testhelpers

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/kirilltitov/go-shortener/internal/app"
	"github.com/kirilltitov/go-shortener/internal/app/auth"
)

func GetValidUserAndToken() (context.Context, uuid.UUID) {
	userID, _ := uuid.NewV6()
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		auth.Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: userID.String(),
			},
		},
	)

	tokenString, _ := token.SignedString([]byte(auth.JWTSecret))

	ctx := NewContextWithValue("token", tokenString)
	ctx = context.WithValue(ctx, app.CtxUserIDKey{}, userID)

	return ctx, userID
}
