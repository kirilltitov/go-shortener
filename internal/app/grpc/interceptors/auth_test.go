package interceptors

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/kirilltitov/go-shortener/internal/app"
	"github.com/kirilltitov/go-shortener/internal/app/grpc/gen"
	"github.com/kirilltitov/go-shortener/internal/app/grpc/testhelpers"
)

func Test_authenticate(t *testing.T) {
	authenticatedContext, userID := testhelpers.GetValidUserAndToken()

	emptyInfo := &grpc.UnaryServerInfo{}

	type args struct {
		ctx  context.Context
		info *grpc.UnaryServerInfo
	}
	type want struct {
		userIDIssued bool
		userID       *uuid.UUID
		outputToken  string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "No context",
			args: args{
				ctx:  context.Background(),
				info: emptyInfo,
			},
			want: want{
				userID:      nil,
				outputToken: "",
			},
		},
		{
			name: "No token, don't create new",
			args: args{
				ctx:  testhelpers.NewContextWithValue("token", ""),
				info: emptyInfo,
			},
			want: want{
				userID:      nil,
				outputToken: "",
			},
		},
		{
			name: "No token, create new",
			args: args{
				ctx: testhelpers.NewContextWithValue("", ""),
				info: &grpc.UnaryServerInfo{
					FullMethod: gen.Shortener_GetURL_FullMethodName,
				},
			},
			want: want{
				userIDIssued: true,
			},
		},
		{
			name: "No token, strict method, don't create new",
			args: args{
				ctx: testhelpers.NewContextWithValue("", ""),
				info: &grpc.UnaryServerInfo{
					FullMethod: gen.Shortener_GetUserURLs_FullMethodName,
				},
			},
			want: want{
				userID:      nil,
				outputToken: "",
			},
		},
		{
			name: "Invalid token",
			args: args{
				ctx:  testhelpers.NewContextWithValue("token", "I.AM.INVALID"),
				info: emptyInfo,
			},
			want: want{
				userID:      nil,
				outputToken: "",
			},
		},
		{
			name: "Valid token",
			args: args{
				ctx:  authenticatedContext,
				info: emptyInfo,
			},
			want: want{
				userIDIssued: true,
				userID:       &userID,
			},
		},
	}
	for _, tt := range tests {
		result := authenticate(tt.args.ctx, tt.args.info)

		resultUserID, ok := result.Value(app.CtxUserIDKey{}).(uuid.UUID)
		if tt.want.userIDIssued {
			require.NotNil(t, resultUserID)
			if !ok {
				fmt.Printf("what")
			}
			require.True(t, ok)
		}

		if tt.want.userID != nil {
			require.NotNil(t, resultUserID)
			require.True(t, ok)
			require.Equal(t, tt.want.userID, &resultUserID)
		}

		if tt.want.outputToken != "" {
			md, mdOk := metadata.FromIncomingContext(result)
			require.True(t, mdOk)
			tokenSlice := md.Get("token")
			require.Len(t, tokenSlice, 1)
			require.Equal(t, tt.want.outputToken, tokenSlice[0])
		}
	}
}
