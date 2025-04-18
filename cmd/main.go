package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/imaginapp/devin"
	grpcserver "github.com/imaginapp/devin/grpc/server"
	"github.com/imaginapp/devin/grpc/server/interceptor"
	ghandler "github.com/imaginapp/devin/handler/grpc"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

	services := &devin.Services{}

	grpcConfig := grpcserver.ServerConfig{
		Address:        fmt.Sprintf(":%s", getEnvOrDefaultString("APP_PORT", defaultGrpcPort)),
		MaxRecvMsgSize: 1024 * 1024,
		MaxSendMsgSize: 1024 * 1024,
	}
	gs := grpcserver.New(
		grpcConfig,
		ghandler.New(services, withReflection),
		interceptor.RequestID(),
		interceptor.Logging(log.Logger),
		interceptor.APIKey(os.Getenv("API_KEY")),
	)

	log.Info().Str("address", grpcConfig.Address).Msg("staring grpc server")
	if err := gs.Start(); err != nil {
		log.Fatal().Err(err).Msg("failed to start grpc server")
	}
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
