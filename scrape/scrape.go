package scrape

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/MITSUBOSHI/baseball-rss/config"
	"github.com/MITSUBOSHI/baseball-rss/feed"
)

func ScrapeAll(ctx context.Context, sites []config.Site) []feed.Article {
	var (
		mu       sync.Mutex
		articles []feed.Article
		wg       sync.WaitGroup
		sem      = make(chan struct{}, 3)
	)

	for _, s := range sites {
		wg.Add(1)
		go func(s config.Site) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			items, err := scrapeSite(ctx, s)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to scrape %s (%s): %v\n", s.Name, s.URL, err)
				return
			}

			mu.Lock()
			articles = append(articles, items...)
			mu.Unlock()
		}(s)
	}

	wg.Wait()
	return articles
}

func scrapeSite(ctx context.Context, s config.Site) ([]feed.Article, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", "baseball-rss/1.0")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse HTML: %w", err)
	}

	baseURL, err := url.Parse(s.URL)
	if err != nil {
		return nil, fmt.Errorf("parse base URL: %w", err)
	}

	seen := make(map[string]bool)
	var articles []feed.Article

	doc.Find(s.LinkSelector).Each(func(_ int, sel *goquery.Selection) {
		href, exists := sel.Attr("href")
		if !exists || href == "" {
			return
		}

		link := resolveURL(baseURL, href)
		if seen[link] {
			return
		}
		seen[link] = true

		title := extractTitle(sel, s.TitleSelector)
		if title == "" || len([]rune(title)) < 5 {
			return
		}
		// Truncate overly long titles (likely including description text)
		if runes := []rune(title); len(runes) > 80 {
			title = string(runes[:80]) + "…"
		}

		articles = append(articles, feed.Article{
			Title:       strings.TrimSpace(title),
			Link:        link,
			Description: "",
			Published:   time.Time{},
			Source:      s.Name,
		})
	})

	return articles, nil
}

func extractTitle(sel *goquery.Selection, titleSelector string) string {
	if titleSelector != "" {
		if t := sel.Find(titleSelector).First().Text(); t != "" {
			return cleanText(t)
		}
	}
	// Fallback: use the link text itself (excluding image alt text noise)
	clone := sel.Clone()
	clone.Find("img, script, style").Remove()
	clone.Find(".category, .time, .date").Remove()
	return cleanText(clone.Text())
}

func cleanText(s string) string {
	// Collapse whitespace (newlines, tabs, multiple spaces)
	var b strings.Builder
	prevSpace := false
	for _, r := range s {
		if r == '\n' || r == '\r' || r == '\t' || r == ' ' || r == '\u00a0' {
			if !prevSpace {
				b.WriteRune(' ')
				prevSpace = true
			}
			continue
		}
		b.WriteRune(r)
		prevSpace = false
	}
	return strings.TrimSpace(b.String())
}

func resolveURL(base *url.URL, href string) string {
	parsed, err := url.Parse(href)
	if err != nil {
		return href
	}
	return base.ResolveReference(parsed).String()
}
