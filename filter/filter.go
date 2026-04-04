package filter

import (
	"strings"
	"time"

	"github.com/MITSUBOSHI/baseball-rss/feed"
)

type MatchedArticle struct {
	feed.Article
	MatchedKeywords []string
}

func Filter(articles []feed.Article, keywords []string, since time.Duration) []MatchedArticle {
	cutoff := time.Now().Add(-since)
	seen := make(map[string]bool)
	var matched []MatchedArticle

	for _, a := range articles {
		if !a.Published.IsZero() && a.Published.Before(cutoff) {
			continue
		}

		if seen[a.Link] {
			continue
		}
		seen[a.Link] = true

		var hits []string
		text := a.Title + " " + a.Description
		for _, kw := range keywords {
			if strings.Contains(text, kw) {
				hits = append(hits, kw)
			}
		}

		if len(hits) > 0 {
			matched = append(matched, MatchedArticle{
				Article:         a,
				MatchedKeywords: hits,
			})
		}
	}

	return matched
}
