package interceptor

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	grpcstatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func GetLength(in interface{}) int {
	if pm, ok := in.(proto.Message); ok {
		b, _ := proto.Marshal(pm)
		return len(b)
	}
	return 0
}

// convert GRPC status code to HTTP status code
func getStatusCode(err error) int {
	if err == nil {
		return int(codes.OK)
	}

	statusCode := codes.Internal
	if grpcErr, ok := grpcstatus.FromError(err); ok {
		statusCode = grpcErr.Code()
	}

	return int(statusCode)
}

func Logging(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		var err error

		var requestID string
		if val, ok := ctx.Value(ContextRequestIDKey).(string); ok {
			requestID = val
		}
		var cacheControl string
		if val, ok := ctx.Value(ContextCacheControlKey).(string); ok {
			cacheControl = val
		}

		clientIP := "unknown"
		if peer, ok := peer.FromContext(ctx); ok {
			clientIP = peer.Addr.String()
		}

		// setup the access log
		accessLog := logger.With().
			Str("requestID", requestID).
			Str("clientIP", clientIP).
			Str("uri", info.FullMethod).
			Int("requestContentLength", GetLength(req)).
			Logger()

		ctx = log.With().Str("requestID", requestID).Logger().WithContext(ctx)

		// call the handler
		h, err := handler(ctx, req)

		statusCode := getStatusCode(err)

		timeTaken := time.Since(start)

		if err != nil {
			accessLog.Error().
				Int("status", statusCode).
				Dur("responseTime", timeTaken).
				Err(err).Msg("access log")

			return nil, err
		}

		if info.FullMethod == "/grpc.health.v1.Health/Check" {
			accessLog.Debug().
				Int("status", statusCode).
				Int("responseContentLength", GetLength(h)).
				Dur("responseTime", timeTaken).
				Msg("access log")

			return h, nil
		}

		accessLog.Info().
			Int("status", statusCode).
			Int("responseContentLength", GetLength(h)).
			Dur("responseTime", timeTaken).
			Str("cacheControl", cacheControl).
			Msg("access log")
		return h, nil
	}
}
