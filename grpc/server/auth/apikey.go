package auth

import (
	"context"
	"errors"
	"os"

	"google.golang.org/grpc/metadata"
)

func CheckAPIKey(ctx context.Context) error {
	apiKey := os.Getenv("API_KEY")

	if apiKey == "" {
		return errors.New("API_KEY missing from env")
	}

	reqAPIKey, err := getAPIKey(ctx)
	if err != nil {
		return err
	}

	if reqAPIKey != apiKey {
		return ErrInvalidAPIKey
	}

	return nil
}

const metadataAPIKEYKey = "x-api-key"

func getAPIKey(ctx context.Context) (string, error) {
	var md metadata.MD
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ErrAPIKeyNotFound
	}
	apiKeys := md.Get(metadataAPIKEYKey)
	if len(apiKeys) != 1 {
		return "", ErrAPIKeyNotFound
	}

	return apiKeys[0], nil
}
