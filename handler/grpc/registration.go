package handler

import (
	"context"
	"time"

	"github.com/imaginapp/devin"
	"github.com/imaginapp/proto/go/gen/imagin/external/service/v1"
	"github.com/rs/zerolog"
)

func (h *Handler) ValidateInviteCode(ctx context.Context, in *service.ValidateInviteCodeRequest) (*service.ValidateInviteCodeResponse, error) {
	if err := checkAPIKey(ctx); err != nil {
		return nil, err
	}
	logger := zerolog.Ctx(ctx)
	logger.Debug().Str("code", in.Code).Msg("ValidateInviteCode")
	defer devin.RandomDelay()

	var cacheResponse service.ValidateInviteCodeResponse
	cacheKey := cacheKey(in)
	if err := h.cacheGet(ctx, cacheKey, &cacheResponse); err == nil {
		logger.Debug().Bool("isCached", true).Msg("got cached invite code")
		return &cacheResponse, nil
	} else {
		logger.Debug().Err(err).Msg("account not found in cache")
	}
	// check the invite code
	isCodeValid := h.s.Invites.IsCodeValid(ctx, in.Code)

	response := &service.ValidateInviteCodeResponse{IsValid: isCodeValid}
	h.cacheSet(ctx, cacheKey, response, time.Second*10)
	return response, nil
}
