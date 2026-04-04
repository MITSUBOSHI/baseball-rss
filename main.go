package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/MITSUBOSHI/baseball-rss/config"
	"github.com/MITSUBOSHI/baseball-rss/feed"
	"github.com/MITSUBOSHI/baseball-rss/filter"
	"github.com/MITSUBOSHI/baseball-rss/output"
	"github.com/MITSUBOSHI/baseball-rss/summarize"
)

func main() {
	var (
		configPath = flag.String("config", "config.yaml", "path to config file")
		format     = flag.String("format", "terminal", "output format: terminal or markdown")
		since      = flag.Duration("since", 24*time.Hour, "only show articles from this duration ago")
		noSummary  = flag.Bool("no-summary", false, "skip AI summarization")
	)
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if !*noSummary && os.Getenv("ANTHROPIC_API_KEY") == "" {
		fmt.Fprintf(os.Stderr, "Error: ANTHROPIC_API_KEY is not set (use --no-summary to skip summarization)\n")
		os.Exit(1)
	}

	ctx := context.Background()

	articles := feed.FetchAll(ctx, cfg.Feeds)
	if len(articles) == 0 {
		fmt.Fprintf(os.Stderr, "Error: failed to fetch any articles from all feeds\n")
		os.Exit(1)
	}

	matched := filter.Filter(articles, cfg.Watch.Keywords(), *since)
	if len(matched) == 0 {
		fmt.Println("マッチする記事はありませんでした。")
		return
	}

	var summary string
	if !*noSummary {
		summary, err = summarize.Summarize(ctx, cfg.Anthropic, matched)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: summarization failed: %v\n", err)
		}
	}

	switch *format {
	case "markdown":
		output.PrintMarkdown(os.Stdout, matched, summary)
	default:
		output.PrintTerminal(os.Stdout, matched, summary)
	}
}
