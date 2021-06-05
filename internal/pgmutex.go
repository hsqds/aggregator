package internal

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
)

const updateLockID = 55

// PGMutex represents
type PGMutex struct {
	logger *log.Logger
	dbconn *pgx.Conn
}

// NewPGMutex
func NewPGMutex(dbconn *pgx.Conn) (*PGMutex, error) {
	mtx := &PGMutex{
		log.New(os.Stderr, pgMutextLogPrefix, loggerFlags),
		dbconn,
	}

	return mtx, nil
}

// lock
func (mtx *PGMutex) Lock(ctx context.Context) error {
	row := mtx.dbconn.QueryRow(ctx, "SELECT pg_try_advisory_lock($1)", updateLockID)

	var locked bool

	if err := row.Scan(&locked); err != nil {
		return fmt.Errorf("could not get advisory lock result: %w", err)
	}

	if !locked {
		return fmt.Errorf("could not set advisory lock (id: %d): ", updateLockID)
	}

	mtx.logger.Printf("successfully locked (id: %d)", updateLockID)

	return nil
}

// unlock
func (mtx *PGMutex) Unlock(ctx context.Context) {
	_, err := mtx.dbconn.Query(ctx, "SELECT pg_advisory_unlock($1)", updateLockID)
	if err != nil {
		mtx.logger.Printf("could not unset advisory lock (id: %d): %s", updateLockID, err)

		return
	}

	mtx.logger.Printf("successfully unlocked (id: %d)", updateLockID)
}
