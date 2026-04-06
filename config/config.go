package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Watch     Watch      `yaml:"watch"`
	Feeds     []Feed     `yaml:"feeds"`
	Sites     []Site     `yaml:"sites"`
	Anthropic Anthropic  `yaml:"anthropic"`
}

type Watch struct {
	Teams   []string `yaml:"teams"`
	Players []string `yaml:"players"`
}

func (w Watch) Keywords() []string {
	return append(w.Teams, w.Players...)
}

type Feed struct {
	URL  string `yaml:"url"`
	Name string `yaml:"name"`
}

type Site struct {
	URL           string `yaml:"url"`
	Name          string `yaml:"name"`
	LinkSelector  string `yaml:"link_selector"`
	TitleSelector string `yaml:"title_selector"`
}

type Anthropic struct {
	Model string `yaml:"model"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if len(cfg.Feeds) == 0 && len(cfg.Sites) == 0 {
		return nil, fmt.Errorf("no feeds or sites configured")
	}

	if len(cfg.Watch.Teams) == 0 && len(cfg.Watch.Players) == 0 {
		return nil, fmt.Errorf("no watch targets configured")
	}

	if cfg.Anthropic.Model == "" {
		cfg.Anthropic.Model = "claude-sonnet-4-20250514"
	}

	return &cfg, nil
}
