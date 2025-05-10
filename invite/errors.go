package invite

import "errors"

var ErrInvalidCode = errors.New("code is invalid")
var ErrRedeemCodeFailed = errors.New("failed to redeem code")
var ErrInsertFailed = errors.New("failed to add code")
var ErrInvalidCodeLength = errors.New("invalid code length")
var ErrGetNextInviteCode = errors.New("faild to get invite code")
