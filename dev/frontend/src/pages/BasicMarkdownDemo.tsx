import { MarkdownRenderer } from "../components/MarkdownRenderer";

const basicMarkdownContent = `# 基本的なMarkdown機能のデモ

このページでは、基本的なMarkdown記法の表示例を紹介します。

## 見出し

### レベル3の見出し
#### レベル4の見出し
##### レベル5の見出し
###### レベル6の見出し

## テキストの装飾

通常のテキストです。**太字のテキスト**や*斜体のテキスト*、そして***太字かつ斜体のテキスト***も表示できます。

## リスト

### 順序なしリスト

- 項目1
- 項目2
  - ネストされた項目2-1
  - ネストされた項目2-2
- 項目3

### 順序付きリスト

1. 最初の項目
2. 2番目の項目
3. 3番目の項目
   1. ネストされた項目3-1
   2. ネストされた項目3-2

## リンク

[React公式サイト](https://react.dev)にアクセスできます。

## 引用

> これは引用ブロックです。
> 複数行にわたる引用も可能です。
>
> 段落を分けることもできます。

## 水平線

以下は水平線です：

---

## インラインコード

JavaScriptでは \`const greeting = "Hello, World!";\` のようにコードを書けます。

## コードブロック

\`\`\`javascript
function greet(name) {
  return \`Hello, \${name}!\`;
}

console.log(greet("React"));
\`\`\`
`;

/**
 * 基本的なMarkdown機能のデモページ
 */
export function BasicMarkdownDemo() {
  return (
    <div className="demo-page">
      <MarkdownRenderer
        content={basicMarkdownContent}
        className="markdown-content"
      />
    </div>
  );
}
