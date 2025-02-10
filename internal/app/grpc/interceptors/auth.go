package interceptors

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/kirilltitov/go-shortener/internal/app"
	"github.com/kirilltitov/go-shortener/internal/app/auth"
	"github.com/kirilltitov/go-shortener/internal/app/grpc/gen"
	"github.com/kirilltitov/go-shortener/internal/logger"
)

const mdTokenKey = "token"

var strictAuthMethods = map[string]struct{}{
	gen.Shortener_GetUserURLs_FullMethodName:    {},
	gen.Shortener_DeleteUserURLs_FullMethodName: {},
}

func authenticate(ctx context.Context, info *grpc.UnaryServerInfo) context.Context {
	logger.Log.Infof("Will try to authenticate user")

	_, isStrictAuthMethod := strictAuthMethods[info.FullMethod]

	md, metadataOk := metadata.FromIncomingContext(ctx)
	if !metadataOk {
		return ctx
	}

	values := md.Get(mdTokenKey)
	var tokenString string
	if len(values) == 0 {
		tokenString = ""
	} else {
		tokenString = values[0]
	}

	if tokenString == "" {
		if isStrictAuthMethod {
			logger.Log.Infof("Auth token not found, method is strict, will not issue new token")
			return ctx
		}
		logger.Log.Infof("Auth token not found, will issue new")

		userID, err := uuid.NewV6()
		if err != nil {
			return ctx
		}

		newToken, err := auth.IssueNewToken(userID)
		if err != nil {
			return ctx
		}
		if err = grpc.SetHeader(ctx, metadata.New(map[string]string{mdTokenKey: *newToken})); err != nil {
			logger.Log.WithError(err).Error("Could not set token")
			return ctx
		}
		ctx = context.WithValue(ctx, app.CtxUserIDKey{}, userID)
		return ctx
	}

	claims := &auth.Claims{}
	token, err := auth.ParseTokenString(tokenString, claims)

	if err != nil {
		logger.Log.WithError(err).Error("Could not parse token string")
		return ctx
	}
	if !token.Valid {
		logger.Log.Infof("Auth token is not valid")
		return ctx
	}

	if claims.Subject == "" {
		return ctx
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		logger.Log.WithError(err).Error("Could not parse userID")
		return ctx
	}

	ctx = context.WithValue(ctx, app.CtxUserIDKey{}, userID)
	logger.Log.Infof("Authenticated user %s by JWT", userID.String())

	return ctx
}

func UnaryAuthInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	ctx = authenticate(ctx, info)

	return handler(ctx, req)
}
