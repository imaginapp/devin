package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/imaginapp/devin"
	grpcserver "github.com/imaginapp/devin/grpc/server"
	"github.com/imaginapp/devin/grpc/server/interceptor"
	ghandler "github.com/imaginapp/devin/handler/grpc"
	"github.com/imaginapp/devin/invite"
	"github.com/imaginapp/devin/sqlite"
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
