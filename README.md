# Google Generative AI Go Patterns

`google/generative-ai-go` SDKを使用したGemini APIの基本的な使用パターンを示すサンプルプロジェクトです。

## 前提条件

- Go 1.21以上
- Google AI StudioのAPIキー

## セットアップ

1. リポジトリをクローンします。
2. 依存関係をインストールします。

```bash
go mod tidy
```

3. 環境変数を設定します。

```bash
export GEMINI_API_KEY="your_api_key_here"
```

## 実行方法

```bash
go run cmd/gemini-patterns/main.go
```

## パターン一覧

このプロジェクトには以下の4つのパターンが含まれています。

1.  **Simple Text Generation** (`internal/patterns/simple_text.go`)
    - 基本的なテキスト生成の例です。
2.  **Chat Session** (`internal/patterns/chat.go`)
    - 履歴を保持するチャットセッションの例です。
3.  **Multimodal** (`internal/patterns/multimodal.go`)
    - 画像とテキストを組み合わせたマルチモーダル入力の例です。
    - 実行ディレクトリに `image.png` を配置する必要があります。
4.  **Streaming Response** (`internal/patterns/streaming.go`)
    - 生成されたテキストをストリーミングで順次表示する例です。

## ディレクトリ構造

- `cmd/`: アプリケーションのエントリーポイント
- `internal/`: 内部パッケージ
    - `client/`: Geminiクライアントの初期化
    - `config/`: 設定管理
    - `patterns/`: 各使用パターンの実装