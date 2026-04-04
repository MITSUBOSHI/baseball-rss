package summarize

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/liushuangls/go-anthropic/v2"

	"github.com/MITSUBOSHI/baseball-rss/config"
	"github.com/MITSUBOSHI/baseball-rss/filter"
)

func Summarize(ctx context.Context, cfg config.Anthropic, articles []filter.MatchedArticle) (string, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("ANTHROPIC_API_KEY environment variable is not set")
	}

	client := anthropic.NewClient(apiKey)

	var sb strings.Builder
	for i, a := range articles {
		fmt.Fprintf(&sb, "### 記事 %d\n", i+1)
		fmt.Fprintf(&sb, "- タイトル: %s\n", a.Title)
		fmt.Fprintf(&sb, "- ソース: %s\n", a.Source)
		fmt.Fprintf(&sb, "- URL: %s\n", a.Link)
		fmt.Fprintf(&sb, "- マッチキーワード: %s\n", strings.Join(a.MatchedKeywords, ", "))
		if a.Description != "" {
			fmt.Fprintf(&sb, "- 内容: %s\n", a.Description)
		}
		sb.WriteString("\n")
	}

	resp, err := client.CreateMessages(ctx, anthropic.MessagesRequest{
		Model: anthropic.Model(cfg.Model),
		MaxTokens: 4096,
		System: "あなたはNPB（日本プロ野球）のニュース要約アシスタントです。" +
			"与えられた記事一覧を読み、各記事を1〜2文で簡潔に日本語で要約してください。" +
			"マッチしたキーワード（球団名・選手名）に関連する情報を重点的に要約してください。",
		Messages: []anthropic.Message{
			{
				Role: anthropic.RoleUser,
				Content: []anthropic.MessageContent{
					anthropic.NewTextMessageContent("以下の記事を要約してください:\n\n" + sb.String()),
				},
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("Claude API error: %w", err)
	}

	for _, block := range resp.Content {
		if block.Type == "text" {
			return *block.Text, nil
		}
	}

	return "", fmt.Errorf("no text response from Claude API")
}
