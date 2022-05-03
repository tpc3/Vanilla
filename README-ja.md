# Vanilla
[![Go Report Card](https://goreportcard.com/badge/github.com/tpc3/vanilla)](https://goreportcard.com/report/github.com/tpc3/vanilla)
[![Docker Image CI](https://github.com/tpc3/Vanilla/actions/workflows/docker-image.yml/badge.svg)](https://github.com/tpc3/Vanilla/actions/workflows/docker-image.yml)
[![Go](https://github.com/tpc3/Vanilla/actions/workflows/go.yml/badge.svg)](https://github.com/tpc3/Vanilla/actions/workflows/go.yml)

絵文字の統計を取るためのDiscordBotです。

[english](https://github.com/tpc3/Vanilla/README.md)

## 機能
- 任意の期間の使用回数に応じたランキング表示
- メッセージ、リアクションのそれぞれに1回あたりの重みを設定
- 説明文を保存しmd形式でまとめて入出力
- 2言語対応
- Botによるイベントの記録を切替可能
- 絵文字の名前変更、同名での再登録に自動追従
- 複数Guild対応

## Use
### Simple
1. [Releases](https://github.com/tpc3/Vanilla/releases)から実行ファイルをダウンロード
    - 実行ファイルは`artifact.zip`の中にあるかもしれません。
    - 最新版が使いたいですか？[Actions](https://github.com/tpc3/Vanilla/actions/workflows/go.yml)からダウンロードできます。
1. [config.yamlをダウンロード](https://raw.githubusercontent.com/tpc3/Vanilla/master/config.yaml)
1. config.yamlにBotのTokenを入力
1. `./Vanilla`
1. Discordサーバーで`emoji.sync`コマンドを実行

### Docker
1. [config.yamlをダウンロード](https://raw.githubusercontent.com/tpc3/Vanilla/master/config.yaml)
1. config.yamlにBotのTokenを入力
1. `docker run --rm -it -v $(PWD):/data ghcr.io/tpc3/vanilla`
1. Discordサーバーで`emoji.sync`コマンドを実行

## ビルド
1. このリポジトリをcloneする
1. `go build`
### 必要アプリケーション
- git
- golang
- gcc

## Contribute
このリポジトリに貢献したいですか？  
IssueによるバグレポートやPull Requestは歓迎しています。  
英語でも日本語でも構いません。