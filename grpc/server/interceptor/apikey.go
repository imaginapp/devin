package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func APIKey(apiKey string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if apiKey == "" || apiKey != getMetadataValue(ctx, MetadataAPIKeyKey) {
			return nil, status.Error(codes.Unauthenticated, "Invalid API key")
		}
		return handler(ctx, req)
	}
}
