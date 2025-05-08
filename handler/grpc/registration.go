package handler

import (
	"context"

	"github.com/imaginapp/devin"
	"github.com/imaginapp/proto/go/gen/imagin/external/service/v1"
	"github.com/rs/zerolog"
)

func (h *Handler) ValidateInviteCode(ctx context.Context, in *service.ValidateInviteCodeRequest) (*service.ValidateInviteCodeResponse, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Str("code", in.Code).Msg("CheckInviteCode")
	devin.RandomDelay()

	// check the invite code
	isCodeValid := h.s.Invites.IsCodeValid(ctx, in.Code)

	return &service.ValidateInviteCodeResponse{IsValid: isCodeValid}, nil
}
