package internal

import (
	"context"
	"fmt"
	"time"

	pgx "github.com/jackc/pgx/v4"
)

const (
	pgConnectTimeout    = 3 * time.Second
	pgDisconnectTimeout = 3 * time.Second
)

// NewPGXConn
func NewPGXConn(cfg *DBConfig) (*pgx.Conn, func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), pgConnectTimeout)
	defer cancel()

	conn, err := pgx.Connect(ctx, cfg.String())
	if err != nil {
		return nil, func() {}, fmt.Errorf("could not connect to postgres: %w", err)
	}

	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), pgDisconnectTimeout)
		defer cancel()

		conn.Close(ctx)
	}

	return conn, cleanup, nil
}
