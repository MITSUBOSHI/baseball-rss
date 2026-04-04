package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/MITSUBOSHI/baseball-rss/filter"
)

func PrintTerminal(w io.Writer, articles []filter.MatchedArticle, summary string) {
	fmt.Fprintf(w, "=== Baseball News ===\n\n")
	fmt.Fprintf(w, "マッチした記事: %d件\n\n", len(articles))

	for _, a := range articles {
		fmt.Fprintf(w, "[%s] %s\n", a.Source, a.Title)
		fmt.Fprintf(w, "  キーワード: %s\n", strings.Join(a.MatchedKeywords, ", "))
		if !a.Published.IsZero() {
			fmt.Fprintf(w, "  日時: %s\n", a.Published.Format("2006-01-02 15:04"))
		}
		fmt.Fprintf(w, "  → %s\n\n", a.Link)
	}

	if summary != "" {
		fmt.Fprintf(w, "--- 要約 ---\n\n%s\n", summary)
	}
}

func PrintMarkdown(w io.Writer, articles []filter.MatchedArticle, summary string) {
	fmt.Fprintf(w, "# Baseball News\n\n")
	fmt.Fprintf(w, "マッチした記事: **%d件**\n\n", len(articles))

	fmt.Fprintf(w, "## 記事一覧\n\n")
	for _, a := range articles {
		fmt.Fprintf(w, "### %s\n\n", a.Title)
		fmt.Fprintf(w, "- **ソース**: %s\n", a.Source)
		fmt.Fprintf(w, "- **キーワード**: %s\n", strings.Join(a.MatchedKeywords, ", "))
		if !a.Published.IsZero() {
			fmt.Fprintf(w, "- **日時**: %s\n", a.Published.Format("2006-01-02 15:04"))
		}
		fmt.Fprintf(w, "- **リンク**: %s\n\n", a.Link)
	}

	if summary != "" {
		fmt.Fprintf(w, "## 要約\n\n%s\n", summary)
	}
}
