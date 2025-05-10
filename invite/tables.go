package invite

import (
	"database/sql"
)

var inviteTable = `CREATE TABLE IF NOT EXISTS invites (
    code TEXT NOT NULL PRIMARY KEY,
    group_id TEXT DEFAULT NULL,
    account_id TEXT DEFAULT NULL,
    transaction_hash TEXT DEFAULT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT DEFAULT NULL
);`

var inviteGroupIdx = `CREATE INDEX IF NOT EXISTS invite_group_id_idx ON invites ("group_id");`

func setup(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(inviteTable); err != nil {
		return err
	}
	if _, err := tx.Exec(inviteGroupIdx); err != nil {
		return err
	}

	return tx.Commit()
}
