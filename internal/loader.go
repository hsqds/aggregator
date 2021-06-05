package internal

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

// Loader represents
type Loader struct {
	logger      *log.Logger
	feedConfigs []FeedConfig
}

// NewLoader
func NewLoader(feedConfigs []FeedConfig) *Loader {
	return &Loader{
		log.New(os.Stderr, loaderLogPrefix, loggerFlags),
		feedConfigs,
	}
}

// Load
func (l *Loader) Load(ctx context.Context, feeds chan<- Feed) {
	l.logger.Printf("start loading: %#v", l.feedConfigs)

	wg := &sync.WaitGroup{}
	wg.Add(len(l.feedConfigs))

	l.logger.Println("loading feeds...")
	for i := range l.feedConfigs {
		go l.loadFeed(ctx, l.feedConfigs[i], wg, feeds)
	}

	wg.Wait()

	l.logger.Println("Closing feeds channel")
	close(feeds)

	l.logger.Println("loading done")
}

// loadFeed
func (l *Loader) loadFeed(ctx context.Context, feedCfg FeedConfig, wg *sync.WaitGroup, feeds chan<- Feed) {
	requestCtx, cancel := context.WithCancel(ctx)

	candidates := feedCfg.URLCandidates()

	content := make(chan io.Reader)

	for _, url := range candidates {
		go l.fetch(requestCtx, cancel, url, content)
	}

	for feedContent := range content {
		feed := InitFeed(feedContent, feedCfg)
		feeds <- feed
	}

	wg.Done()
}

// concurrentLoad
func (l *Loader) fetch(ctx context.Context, cancel func(), url string, content chan<- io.Reader) {
	l.logger.Printf("loading %q", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		l.logger.Printf("could not create request: %s", err)

		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		l.logger.Printf("could not do request: %s", err)

		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		l.logger.Printf("skipping candidate (%q): response status is not OK", url)

		return
	}

	ct, ok := resp.Header["Content-Type"]
	if !ok || len(ct) == 0 || !strings.Contains(ct[0], "xml") {
		l.logger.Printf("skipping candidate (%q): Content-Type should contain \"xml\"", url)

		return
	}

	l.logger.Printf("%q successfully loaded", url)

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		l.logger.Printf("could not read response body: %s", err)

		return
	}

	// candidate has passed checks: stop all other concurrent requests
	cancel()

	content <- bytes.NewReader(b)
	close(content)
}
