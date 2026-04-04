# baseball-rss

NPB野球ニュースのRSSフィードからwatch対象の球団・選手に関する記事を取得し、Claude APIで要約するCLIツール。

## セットアップ

```bash
go build -o baseball-rss .
cp config.yaml.example config.yaml
# config.yaml を編集してwatch対象を設定
```

## 設定

`config.yaml` でwatch対象とフィードURLを管理します。

```yaml
watch:
  teams:
    - 横浜DeNAベイスターズ
    - DeNA
    - ベイスターズ
  players:
    - 山本祐大
    - 坂本裕哉

feeds:
  - url: https://full-count.jp/feed/
    name: Full-Count
  - url: https://www.baseballchannel.jp/feed/
    name: Baseball Channel
  - url: https://baseballking.jp/feed
    name: Baseball King
```

## 使い方

### CLI

```bash
# 要約なし（APIキー不要）
./baseball-rss --config config.yaml --no-summary

# 要約あり
ANTHROPIC_API_KEY=sk-... ./baseball-rss --config config.yaml

# Markdown出力、過去1週間分
./baseball-rss --config config.yaml --format markdown --since 168h
```

### Claude Code スキル

このリポジトリ内で `/baseball-rss` を実行すると、CLIツールを使って記事を取得し、Claudeが要約します。
