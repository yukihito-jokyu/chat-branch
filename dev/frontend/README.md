# devcontainer-react

## プロジェクト概要

このプロジェクトは、React と Vite を使用した Web アプリケーション開発のためのテンプレートです。
DevContainer を使用して、一貫した開発環境を迅速に構築できるように設計されています。

## 環境

### Docker コンテナ

- **Base Image**: `node:24-slim`
- **Shell**: `zsh` (oh-my-zsh 導入済み)
  - プラグイン: `git`, `zsh-autosuggestions`, `zsh-syntax-highlighting`
- **Package Manager**: `pnpm`
- **Tools**: `git`, `curl`, `wget`, `task`

### DevContainer 拡張機能

この環境には、以下の VS Code 拡張機能がプリインストールされています。

- **Prettier** (`esbenp.prettier-vscode`): コードフォーマッター
- **ESLint** (`dbaeumer.vscode-eslint`): JavaScript/TypeScript の静的解析
- **Indent Rainbow** (`oderwat.indent-rainbow`): インデントを色分けして可読性を向上
- **Code Spell Checker** (`streetsidesoftware.code-spell-checker`): スペルチェック
- **Git Graph** (`mhutchie.git-graph`): Git の履歴をグラフで表示
- **EditorConfig** (`editorconfig.editorconfig`): 異なるエディタ間でコーディングスタイルを統一
- **Trailing Spaces** (`shardulm94.trailing-spaces`): 行末の不要なスペースをハイライト
- **Japanese Language Pack** (`ms-ceintl.vscode-language-pack-ja`): VS Code の日本語化
- **ES7+ React/Redux/React-Native snippets** (`dsznajder.es7-react-js-snippets`): React 開発用のスニペット
- **Pretty TypeScript Errors** (`yoavbls.pretty-ts-errors`): TypeScript のエラーを見やすく表示

## 実行手順

### 開発サーバーの起動

```bash
pnpm dev
```

### ビルド

```bash
pnpm build
```

## task コマンド

### 拡張機能のインストール

```bash
task init
```
