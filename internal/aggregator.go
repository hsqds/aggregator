package internal

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

// Aggregator represents
type Aggregator struct {
	logger     *log.Logger
	loader     *Loader
	parser     *Parser
	repository *PostRepository
	pgmtx      *PGMutex
}

// NewRSSReader
func NewAggregator(loader *Loader, parser *Parser, repository *PostRepository, pgmutex *PGMutex) *Aggregator {
	return &Aggregator{
		log.New(os.Stderr, aggregatorLogPrefix, loggerFlags),
		loader,
		parser,
		repository,
		pgmutex,
	}
}

// Update
func (a *Aggregator) Update(ctx context.Context) error {
	a.logger.Println("starting update")

	const (
		loadTimeout = 2
		feedsChSize = 5
		itemsChSize = 5 * 30
	)

	if err := a.pgmtx.Lock(ctx); err != nil {
		return fmt.Errorf("could not set lock: %w", err)
	}

	loadCtx, cancel := context.WithTimeout(ctx, loadTimeout*time.Second)
	defer cancel()

	feeds := make(chan Feed, feedsChSize)
	posts := make(chan Post, itemsChSize)
	done := make(chan struct{})

	go a.loader.Load(loadCtx, feeds)
	go a.parser.Parse(ctx, feeds, posts)
	go a.repository.Save(ctx, posts, done)

	select {
	case <-ctx.Done():
		cancel()

		return errors.New("update timeout exceeded")
	case <-done:
		return nil
	}
}
