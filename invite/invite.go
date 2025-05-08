package invite

import (
	"database/sql"
	"time"

	"github.com/speps/go-hashids/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Client struct {
	db     *sql.DB
	hashID *hashids.HashID
}

func New(db *sql.DB, salt string) (*Client, error) {
	hashID, err := initHashID(salt)
	if err != nil {
		return nil, err
	}

	// setup the tables
	if err := setup(db); err != nil {
		return nil, err
	}

	return &Client{db: db, hashID: hashID}, nil
}

func timeNowString() string {
	return formatTime(time.Now())
}

func formatTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

func stringToTimestamp(s string) *timestamppb.Timestamp {
	if s == "" {
		return nil
	}
	createdAt, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil
	}
	if createdAt.IsZero() {
		return nil
	}

	return timestamppb.New(createdAt)
}
