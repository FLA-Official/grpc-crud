package middleware

import (
	"context"
	"strings"

	"grpc-crud/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthInterceptor(secret string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		//  Public routes (no auth)
		public := map[string]bool{
			"/user.v1.UserService/Login":      true,
			"/user.v1.UserService/CreateUser": true,
		}

		if public[info.FullMethod] {
			return handler(ctx, req)
		}

		//  Read metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeader := md["authorization"]
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing token")
		}

		// Bearer token
		token := strings.TrimPrefix(authHeader[0], "Bearer ")

		// 🔍 Verify JWT
		payload, err := utils.VerifyJWT(secret, token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		//  Store in context
		ctx = context.WithValue(ctx, UserContextKey, payload)

		return handler(ctx, req)
	}
}
