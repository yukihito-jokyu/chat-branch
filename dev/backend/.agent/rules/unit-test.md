---
trigger: model_decision
description: 単体テスト作成時に必読するルール
---

# 単体テスト作成ルール

このドキュメントでは、バックエンドアプリケーションにおける単体テストの実装ルールを定めます。

## 1. 基本方針

- **テスト手法**: 全ての単体テストは **テーブル駆動テスト (Table-Driven Tests)** 形式で実装します。
- **モックの活用**: テスト対象が依存するコンポーネント（インターフェース）は全てモック化し、テスト対象の責務のみを検証します。
- **パッケージ構成**: テストファイルは実装ファイルと同じディレクトリに配置し、ファイル名は `*_test.go` とします。
- C1レベルのテストを実装するものとする。

## 2. 推奨ライブラリ

テストの実装には以下のライブラリの使用を推奨します。

- **アサーション**: `github.com/stretchr/testify/assert`
- **モック**: `github.com/stretchr/testify/mock`
- **DBテストドライバ (Repository層)**: `github.com/glebarez/sqlite` (Pure Go SQLite driver)

## 3. 各層のテスト実装ルール

### 3.1 Handler層 (`internal/handler`)

- **責務**: HTTPリクエストの受け取り、パラメータのパース、Usecaseの呼び出し、HTTPレスポンスの返却。
- **依存**: `Usecase` インターフェース。
- **テスト内容**:
    - ステータスコードが期待通りか。
    - レスポンスボディ（JSON等）が期待通りか。
    - Usecaseが適切な引数で呼び出されているか。
    - エラー時に適切なHTTPエラーを返しているか。
- **実装方法**:
    - `net/http/httptest` と Echoの `NewContext` を使用してリクエスト/レスポンスをシミュレートします。
    - Usecaseのモックを注入します。

### 3.2 Usecase層 (`internal/usecase`)

- **責務**: ビジネスロジックの実行、Repositoryの呼び出し。
- **依存**: `Repository` インターフェース。
- **テスト内容**:
    - ビジネスロジックが正しく機能しているか。
    - Repositoryが適切な引数で呼び出されているか。
    - Repositoryからの戻り値（成功/エラー）に応じた挙動をするか。
- **実装方法**:
    - Repositoryのモックを注入します。

### 3.3 Repository層 (`internal/repository`)

- **責務**: データベースへのクエリ発行、データマッピング。
- **依存**: `gorm.DB`。
- **テスト内容**:
    - 実際にデータが保存・取得できるか（GORMの挙動確認）。
    - DBからの返却値が正しくモデルにマッピングされるか。
    - 意図したデータのみが取得できているか（ScopeやWhere句の検証）。
- **実装方法**:
    - `github.com/glebarez/sqlite` を使用して、インメモリデータベース (`:memory:`) を構築します。
    - テストごとに `AutoMigrate` を実行し、必要な初期データを投入して検証を行います。
    - ※ MySQL固有の機能を使用している場合は、可能な限り標準的なSQLに寄せるか、テスト時のみ代替手段を検討します。

## 4. テーブル駆動テストのテンプレート構造

各テスト関数は以下の構造を持つようにします。

```go
func Test_TargetFunction(t *testing.T) {
    // テストケースの定義構造体
    type args struct {
        // テスト対象関数の引数
        arg1 string
        arg2 int
    }
    
    tests := []struct {
        name      string          // テストケース名
        args      args            // 入力引数
        setupMock func(m *mocks)  // モックの期待値設定
        wantErr   bool            // エラーを期待するか
        assertion func(t *testing.T, got ResultType) // 結果の検証
    }{
        {
            name: "正常系: 〇〇の場合",
            args: args{...},
            setupMock: func(m *mocks) {
                // モックの呼び出し期待値と返却値を設定
                m.dependency.On("Method", ...).Return(...)
            },
            wantErr: false,
            assertion: func(t *testing.T, got ResultType) {
                assert.Equal(t, expected, got)
            },
        },
        // ... 他のテストケース
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 1. モックの初期化
            // 2. テスト対象の初期化 (依存性の注入)
            // 3. モックの設定 (tt.setupMock)
            // 4. テスト対象の実行
            // 5. アサーション (tt.wantErr, tt.assertion)
        })
    }
}
```