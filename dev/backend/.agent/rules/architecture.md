---
trigger: model_decision
description: バックエンドディレクトリ構造についてのルール。実装タスクにおいて、どの処理をどのディレクトリのどのファイルに記載するのかを決める判断材料として使用する。
---

# ディレクトリ構造

```
backend/
├── cmd/
│   └── api/
│       └── main.go           # エントリーポイント（main関数）
├── config/
│   └── config.go             # 環境変数や設定ファイルの読み込み
├── internal/                 # 外部からimportできない非公開コード（アプリの核心）
│   ├── domain/               # ドメインモデル（struct）やインターフェース定義
│   │   ├── model/          # ドメインエンティティ
│   │   ├── repository/     # リポジトリのインターフェース
│   │   └──  usecase/     # ユースケースのインターフェース
│   ├── handler/              # Echoのハンドラー (HTTP層)
│   │   ├── user.go
│   │   └── model/        # レスポンス、リクエストモデル
│   ├── usecase/              # ビジネスロジック (Service層)
│   │   ├── user.go
│   │   └── ...
│   ├── repository/           # DB操作の実装 (Infrastructure層)
│   │   ├── user.go
│   │   └── ...
│   └── router/               # ルーティング定義
│       └── router.go
├── db/migrations/               # DBマイグレーションファイル
├── go.mod
├── go.sum
└── Dockerfile
```

# 各層の説明

## 1. Domain 層(`internal/domain`)

アプリの「核」となる部分で、ここが他の層に依存してはいけない。

### 責務・ルール

- アプリケーションで扱うデータ構造(EntityModel)を定義する。
- Infrastructure 層のインターフェースを定義する。
- 外部ライブラリへの依存を極力避ける。

```go
package domain

import "context"

// internal/domain/model
// User: アプリケーション内でのユーザー表現
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// internal/domain/repository
// UserRepository: データ保存層が守るべきルール（インターフェース）
// ※ 実装はここには書かない
type UserRepository interface {
    GetByID(ctx context.Context, id int) (*User, error)
    Store(ctx context.Context, user *User) error
}

// UserUsecase: ビジネスロジック層が守るべきルール
type UserUsecase interface {
    GetByID(ctx context.Context, id int) (*User, error)
    Register(ctx context.Context, user *User) error
}
```

## 2. Infrastructure 層(`internal/repository`)

データの永続化（DB 操作）を担当します。

### 責務・ルール

- インターフェースで定義した処理の実装部分を担当する。
- SQL の実行や ORM の操作を行う。
- DB のデータ構造をドメインエンティティに変換して返す。

```go
package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"backend/internal/domain"
)

// userORM: DBのテーブル定義とマッピングするための構造体
// この層（Infrastructure）だけで使用し、外部には公開しません。
type userORM struct {
	ID    int    `gorm:"primaryKey"`
	Name  string `gorm:"size:255"`
	Email string `gorm:"unique;not null"`
	// 必要であれば CreatedAt time.Time などを追加
}

// TableName: テーブル名を明示的に指定する場合（任意）
func (userORM) TableName() string {
	return "users"
}

// toDomain: DB用構造体からドメインモデルへの変換メソッド
func (orm *userORM) toDomain() *domain.User {
	return &domain.User{
		ID:    orm.ID,
		Name:  orm.Name,
		Email: orm.Email,
	}
}

// fromDomain: ドメインモデルからDB用構造体への変換（保存時などに使用）
func fromDomain(u *domain.User) *userORM {
	return &userORM{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

// -------------------------------------------------------

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository: GORMのDBインスタンスを受け取るように変更
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	slog.DebugContext(ctx, "〇〇処理を開始", "id", id)
	var orm userORM

	// GORMによる検索
	// WithContext(ctx) を使ってコンテキストを伝播させるのが重要です
	err := r.db.WithContext(ctx).First(&orm, id).Error

	if err != nil {
		// レコードが見つからない場合のエラーハンドリング
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// gormのエラーを返す
			return nil, gorm.ErrRecordNotFound)
		}
		return nil, err
	}

	// DB用構造体(ORM)からドメインモデルへ変換して返す
	return orm.toDomain(), nil
}

func (r *userRepository) Store(ctx context.Context, user *domain.User) error {
	slog.DebugContext(ctx, "〇〇処理を開始", "user_id", user.ID)
	orm := fromDomain(user)

	// 新規作成 (Create)
	err := r.db.WithContext(ctx).Create(orm).Error
	if err != nil {
		return err
	}

	// GORMはIDを自動採番してormに入れるので、必要ならドメインモデルに戻す
	user.ID = orm.ID
	return nil
}
```

## 3. Usecase / Service 層 (internal/usecase)

ビジネスロジック（業務ルール）を担当します。

### 責務・ルール

- domain/usecase で定義したインターフェースの実装部分を担当する。
- HTTP のこと（Echo など）は一切知らない。
- バリデーションやトランザクション制御、外部 API コールなどを行う。
- Repository のインターフェースを通して DB のデータを操作する。

```go
package usecase

import (
    "context"
    "errors"
    "backend/internal/domain"
)

type userUsecase struct {
    userRepo domain.UserRepository
}

// NewUserUsecase: Repositoryを注入してもらう
func NewUserUsecase(userRepo domain.UserRepository) domain.UserUsecase {
    return &userUsecase{userRepo: userRepo}
}

func (u *userUsecase) GuestLogin(ctx context.Context, userID string) (string, error) {
	slog.InfoContext(ctx, "ゲストログイン処理を開始", "user_id", userID)
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("ユーザー検索に失敗: %w", err)
	}

	token, err := u.generateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("トークン生成に失敗: %w", err)
	}

	return token, nil
}

// Registerメソッドの実装は省略
func (u *userUsecase) Register(ctx context.Context, user *domain.User) error { return nil }
```

## 4. Interface Adapter / HTTP 層 (internal/handler)

Echo との接点。

### 責務・ルール

- ここで初めて echo.Context が登場する。
- リクエストパラメータのパース（JSON Bind, Query Param 取得）。
- Usecase を呼び出す。
- Usecase の結果やエラーを、適切な HTTP ステータスコード（200, 400, 500）に変換してレスポンスする。
- 複雑なロジックは書かない。

```
package handler

import (
    "net/http"
    "strconv"
    "github.com/labstack/echo/v4"
    "backend/internal/domain"
)

// UserHandler: 具体的な構造体（インターフェースではないことが多い）
type UserHandler struct {
    uUsecase domain.UserUsecase
}

func NewUserHandler(u domain.UserUsecase) *UserHandler {
    return &UserHandler{uUsecase: u}
}

func (h *UserHandler) GetByID(c echo.Context) error {
    // 1. リクエストのパース
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid id"})
    }

    // 2. コンテキストの取得（トレーシングやタイムアウト用）
    ctx := c.Request().Context()

    // 3. Usecaseの呼び出し
    user, err := h.uUsecase.GetByID(ctx, id)
    if err != nil {
        // 本来はエラーの種類によって 404 か 500 かなどを分岐させる
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
    }

    // 4. レスポンス
    return c.JSON(http.StatusOK, user)
}
```