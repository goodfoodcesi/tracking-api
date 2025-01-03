package grpc

import (
	"context"
	"github.com/goodfoodcesi/tracking-api/pkg/config"
	"github.com/goodfoodcesi/tracking-api/pkg/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func JWTInterceptor(cfg config.Config) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
		}

		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
		}

		tokenStr := token.ExtractTokenFromHeader(authHeader[0])
		if err := token.ValidateToken(tokenStr, cfg); err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		//TODO A FIXER
		//userID, err := token.ExtractTokenID(tokenStr, cfg)
		//if err != nil {
		//	return nil, status.Error(codes.Internal, "failed to process token")
		//}
		//ctx = context.WithValue(ctx, "user_id", userID)

		return handler(ctx, req)
	}
}
