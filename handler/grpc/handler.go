package handler

import (
	"context"
	"os"

	"github.com/imaginapp/devin"
	"github.com/imaginapp/proto/go/gen/imagin/external/service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Handler struct {
	s              *devin.Services
	withReflection bool

	service.UnimplementedImaginServiceServer
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

func (h *Handler) GetVersion(context.Context, *service.GetVersionRequest) (*service.GetVersionResponse, error) {
	return &service.GetVersionResponse{Id: os.Getenv("VERSION")}, nil
}
