package handler

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/imaginapp/devin"
	"github.com/imaginapp/devin/grpc/server/auth"
	"github.com/imaginapp/devin/hash"
	"github.com/imaginapp/proto/go/gen/imagin/external/service/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type Cache interface {
	Get(ctx context.Context, key string, value any) error
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
}

type Handler struct {
	s              *devin.Services
	cache          Cache
	withReflection bool

	service.UnimplementedImaginServiceServer
}

// Caching Setup

func (h *Handler) SetCache(c Cache) {
	h.cache = c
}

func cacheKey(req proto.Message) string {
	bytes, err := proto.Marshal(req)
	if err != nil {
		log.Debug().Err(err).Msg("failed to create cache key")
		return ""
	}
	reqHash := hash.BytesHash(bytes).Base36()
	return fmt.Sprintf("%s:%s", req.ProtoReflect().Descriptor().Name(), reqHash)
}

func (h *Handler) cacheGet(ctx context.Context, key string, value any) error {
	if h.cache == nil {
		return errors.New("no cache configured")
	}
	return h.cache.Get(ctx, key, value)
}

func (h *Handler) cacheSet(ctx context.Context, key string, value any, ttl time.Duration) error {
	if h.cache == nil {
		return errors.New("no cache configured")
	}
	return h.cache.Set(ctx, key, value, ttl)
}

func New(services *devin.Services, withReflection bool) *Handler {
	return &Handler{
		s:              services,
		withReflection: withReflection,
	}
}

func (h *Handler) GRPCHandler(s *grpc.Server) {
	service.RegisterImaginServiceServer(s, h)
	if h.withReflection {
		reflection.Register(s)
	}
}

func checkAPIKey(ctx context.Context) error {
	if err := auth.CheckAPIKey(ctx); err != nil {
		return status.Error(codes.Unauthenticated, err.Error())
	}
	return nil
}

func (h *Handler) GetVersion(ctx context.Context, _ *service.GetVersionRequest) (*service.GetVersionResponse, error) {
	if err := checkAPIKey(ctx); err != nil {
		return nil, err
	}

	return &service.GetVersionResponse{Id: os.Getenv("VERSION")}, nil
}
