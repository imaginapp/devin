package invite

import (
	"database/sql"

	"github.com/imaginapp/proto/go/gen/imagin/external/message/v1"
)

const OpenInviteCodeGroupID = "open"
const OpenInviteMinimumBalance = 30

type InviteStatus string

const InviteStatusAny InviteStatus = "any"
const InviteStatusUsed InviteStatus = "used"
const InviteStatusUnused InviteStatus = "unused"
const InviteStatusLocked InviteStatus = "locked"

func (s InviteStatus) String() string {
	return string(s)
}

type Invite struct {
	Code            string
	AccountID       sql.NullString
	GroupID         sql.NullString
	TransactionHash sql.NullString
	CreatedAt       string
	UpdatedAt       sql.NullString
}

func (i *Invite) IsRedeemed() bool {
	if i.IsLocked() {
		return false
	}
	return i.AccountID.String != ""
}

func (i *Invite) Status() InviteStatus {
	if i.IsLocked() {
		return InviteStatusLocked
	}
	if i.IsRedeemed() {
		return InviteStatusUsed
	}

	return InviteStatusUnused
}

func (i *Invite) IsRedeemedBy(accountID string) bool {
	return i.AccountID.String == accountID
}

func (i *Invite) IsLocked() bool {
	// Check both conditions in a more readable way
	isOpenInviteGroup := i.GroupID.String == OpenInviteCodeGroupID
	isLocked := i.AccountID.String == "locked"

	return isOpenInviteGroup && isLocked
}

func (i *Invite) Proto() *message.Invite {
	invitepb := &message.Invite{
		Code:            i.Code,
		AccountId:       i.AccountID.String,
		TransactionHash: i.TransactionHash.String,
		CreatedAt:       stringToTimestamp(i.CreatedAt),
		UpdatedAt:       stringToTimestamp(i.UpdatedAt.String),
	}

	return invitepb
}
