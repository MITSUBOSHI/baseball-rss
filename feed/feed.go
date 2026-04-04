package feed

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/MITSUBOSHI/baseball-rss/config"
)

type Article struct {
	Title       string
	Link        string
	Description string
	Published   time.Time
	Source      string
}

func FetchAll(ctx context.Context, feeds []config.Feed) []Article {
	var (
		mu       sync.Mutex
		articles []Article
		wg       sync.WaitGroup
		sem      = make(chan struct{}, 3)
	)

	for _, f := range feeds {
		wg.Add(1)
		go func(f config.Feed) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			items, err := fetchFeed(ctx, f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to fetch %s (%s): %v\n", f.Name, f.URL, err)
				return
			}

			mu.Lock()
			articles = append(articles, items...)
			mu.Unlock()
		}(f)
	}

	wg.Wait()
	return articles
}

func fetchFeed(ctx context.Context, f config.Feed) ([]Article, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	parser := gofeed.NewParser()
	parsed, err := parser.ParseURLWithContext(f.URL, ctx)
	if err != nil {
		return nil, fmt.Errorf("parse feed: %w", err)
	}

	var articles []Article
	for _, item := range parsed.Items {
		var published time.Time
		if item.PublishedParsed != nil {
			published = *item.PublishedParsed
		} else if item.UpdatedParsed != nil {
			published = *item.UpdatedParsed
		}

		desc := item.Description
		if desc == "" {
			desc = item.Content
		}

		articles = append(articles, Article{
			Title:       item.Title,
			Link:        item.Link,
			Description: desc,
			Published:   published,
			Source:      f.Name,
		})
	}

	return articles, nil
}
