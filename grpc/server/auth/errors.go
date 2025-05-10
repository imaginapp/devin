package auth

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrAPIKeyNotFound apiKey was not found in request
var ErrAPIKeyNotFound = status.Errorf(codes.Unauthenticated, "apikey not found in request")

// ErrTokenInvalid was unable to ve verified
var ErrTokenInvalid = status.Errorf(codes.Unauthenticated, "token is invalid")

// ErrInvalidAPIKey api key is invalid
var ErrInvalidAPIKey = status.Errorf(codes.Unauthenticated, "invalid key")

// ErrMessageInvalid failed to read the protobuf message
var ErrMessageInvalid = status.Errorf(codes.Internal, "failed to read message")
