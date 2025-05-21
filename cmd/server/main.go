package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/imaginapp/devin"
	"github.com/imaginapp/devin/cache"
	grpcserver "github.com/imaginapp/devin/grpc/server"
	"github.com/imaginapp/devin/grpc/server/interceptor"
	ghandler "github.com/imaginapp/devin/handler/grpc"
	"github.com/imaginapp/devin/invite"
	"github.com/imaginapp/devin/lru"
	"github.com/imaginapp/devin/redis"
	"github.com/imaginapp/devin/sqlite"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

const defaultGrpcPort = "50051"

func main() {
	log.Logger = log.Level(zerolog.InfoLevel)
	var withReflection bool
	if os.Getenv("IS_LOCAL") == "true" {
		withReflection = true
		// override with pretty logging
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).
			Level(zerolog.DebugLevel)

		err := godotenv.Load()
		if err != nil {
			log.Fatal().Msg("error loading .env file")
		}
	}

	// Setup cache for gRPC service
	redisCacheClient, err := redis.NewFromEnv(0)
	if err != nil {
		log.Fatal().Err(err).Msg("failed start service")
	}
	lruClient, err := lru.New()
	if err != nil {
		log.Fatal().Err(err).Msg("failed start service")
	}
	grpcCacheClient := cache.New(lruClient, redisCacheClient)

	sqlPath := os.Getenv("DB_PATH")
	dbToken := os.Getenv("DB_TOKEN")
	if dbToken != "" {
		sqlPath = fmt.Sprintf("%s?authToken=%s", sqlPath, dbToken)
	}
	sqlClient, err := sqlite.New(sqlPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed start service")
	}
	defer sqlClient.Close()

	inviteClient, err := invite.New(sqlClient.DB(), "test123")
	if err != nil {
		log.Fatal().Err(err).Msg("failed start service")
	}

	services := &devin.Services{
		Invites: inviteClient,
	}

	errGrp, egCtx := errgroup.WithContext(context.Background())

	// Set up signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(egCtx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// start the gRPC server
	gh := ghandler.New(services, withReflection)
	gh.SetCache(grpcCacheClient)
	startGRPCServer(ctx, errGrp, gh)

	// keep the application running while servers are running
	if err := errGrp.Wait(); err != nil {
		log.Error().Err(err).Msg("error during shutdown")
	}
	log.Info().Msg("graceful shutdown completed")
}

func getEnvOrDefaultString(envKey string, defaultValue string) string {
	if v := os.Getenv(envKey); v != "" {
		return v
	}
	return defaultValue
}

func getEnvOrDefaultInt(envKey string, defaultValue int) int {
	if v := os.Getenv(envKey); v != "" {
		if intval, err := strconv.Atoi(v); err == nil {
			return intval
		}
	}
	return defaultValue
}

func startGRPCServer(ctx context.Context, errGrp *errgroup.Group, handler *ghandler.Handler) {
	grpcConfig := grpcserver.ServerConfig{
		Address:        fmt.Sprintf(":%s", getEnvOrDefaultString("APP_PORT", defaultGrpcPort)),
		MaxRecvMsgSize: 1024 * 1024,
		MaxSendMsgSize: 1024 * 1024,
	}
	gs := grpcserver.New(
		grpcConfig,
		handler,
		interceptor.RequestID(),
		interceptor.Logging(log.Logger),
	)

	log.Info().Str("address", grpcConfig.Address).Msg("starting grpc server")
	errGrp.Go(func() error {
		if err := gs.Start(); err != nil {
			log.Fatal().Err(err).Msg("failed to start grpc server")
		}
		return nil
	})

	shutdown := func(ctx context.Context) error {
		gs.Close()
		return nil
	}

	gracefulShutdown(ctx, errGrp, shutdown)
}

func gracefulShutdown(ctx context.Context, errGrp *errgroup.Group, shutdownFunc func(ctx context.Context) error) {
	errGrp.Go(func() error {
		<-ctx.Done() // Wait for cancellation signal
		log.Info().Msg("shutting down server")

		// Create a timeout context for shutdown
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := shutdownFunc(shutdownCtx); err != nil {
			return fmt.Errorf("server shutdown error: %w", err)
		}

		log.Info().Msg("server shutdown complete")
		return nil
	})
}
