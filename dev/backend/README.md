# backend

チャットアプリケーションのバックエンドプロジェクトです。Go言語で実装されており、Echoフレームワーク、Gorm、Gemini APIなどを使用しています。

## 開発環境のセットアップ

このプロジェクトはVS CodeのDev Containersを利用して開発環境を構築することを前提としています。

### 1. Dev Containerの起動

VS Codeでこのディレクトリ（`backend`）を開き、左下の緑色のボタンをクリックして「Reopen in Container」を選択するか、コマンドパレットから「Dev Containers: Reopen in Container」を実行してください。
必要なツール（Go, MySQL client, Taskなど）が含まれたDockerコンテナが起動します。

### 2. 設定ファイルの作成

`config` ディレクトリにあるテンプレートファイルをコピーして、設定ファイルを作成します。

```bash
cp config/config.yml.template config/config.yml
```

`config/config.yml` を開き、以下の項目を設定してください。

*   **gemini.apiKey**: Google AI Studioで取得したGemini APIキーを入力してください。
*   **jwt.secret**: JWT署名用のシークレットキー（開発用なら適当な文字列で可）。
*   **database**: データベース接続情報（Dev Container内のDBサービスを使用する場合はデフォルトのままで動作します）。

### 3. 初期化とマイグレーション

Taskfileを使用して初期化とデータベースのマイグレーションを行います。

```bash
# 拡張機能のインストール（初回のみ）
task ex-init

# データベースツールのインストール（goose）
task db-init

# マイグレーションの適用
task db-up
```

## 開発コマンド (Taskfile)

このプロジェクトでは `Taskfile.yml` を使用して各種コマンドを管理しています。

| コマンド | 説明 |
| :--- | :--- |
| `task ex-init` | VS Code拡張機能のインストールスクリプトを実行します。 |
| `task db-init` | マイグレーションツール `goose` をインストールします。 |
| `task db-up` | データベースにマイグレーションを適用します。 |
| `task db-down` | データベースのマイグレーションをロールバックします。 |
| `task db-status` | マイグレーションの適用状況を確認します。 |
| `task db-create -- <name>` | 新しいマイグレーションファイルを作成します。<br>例: `task db-create -- create_users_table` |
| `task db-reset` | データベースをリセット（全削除＆再適用）します。 |
| `task run` | サーバーを起動します（`go run cmd/server/main.go` 相当）。 |
| `task build` | アプリケーションをビルドします。 |
| `task test` | 全テストを実行します。 |
| `task coverage` | テストカバレッジを計測し、HTMLレポートを出力します。 |
| `task lint` | `golangci-lint` を実行してコードをチェックします。 |
| `task format` | `goimports` を使用してコードをフォーマットします。 |

## ディレクトリ構造

クリーンアーキテクチャを意識した構成になっています。

```
backend/
├── .devcontainer/ # Dev Container設定
├── cmd/           # エントリーポイント
│   └── server/    # サーバー起動用メインファイル
├── config/        # 設定ファイルと読み込みロジック
├── db/            # データベース関連
│   └── migrations/# マイグレーションSQLファイル
├── internal/      # 外部からインポートされない内部パッケージ
│   ├── domain/    # ドメイン層（インターフェース、ドメインモデル定義）
│   ├── handler/   # インターフェース層（HTTPハンドラー、リクエスト/レスポンス定義）
│   ├── usecase/   # ユースケース層（ビジネスロジック）
│   ├── repository/# インフラ層（DB操作の実装）
│   └── router/    # ルーティング定義
├── pkg/           # 外部から利用可能な汎用パッケージ
└── Taskfile.yml   # タスクランナー設定
```

## 技術スタック

### バックエンド

| カテゴリ | パッケージ名 | 用途・選定理由 |
| :--- | :--- | :--- |
| 言語 | Go | 静的型付け、高速な実行速度、並行処理に強いため。 |
| Webフレームワーク | Echo (labstack/echo/v4) | 軽量かつ機能豊富なGo用フレームワーク。ルーティングやJSON処理が容易。 |
| ORM | Gorm | データベース操作をGoの構造体で行うため。マイグレーション機能も利用。 |
| 認証 (JWT) | golang-jwt/jwt/v5 | ログイン認証およびAPIリクエストの検証用トークン生成。 |
| AI SDK | google-generative-ai-go | Gemini APIをGoから型安全に呼び出すためのGoogle公式SDK。 |
| 環境変数 | godotenv | .env ファイルからAPIキーやDB接続情報を読み込むため。 |
| UUID | google/uuid | DBの主キー（UUID）をGo側で生成・操作するため。 |
| マイグレーション | pressly/goose | SQLファイルで履歴管理しつつ、サーバー起動時に自動適用するため。 |

### DB・LLM

| カテゴリ | 名称・バージョン | 用途・選定理由 |
| :--- | :--- | :--- |
| データベース | MySQL 8.0 以上 | チャットデータおよびプロジェクトデータの永続化。再帰クエリ(Recursive CTE) を使ってツリー構造を効率的に扱うために8.0以上が必須。 |
| LLM | Gemini 2.5 Flash | 高速・低コストなモデル。チャットの要約やコンテキスト生成で頻繁にAPIを叩くため、レイテンシとコストのバランスが良いFlashモデルを採用。 |
