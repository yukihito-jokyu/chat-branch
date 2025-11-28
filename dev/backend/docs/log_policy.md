# ログ出力方針 (Log Output Policy)

本プロジェクトにおけるログ出力のベストプラクティスとルールを定めます。
各層（Handler, Usecase, Repository）での責務に応じた適切なログ出力を心がけてください。

## 1. ログライブラリとフォーマット

- **ライブラリ**: 標準ライブラリの `log/slog` (Go 1.21+) を使用することを推奨します。
- **フォーマット**:
  - **本番環境 (Production)**: JSON形式 (`slog.JSONHandler`)。機械可読性を高め、ログ収集基盤での解析を容易にするため。
  - **開発環境 (Development)**: テキスト形式 (`slog.TextHandler`)。人間が読みやすくするため。

## 2. ログレベルの基準

| レベル | 説明 | 使用例 |
| :--- | :--- | :--- |
| **ERROR** | システムの動作継続に影響を与えるエラー、または予期しない例外。即時の対応が必要な場合が多い。 | DB接続エラー、外部APIのダウン、パニックのリカバリ、500 Internal Server Errorとなる事象。 |
| **WARN** | システムは継続可能だが、注意が必要な事象。または、クライアント起因の不正なリクエスト（バリデーションエラー等）で頻発するもの。 | 認証失敗、バリデーションエラー（400 Bad Request）、リトライ発生時。 |
| **INFO** | 正常な動作の記録。主要なビジネスプロセスの開始・終了、状態遷移。 | ユーザー登録完了、ログイン成功、バッチ処理の開始・終了。 |
| **DEBUG** | 開発・デバッグ時に詳細な情報を追跡するためのログ。本番環境では通常出力しない。 | 変数の中身、詳細な分岐の通過確認、SQLクエリ（機密情報を除く）。 |

## 3. 共通項目 (Context)

ログには可能な限り以下のコンテキスト情報を含めてください。`slog.With` や `context` から抽出して付与します。

- `request_id`: リクエストを一意に識別するID
- `user_id`: 操作を行っているユーザーのID（認証済みの場合）
- `trace_id`: 分散トレーシング用のID（導入している場合）

## 4. レイヤー別のログ出力方針

各層の責務に応じて、出力すべきログの内容が異なります。`auth.go` を例に説明します。

### 4.1. Handler層 (`internal/handler`)

**責務**: HTTPリクエストの受付、バリデーション、レスポンスの返却。

- **ログ出力のタイミング**:
  - リクエスト処理の開始と終了（通常はMiddlewareで行うため、個別のHandlerメソッド内では不要な場合が多い）。
  - **エラー発生時**: Usecaseからエラーが返ってきた場合、クライアントに返すステータスコードと共にログ出力する。
- **推奨ログ**:
  - **WARN**: バリデーションエラー、認証エラー（401/403）。
  - **ERROR**: 500エラーとなる予期せぬエラー。`err` フィールドを含める。
- **`auth.go` の例**:
  ```go
  func (h *AuthHandler) Login(c echo.Context) error {
      // ...
      if err != nil {
          // 500エラーの場合はERRORログ
          slog.ErrorContext(ctx, "ゲストログインに失敗", "error", err)
          return c.JSON(http.StatusInternalServerError, ...)
      }
      // 正常終了時はMiddlewareのアクセスログで十分な場合が多いが、重要な操作はINFOで出す
      slog.InfoContext(ctx, "ゲストログインに成功", "user_id", userID)
      return c.JSON(http.StatusOK, ...)
  }
  ```

### 4.2. Usecase層 (`internal/usecase`)

**責務**: ビジネスロジックの実行。

- **ログ出力のタイミング**:
  - ビジネスロジックの重要なフローの通過点。
  - 分岐条件の決定打となった情報。
- **推奨ログ**:
  - **INFO**: ビジネス的に意味のある操作の成功（「ユーザー作成完了」「メール送信完了」）。
  - **WARN**: ビジネスルールによる拒否（「在庫切れ」「アカウントロック中」）。
- **`auth.go` の例**:
  ```go
  func (u *authUsecase) GuestSignup(ctx context.Context) (*model.User, string, error) {
      // ...
      if err := u.userRepo.Create(ctx, user); err != nil {
          // ここではログを出さず、エラーをラップして返す（呼び出し元でログ出力する方針も可だが、
          // 詳細なコンテキストを知っているここでWARN/ERRORを出すのもあり。
          // 基本方針としては「エラーは発生源に近いところでログに残す」か「Handlerでまとめて出す」か統一する。
          // 推奨: エラーの文脈付与（fmt.Errorf）を行い、ログはHandlerまたはMiddlewareで集約する。
          // ただし、重要なビジネスイベントの失敗はここでWARNを出しても良い。
          return nil, "", fmt.Errorf("ユーザー作成に失敗: %w", err)
      }
      
      slog.InfoContext(ctx, "ゲストユーザー作成に成功", "user_id", user.ID)
      return user, token, nil
  }
  ```

### 4.3. Repository層 (`internal/repository`)

**責務**: データの永続化（DB操作）。

- **ログ出力のタイミング**:
  - 基本的にはログ出力を行わず、エラーを返すことに専念する。
  - ただし、**DB接続エラー**や**スロークエリ**など、インフラ寄りの重大な問題はここでERRORログを出す場合がある。
- **推奨ログ**:
  - **DEBUG**: 発行されたSQL（開発時のみ有効にする）。
  - **ERROR**: 接続断などの致命的なエラー。
- **`user.go` の例**:
  ```go
  func (r *userRepository) Create(ctx context.Context, user *model.User) error {
      result := r.db.WithContext(ctx).Create(dto)
      if result.Error != nil {
          // 通常はエラーを返すだけにする
          return result.Error
      }
      return nil
  }
  ```
