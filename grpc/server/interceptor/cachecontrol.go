package interceptor

import (
	"context"

	"google.golang.org/grpc"
)

func CacheControl() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		cacheControl := getMetadataValue(ctx, MetadataCacheControlKey)
		ctx = context.WithValue(ctx, ContextCacheControlKey, cacheControl)
		return handler(ctx, req)
	}
}
