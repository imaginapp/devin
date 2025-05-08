package invite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/imaginapp/proto/go/gen/imagin/external/message/v1"
	"github.com/rs/zerolog/log"
)

func (c *Client) IsCodeValid(ctx context.Context, code string) bool {
	if len(code) != codeLength || !isAlphanumeric(code) {
		return false
	}
	inviteData, err := c.getCode(ctx, code)
	if err != nil {
		return false
	}
	// found invite but already used.
	if inviteData.Status() != InviteStatusUnused {
		return false
	}

	return true
}

var getCodeQuery = `SELECT code, group_id, account_id, transaction_hash, updated_at, created_at FROM invites WHERE upper(code) = upper(?)`

func (c *Client) getCode(ctx context.Context, code string) (*Invite, error) {
	var invite Invite
	err := c.db.QueryRowContext(ctx, getCodeQuery, code).Scan(
		&invite.Code,
		&invite.GroupID,
		&invite.AccountID,
		&invite.TransactionHash,
		&invite.UpdatedAt,
		&invite.CreatedAt,
	)
	if err != nil {
		return nil, ErrInvalidCode
	}

	return &invite, nil
}

var redeemCodeQuery = `UPDATE invites SET account_id = ?, transaction_hash = ?, updated_at = ? WHERE upper(code) = upper(?)`

func (c *Client) RedeemCode(ctx context.Context, code, accountID, transactionHash string) error {
	inviteCode, err := c.getCode(ctx, code)
	if err != nil {
		// code not found
		return ErrInvalidCode
	}

	if inviteCode.IsRedeemedBy(accountID) {
		// code already verified by account
		// this is ok, return nil
		return nil
	}

	if inviteCode.Status() == InviteStatusUsed {
		// code already used by another account
		return ErrInvalidCode
	}

	_, err = c.db.ExecContext(ctx, redeemCodeQuery, accountID, transactionHash, timeNowString(), code)
	if err != nil {
		log.Error().Err(err).Msg("failed to update invite")
		// log error and continue so we dont break the flow
	}

	return nil
}

func (c *Client) IsCodeRedeemable(ctx context.Context, code string) bool {
	inviteData, err := c.getCode(ctx, code)
	if err != nil {
		return false
	}
	// found invite but already used.
	if inviteData.Status() == InviteStatusUsed {
		return false
	}

	return true
}

var getNextAvailableCodeQuery = `SELECT code FROM invites WHERE (account_id IS NULL OR account_id = '') AND group_id = ? ORDER BY created_at ASC LIMIT 1`
var lockNextAvailableCodeQuery = `UPDATE invites SET account_id = ?, updated_at = ? WHERE code = ? AND (account_id IS NULL OR account_id = '')`

func (c *Client) GetAndLockNextInviteCode(ctx context.Context, groupID string) (string, error) {
	// Unlock any locked codes that are older than 20 minutes
	if err := c.unlockLockedCodes(ctx); err != nil {
		log.Logger.Error().Err(err).Msg("failed to unlock locked codes")
	}

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var code string
	err = tx.QueryRowContext(ctx, getNextAvailableCodeQuery, groupID).Scan(&code)
	if err != nil {
		log.Logger.Error().Err(err).Msg("failed to get next available invite code")
		return "", ErrGetNextInviteCode
	}

	result, err := tx.ExecContext(ctx, lockNextAvailableCodeQuery, "locked", timeNowString(), code)
	if err != nil {
		log.Logger.Error().Err(err).Msg("failed to lock invite code")
		return "", ErrGetNextInviteCode
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Logger.Error().Err(err).Msg("failed to get rows affected")
		return "", ErrGetNextInviteCode
	}
	if rowsAffected == 0 {
		log.Logger.Warn().Msg("invite was taken by another process")
		return "", ErrGetNextInviteCode
	}

	if err := tx.Commit(); err != nil {
		log.Logger.Error().Err(err).Msg("failed to commit transaction")
		return "", ErrGetNextInviteCode
	}

	return code, nil
}

var unlockLockedCodesQuery = `UPDATE invites SET account_id = '', updated_at = ? WHERE account_id = ? AND updated_at < ?`

func (c *Client) unlockLockedCodes(ctx context.Context) error {
	twentyMinutesAgo := time.Now().Add(-20 * time.Minute)

	result, err := c.db.ExecContext(ctx, unlockLockedCodesQuery, timeNowString(), "locked", formatTime(twentyMinutesAgo))
	if err != nil {
		return fmt.Errorf("failed to unlock locked codes: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		log.Debug().
			Int64("unlocked_codes", rowsAffected).
			Msg("successfully unlocked codes")
	}

	return nil
}

var insertCodeQuery = `INSERT INTO invites (code, group_id, created_at) VALUES (?, ?, ?)`

func (c *Client) GenerateCodes(ctx context.Context, count uint64, groupID string) ([]string, error) {
	currentCount := c.totalInviteCount(ctx)

	codes := []string{}

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	groupIDNullString := sql.NullString{String: groupID, Valid: groupID != ""}
	for i := 0; i < int(count); i++ {
		// get code
		code, err := c.generateCode(i + int(currentCount))
		if err != nil {
			return nil, err
		}
		codes = append(codes, code)
		// insert code
		if _, err := tx.Exec(insertCodeQuery, code, groupIDNullString, timeNowString()); err != nil {
			return nil, ErrInsertFailed
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return codes, nil
}

var totalCodeCountQuery = `SELECT COUNT(1) FROM invites`

func (c *Client) totalInviteCount(ctx context.Context) int64 {
	var count int64
	err := c.db.QueryRowContext(ctx, totalCodeCountQuery).Scan(&count)
	if err != nil {
		return 0
	}
	return count
}

var getInvitesQuery = `SELECT code, group_id, account_id, transaction_hash, updated_at, created_at FROM invites`
var getInvitesWhereUnused = `(account_id == '' OR account_id IS NULL)`
var getInvitesWhereUsed = `(account_id != '' AND account_id IS NOT NULL)`
var getInvitesWhereGroup = `group_id = ?`
var getInvitesWhereNoGroup = `(group_id == '' OR group_id IS NULL)`

func (c *Client) GetInvites(ctx context.Context, status string, limit, offset int, groupID string) ([]*message.Invite, error) {
	var args []any
	var where []string
	if status == InviteStatusUnused.String() {
		where = append(where, getInvitesWhereUnused)
	} else if status == InviteStatusUsed.String() {
		where = append(where, getInvitesWhereUsed)
	}
	if groupID != "" {
		args = append(args, groupID)
		where = append(where, getInvitesWhereGroup)
	} else {
		where = append(where, getInvitesWhereNoGroup)
	}

	query := getInvitesQuery
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}
	args = append(args, limit, offset)
	query += " ORDER BY created_at asc LIMIT ? OFFSET ?"

	var invites []*Invite
	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var invite Invite
		if err := rows.Scan(
			&invite.Code,
			&invite.GroupID,
			&invite.AccountID,
			&invite.TransactionHash,
			&invite.UpdatedAt,
			&invite.CreatedAt,
		); err != nil {
			return nil, err
		}
		invites = append(invites, &invite)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var protoInvites []*message.Invite
	for _, invite := range invites {
		protoInvites = append(protoInvites, invite.Proto())
	}

	return protoInvites, nil
}
