package internal

import (
	"io"
	"net/url"
	"strings"
)

// Feed represents
type Feed struct {
	content io.Reader
	config  FeedConfig
	name    string
}

// NewFeed
func InitFeed(content io.Reader, config FeedConfig) Feed {
	var name string

	parsed, err := url.Parse(config.URL)
	if err != nil {
		name = strings.ReplaceAll(config.URL, "/", "")
	}

	name = parsed.Hostname()

	return Feed{content, config, name}
}

// Reader
func (f *Feed) Reader() io.Reader {
	return f.content
}
