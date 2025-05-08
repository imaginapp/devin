package devin

import "context"

type Services struct {
	Invites Invites
}

type Invites interface {
	IsCodeValid(ctx context.Context, code string) bool
	RedeemCode(ctx context.Context, code, accountID, transactionHash string) error
	IsCodeRedeemable(ctx context.Context, code string) bool
	GetAndLockNextInviteCode(ctx context.Context, groupID string) (string, error)
}
