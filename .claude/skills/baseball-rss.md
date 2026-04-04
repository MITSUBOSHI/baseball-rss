---
name: baseball-rss
description: NPB野球ニュースのRSSフィードからwatch対象の球団・選手に関する記事を取得・要約する
user_invocable: true
---

# Baseball RSS Reader

CLIツール `baseball-rss` を使って、watch対象の球団・選手に関するNPB野球ニュースを取得・要約する。

## 手順

1. まず CLIツールをビルドする（未ビルドの場合）:
   ```
   go build -o baseball-rss .
   ```

2. `--no-summary` モードでCLIを実行し、マッチした記事一覧を取得する:
   ```
   ./baseball-rss --config config.yaml --no-summary --since 24h
   ```

3. CLIの出力結果をもとに、各記事を1〜2文で日本語で要約する

4. 最後に全体のまとめを3〜5文で出力する

## オプション

ユーザーが引数を渡した場合はCLIに転送する:
- `--since 168h` → 過去1週間分
- `--format markdown` → Markdown形式で出力

引数がなければデフォルト（過去24時間、ターミナル出力）で実行する。

## 注意事項

- CLIがフィード取得失敗の警告を出した場合、その旨をユーザーに伝える
- マッチする記事がない場合は「マッチする記事はありませんでした」と返す
- 要約はwatch対象の球団・選手に関連する情報を重点的に含める
