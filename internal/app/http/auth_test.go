package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/kirilltitov/go-shortener/internal/app/auth"
	"github.com/kirilltitov/go-shortener/internal/container"
	"github.com/kirilltitov/go-shortener/internal/shortener"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kirilltitov/go-shortener/internal/config"
)

func TestApplication_authenticate(t *testing.T) {
	cfg := config.NewWithoutParsing()
	cfg.DatabaseDSN = ""
	cfg.FileStoragePath = ""
	cnt, err := container.New(context.Background(), cfg)
	require.NoError(t, err)
	service := shortener.New(cfg, cnt)
	a := New(service, &sync.WaitGroup{})

	type want struct {
		userID        *uuid.UUID
		cookieWritten bool
	}
	type testCase struct {
		name  string
		input *string
		want  want
	}

	getJWT := func(userID string) *string {
		token := jwt.NewWithClaims(
			jwt.SigningMethodHS256,
			auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject: userID,
				},
			},
		)

		tokenString, _ := token.SignedString([]byte(auth.JWTSecret))

		return &tokenString
	}

	userID, _ := uuid.Parse("0f227b5e-81a6-11ee-b962-0242ac120abd")
	invalidJWT := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIwZjIyN2I1ZS04MWE2LTExZWUtYjk2Mi0wMjQyYWMxMjBhYmQifQ.I AM INVALID"

	tests := []testCase{
		{
			name:  "Positive",
			input: getJWT(userID.String()),
			want: want{
				userID:        &userID,
				cookieWritten: false,
			},
		},
		{
			name:  "Negative (no cookie)",
			input: nil,
			want: want{
				userID:        nil,
				cookieWritten: true,
			},
		},
		{
			name:  "Negative (invalid JWT)",
			input: &invalidJWT,
			want: want{
				userID:        nil,
				cookieWritten: true,
			},
		},
		{
			name:  "Unauthorized",
			input: getJWT(""),
			want: want{
				userID:        nil,
				cookieWritten: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			if tt.input != nil {
				r.AddCookie(&http.Cookie{
					Name:  JWTCookieName,
					Value: *tt.input,
				})
			}

			userID, err := a.authenticate(r, w, false)
			require.NoError(t, err)

			if tt.want.userID != nil {
				assert.Equal(t, tt.want.userID, userID)
			}

			if tt.want.cookieWritten == true {
				resp := w.Result()
				defer resp.Body.Close()
				_, ok := resp.Header["Set-Cookie"]
				assert.True(t, ok)
			}
		})
	}
}
