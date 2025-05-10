package grpcserver

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	// Installing the gzip encoding registers it as an available compressor.
	// gRPC will automatically negotiate and use gzip if the client supports it.
	_ "google.golang.org/grpc/encoding/gzip"
)

type Handler interface {
	GRPCHandler(s *grpc.Server)
}

type ServerConfig struct {
	Address        string
	Hostname       string
	MaxRecvMsgSize int // Control message size
	MaxSendMsgSize int // Control message size
}

type Server struct {
	address string
	server  *grpc.Server
}

func (s *Server) Close() error {
	if s.server != nil {
		s.server.GracefulStop()
	}
	return nil
}

func New(config ServerConfig, handler Handler, interceptors ...grpc.UnaryServerInterceptor) *Server {
	serverOptions := []grpc.ServerOption{}
	if config.MaxRecvMsgSize > 0 {
		serverOptions = append(serverOptions, grpc.MaxRecvMsgSize(config.MaxRecvMsgSize))
	}
	if config.MaxSendMsgSize > 0 {
		serverOptions = append(serverOptions, grpc.MaxRecvMsgSize(config.MaxSendMsgSize))
	}
	if len(interceptors) > 0 {
		serverOptions = append(serverOptions, grpc.ChainUnaryInterceptor(interceptors...))
	}

	s := grpc.NewServer(serverOptions...)

	// add health check
	healthService := NewHealthChecker()
	grpc_health_v1.RegisterHealthServer(s, healthService)

	// set up the handler
	handler.GRPCHandler(s)

	return &Server{
		address: config.Address,
		server:  s,
	}
}

func (s *Server) Start() error {
	// start the GRPC server
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to listen to gRPC: %w", err)
	}
	defer lis.Close()

	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve gRPC: %w", err)
	}

	return nil
}
