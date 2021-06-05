package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// Rule represents
type Rules struct {
	PostPath        string `json:"postPath"`
	TitlePath       string `json:"titlePath"`
	LinkPath        string `json:"linkPath"`
	DescriptionPath string `json:"descriptionPath"`
}

type FeedConfig struct {
	URL   string `json:"url"`
	Rules Rules  `json:"rules"`
}

// URLCandidates
func (conf *FeedConfig) URLCandidates() []string {
	url := strings.TrimRight(conf.URL, "/")

	return []string{
		url,
		fmt.Sprintf("%s/feed", url),
		fmt.Sprintf("%s/rss.xml", url),
		fmt.Sprintf("%s/feed.atom", url),
	}
}

type DBConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DBname   string `json:"dbname"`
}

// ConnectURL
func (dbc *DBConfig) String() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		dbc.Username,
		dbc.Password,
		dbc.Host,
		dbc.Port,
		dbc.DBname,
	)
}

// Config represents aggregator config
type Config struct {
	DB          DBConfig     `json:"db"`
	FeedConfigs []FeedConfig `json:"feeds"`
}

// FeedConfigs
func GetFeedConfigs(cfg *Config) []FeedConfig {
	return cfg.FeedConfigs
}

// DBConfig
func GetDBConfig(cfg *Config) *DBConfig {
	return &cfg.DB
}

// InitConfig
func InitConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %w", err)
	}
	defer f.Close()

	rawCfg, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("could not read config: %w", err)
	}

	cfg := &Config{}
	err = json.Unmarshal(rawCfg, cfg)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshal config data: %w", err)
	}

	return cfg, nil
}
