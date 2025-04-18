package interceptor

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func RequestID() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		requestID := getMetadataValue(ctx, MetadataRequestIDKey)
		if requestID == "" {
			if newUUID, err := uuid.NewRandom(); err == nil {
				requestID = newUUID.String()
			}
		}
		ctx = context.WithValue(ctx, ContextRequestIDKey, requestID)
		return handler(ctx, req)
	}
}
