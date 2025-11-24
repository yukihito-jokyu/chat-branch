import { MarkdownRenderer } from "../components/MarkdownRenderer";

const gfmContent = `# GitHub Flavored Markdown (GFM) 機能のデモ

このページでは、remark-gfmプラグインによって有効化されるGFM拡張機能を紹介します。

## テーブル

| 機能 | 説明 | サポート状況 |
|------|------|------------|
| テーブル | Markdown形式のテーブル表示 | ✅ サポート |
| タスクリスト | チェックボックス付きリスト | ✅ サポート |
| 打ち消し線 | テキストに取り消し線 | ✅ サポート |
| 自動リンク | URLの自動リンク化 | ✅ サポート |

### 複雑なテーブルの例

| 左揃え | 中央揃え | 右揃え |
|:-------|:-------:|-------:|
| データ1 | データ2 | データ3 |
| 長いデータ項目 | 中央 | 100 |
| A | B | C |

## タスクリスト

プロジェクトの進捗状況：

- [x] プロジェクトのセットアップ
- [x] 基本コンポーネントの実装
- [x] GFM機能の統合
- [ ] テストの作成
- [ ] ドキュメントの整備
- [ ] デプロイ

## 打ち消し線

~~この機能は非推奨です。~~ 新しいAPIを使用してください。

価格: ~~¥5,000~~ **¥3,500** (30%オフ!)

## 自動リンク

URLは自動的にリンクになります：

- https://github.com
- https://react.dev
- https://vitejs.dev

メールアドレスも自動リンク化されます：
contact@example.com

## 脚注

GitHubのMarkdown[^1]は、標準のMarkdownを拡張したものです[^2]。

[^1]: GitHub Flavored Markdown (GFM)
[^2]: 詳細は https://github.github.com/gfm/ を参照してください

## 組み合わせ例

以下は複数の機能を組み合わせた例です：

| タスク | 状態 | 備考 |
|--------|------|------|
| ~~古いAPIの削除~~ | 完了 | https://github.com/example/pr/123 |
| 新しいAPIの実装 | 進行中 | - [x] 設計<br>- [ ] 実装 |
| ドキュメント更新 | 未着手 | contact@example.com に連絡 |
`;

/**
 * GitHub Flavored Markdown (GFM) 機能のデモページ
 */
export function GfmDemo() {
  return (
    <div className="demo-page">
      <MarkdownRenderer content={gfmContent} className="markdown-content" />
    </div>
  );
}
