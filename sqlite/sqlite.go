package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"

	"github.com/rs/zerolog/log"
)

type Client struct {
	db *sql.DB
}

func (c *Client) DB() *sql.DB {
	return c.db
}

func (c *Client) Close() error {
	if c.db != nil {
		c.db.Close()
	}
	return nil
}

func New(dbPath string) (*Client, error) {
	log.Debug().Str("dbPath", dbPath).Msg("connecting to sqlite db")
	db, err := sql.Open("libsql", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open db %s: %w", dbPath, err)
	}
	if db == nil {
		return nil, errors.New("db is nil")
	}

	err = onConnect(db)
	if err != nil {
		return nil, fmt.Errorf("faild running on connect: %w", err)
	}

	return &Client{
		db: db,
	}, nil
}

func onConnect(db *sql.DB) error {
	_, err := db.ExecContext(context.Background(), "PRAGMA foreign_keys = ON")
	if err != nil {
		return err
	}
	err = db.PingContext(context.Background())
	if err != nil {
		return err
	}

	return nil
}
