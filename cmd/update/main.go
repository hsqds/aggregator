package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/hsqds/rss-reader/internal"
)

func main() {
	log.Println("Starting")

	var configPath string

	flag.StringVar(&configPath, "c", "./config.json", "config file path")
	flag.Parse()

	cfg, err := internal.InitConfig(configPath)
	if err != nil {
		log.Fatalf("could not initialize config: %s", err)
	}

	aggregator, cleanup, err := internal.BuildAggregator(cfg)
	if err != nil {
		log.Fatalf("could not initialize aggregator: %s", err)
	}

	defer cleanup()

	const updateTimeout = 10 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), updateTimeout)
	defer cancel()

	err = aggregator.Update(ctx)
	if err != nil {
		log.Fatalf("could not update: %s", err)
	}

	log.Println("Done")
}
