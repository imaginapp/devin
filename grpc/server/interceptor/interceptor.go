package interceptor

import (
	"context"

	"google.golang.org/grpc/metadata"
)

const MetadataCacheControlKey = "cache-control"
const MetadataAPIKeyKey = "x-api-key"
const MetadataRequestIDKey = "x-request-id"
const MetadataAccountIDKey = "x-account-id"

type contextKey string

var (
	ContextCacheControlKey = contextKey(MetadataCacheControlKey)
	ContextRequestIDKey    = contextKey(MetadataRequestIDKey)
	ContextAccountIDKey    = contextKey(MetadataAccountIDKey)
)

func getMetadataValue(ctx context.Context, key string) string {
	var md metadata.MD
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	values := md.Get(key)
	if len(values) != 1 {
		return ""
	}
	return values[0]
}
