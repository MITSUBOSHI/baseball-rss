package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/MITSUBOSHI/baseball-rss/filter"
)

func PrintTerminal(w io.Writer, articles []filter.MatchedArticle, summary string) {
	_, _ = fmt.Fprintf(w, "=== Baseball News ===\n\n")
	_, _ = fmt.Fprintf(w, "マッチした記事: %d件\n\n", len(articles))

	for _, a := range articles {
		_, _ = fmt.Fprintf(w, "[%s] %s\n", a.Source, a.Title)
		_, _ = fmt.Fprintf(w, "  キーワード: %s\n", strings.Join(a.MatchedKeywords, ", "))
		if !a.Published.IsZero() {
			_, _ = fmt.Fprintf(w, "  日時: %s\n", a.Published.Format("2006-01-02 15:04"))
		}
		_, _ = fmt.Fprintf(w, "  → %s\n\n", a.Link)
	}

	if summary != "" {
		_, _ = fmt.Fprintf(w, "--- 要約 ---\n\n%s\n", summary)
	}
}

func PrintMarkdown(w io.Writer, articles []filter.MatchedArticle, summary string) {
	_, _ = fmt.Fprintf(w, "# Baseball News\n\n")
	_, _ = fmt.Fprintf(w, "マッチした記事: **%d件**\n\n", len(articles))

	_, _ = fmt.Fprintf(w, "## 記事一覧\n\n")
	for _, a := range articles {
		_, _ = fmt.Fprintf(w, "### %s\n\n", a.Title)
		_, _ = fmt.Fprintf(w, "- **ソース**: %s\n", a.Source)
		_, _ = fmt.Fprintf(w, "- **キーワード**: %s\n", strings.Join(a.MatchedKeywords, ", "))
		if !a.Published.IsZero() {
			_, _ = fmt.Fprintf(w, "- **日時**: %s\n", a.Published.Format("2006-01-02 15:04"))
		}
		_, _ = fmt.Fprintf(w, "- **リンク**: %s\n\n", a.Link)
	}

	if summary != "" {
		_, _ = fmt.Fprintf(w, "## 要約\n\n%s\n", summary)
	}
}
