//+build wireinject

package internal

import (
	"github.com/google/wire"
)

// BuildAggregator
func BuildAggregator(cfg *Config) (*Aggregator, func(), error) {
	wire.Build(
		GetDBConfig,
		GetFeedConfigs,
		NewAggregator,
		NewLoader,
		NewParser,
		NewPostRepository,
		NewPGXConn,
		NewPGMutex,
	)

	return &Aggregator{}, func() {}, nil
}

// BuildRepository
func BuildRepository(cfg *Config) (*PostRepository, func(), error) {
	wire.Build(
		GetDBConfig,
		NewPostRepository,
		NewPGXConn,
	)

	return &PostRepository{}, func() {}, nil
}
