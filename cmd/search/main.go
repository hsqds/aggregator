package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/hsqds/rss-reader/internal"
)

func main() {
	var (
		configPath    string
		searchRequest string
		limit         int
	)

	flag.StringVar(&configPath, "c", "./config.json", "config file path")
	flag.StringVar(&searchRequest, "s", "", "search request")
	flag.IntVar(&limit, "l", 10, "results limit")
	flag.Parse()

	cfg, err := internal.InitConfig(configPath)
	if err != nil {
		log.Fatalf("could not initialize config: %s", err)
	}

	repo, cleanup, err := internal.BuildRepository(cfg)
	if err != nil {
		log.Printf("could not create repository: %s", err)
	}

	defer cleanup()

	const searchTimeout = time.Second

	ctx, cancel := context.WithTimeout(context.Background(), searchTimeout)
	defer cancel()

	pp, err := repo.Search(ctx, searchRequest, limit)
	if err != nil {
		log.Printf("could not find: %s", err)
	}

	if len(pp) == 0 {
		log.Println("nothing was found")
	}

	var shortDesc string

	const shortLen = 90

	for i := range pp {
		p := pp[i]

		if len(p.Description) > shortLen {
			shortDesc = p.Description[:shortLen]
		} else {
			shortDesc = p.Description
		}

		fmt.Printf("Title: %s\nLink: %s\nDescription: %s\n\n", p.Title, p.Link, shortDesc)
	}
}
