---
trigger: model_decision
description: フロントエンドディレクトリ構造についてのルール。実装タスクにおいて、どの画面・機能をどのディレクトリのどのファイルに記載するのかを決める判断材料として使用する。
---

# Frontend Architecture Guidelines

## 1\. ディレクトリ構造の青写真

Feature-based Architecture（機能単位の分割）を採用します。

```text
src/
├── app/                  # アプリケーションの構成とルーティング
│   ├── routes/           # ★ TanStack Router: ファイルベースルーティング定義
│   ├── provider.tsx      # 全てのContext Provider (Query, Theme, Tooltip等)
│   └── router.tsx        # Routerインスタンス生成と設定
│
├── components/           # アプリケーション全体で共有するUI部品
│   ├── ui/               # ★ shadcn/ui: 自動生成されるUIパーツ (Button, Input等)
│   ├── layout/           # レイアウトコンポーネント (Header, Sidebar)
│   └── common/           # プロジェクト独自の共有UI (MarkdownRenderer, LoadingSpinner)
│
├── features/             # ★ ドメイン（機能）ごとの実装。コードの8割はここ。
│   ├── [feature-name]/   # 例: chat, auth, canvas
│   │   ├── api/          # API接続関数 (Axios)
│   │   ├── components/   # この機能内でのみ使うコンポーネント
│   │   ├── hooks/        # カスタムフック
│   │   ├── queries.ts    # ★ TanStack Query: Query Options & Keys定義
│   │   ├── stores/       # ★ Zustand: 機能固有のクライアントState
│   │   └── types.ts      # 機能固有の型定義
│   └── ...
│
├── lib/                  # 外部ライブラリの設定・ラッパー・ユーティリティ
│   ├── api-client.ts     # Axiosインスタンス (Interceptor設定)
│   ├── query-client.ts   # QueryClient設定
│   ├── utils.ts          # shadcn用 cn() 関数
│   └── i18n.ts           # i18next設定
│
├── hooks/                # 全体で汎用的に使えるHooks (useDebounce, useMediaQuery)
├── stores/               # グローバルなZustand Store (UserSession, Settings)
├── types/                # 全体で共有する型定義 (User, APIResponse)
└── main.tsx              # エントリーポイント
```

-----

## 2\. パッケージの責務範囲（Responsibility Matrix）

どのライブラリを使って何を実装すべきか、迷った時の判断基準です。

| カテゴリ | パッケージ | 責務・役割 | ルール |
| :--- | :--- | :--- | :--- |
| **Routing** | `TanStack Router` | URLの管理、ページ遷移、**データのプリフェッチ(Loader)** | ページコンポーネント自体は薄く保ち、ロジックはFeatureへ委譲する。 |
| **Server State** | `TanStack Query` | APIデータの取得、キャッシュ、同期、Loading/Error状態管理 | `useEffect`でのデータ取得は禁止。必ずQuery経由で行う。 |
| **Client State** | `Zustand` | UIの状態管理（モーダル開閉、React Flowのズーム率など） | サーバーデータはQueryに任せ、Zustandには入れない。 |
| **UI System** | `shadcn/ui` + `Tailwind` | 基本的なUIパーツの提供とスタイリング | `components/ui` 内のファイルは手動でロジックを追加しない（スタイル調整はOK）。 |
| **Canvas** | `React Flow` | ノードベースUI、無限キャンバスの描画 | ロジックが複雑になるため、必ず `features/canvas` 等に隔離する。 |
| **API Client** | `Axios` | HTTPリクエストの発行、トークン管理 | 直接コンポーネントで呼ばず、`features/*/api` で関数化してから呼ぶ。 |
| **Schema/Type** | `Zod` (想定) | バリデーション、型推論 | フォームやAPIレスポンスの検証に使用する。 |

-----

## 3\. 実装ルールと判断フロー

### ルール①：コンポーネントの配置場所

「このコンポーネントはどこに置くべきか？」

1.  **`src/components/ui`**: `Button` や `Dialog` など、shadcn/ui で追加したもの。**（ビジネスロジックを持たせない）**
2.  **`src/features/[feature]/components`**: 特定の機能（例：チャット画面）でしか使わないもの。**（原則ここに置く）**
3.  **`src/components/common`**: 複数の機能で使い回すもの（例：Markdown表示、独自のアバター表示）。

### ルール②：状態管理の使い分け

「このデータはどこで管理すべきか？」

1.  **URLパラメータで表現できる？** (ID, 検索クエリ, タブ)
      * YES → **TanStack Router** (`searchParams` / `params`)
2.  **サーバーにあるデータ？** (ユーザー一覧, メッセージログ)
      * YES → **TanStack Query**
3.  **アプリを閉じても保持したい設定？** (テーマ, 言語)
      * YES → **Zustand** (persist middleware)
4.  **一時的なUI操作？** (入力中のテキスト, ドラッグ中の位置)
      * YES → **Zustand** (または `useState`)

### ルール③：TanStack Routerとデータ取得

データ取得は「Render-as-you-fetch」パターンを徹底します。

  * **NG**: コンポーネント内の `useEffect` で fetch する。
  * **OK**: `useQuery` をコンポーネント内で使う（Loadingが発生する）。
  * **BEST**: Routeの **`loader`** で `queryClient.ensureQueryData` を実行し、コンポーネントでは `useSuspenseQuery` を使う（Loadingなしで即表示）。

### ルール④：機能（Feature）の独立性

  * Featureディレクトリ間の相互インポートは極力避ける。
  * もし `features/chat` のコンポーネントを `features/dashboard` で使いたくなったら、それは「汎用コンポーネント」に昇格させる合図 → `src/components/common` へ移動。

-----

## 4\. ファイル作成のテンプレート（Checklist）

新しい機能（例: 「レポート機能」）を追加する際の手順です。

1.  **フォルダ作成**: `src/features/report/` を作成。
2.  **API定義**: `src/features/report/api/index.ts` にAxiosの関数を書く。
3.  **Query定義**: `src/features/report/queries.ts` にQuery KeyとOptions定義を書く。
4.  **コンポーネント作成**: `src/features/report/components/ReportView.tsx` を実装。
5.  **ルーティング追加**: `src/app/routes/report.tsx` を作成し、LoaderとComponentを紐付ける。

-----

このルールセットに従うことで、コードの予測可能性が高まり、プロジェクトが大規模化してもメンテナンス性を維持できます。