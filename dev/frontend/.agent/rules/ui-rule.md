---
trigger: model_decision
description: shadcn/uiとTailwind CSSを活用したUI作成時のルール
---

# 🎨 UI Implementation Guidelines (shadcn/ui + Tailwind CSS)

このガイドラインは、UIの一貫性、可読性、保守性を最大化するための絶対ルールです。  
Tailwind CSSの柔軟性が引き起こす「カオス（Utility Soup）」を防ぎ、shadcn/uiのデザインシステムを厳格に適用することを目的とします。

---

## 1. 🎯 Design Tokens & Semantic Colors（最重要）

**原則: 色と角丸（Radius）のハードコーディングを禁止し、必ずSemantic Tokenを使用する。**  
これにより、ダークモード対応とテーマ変更が自動的に機能します。

### ❌ 禁止（Strictly Prohibited）

以下のユーティリティは使用禁止です：

- Raw Colors: `bg-blue-500`, `text-gray-600`, `border-red-400`
- Raw Radius: `rounded-lg`, `rounded-md`, `rounded-[8px]`
- Dark Mode Prefix: `dark:bg-gray-900`（セマンティックトークンを使えば不要）

### ✅ 必須（Mandatory）

必ず以下のセマンティックトークンを使用してください。

| カテゴリ | クラス名 | 用途 |
|----------|-----------|------|
| 背景 | `bg-background` | ページ全体の背景 |
| 背景 | `bg-card` | カード、パネルの背景 |
| 背景 | `bg-muted` | 二次的な背景、無効状態 |
| 背景 | `bg-primary` | メインアクション（ボタン等） |
| 背景 | `bg-accent` | ホバー時のハイライト |
| テキスト | `text-foreground` | 基本テキスト |
| テキスト | `text-muted-foreground` | 補足、サブテキスト |
| テキスト | `text-primary-foreground` | Primary背景上のテキスト |
| ボーダー | `border-border` | 基本的な枠線 |
| ボーダー | `border-input` | 入力フォームの枠線 |
| 形状 | `rounded-[var(--radius)]` | 角丸（コンポーネント標準に従う） |

---

## 2. 🧩 Component Architecture & Utilities

**原則: div による再発明を避け、shadcn/ui コンポーネントを優先する。**

### 2.1 コンポーネント優先

可能な限り `div + className` で実装せず、UIコンポーネントを使用してください。

- ❌ `<div className="border p-4 rounded...">`
- ✅ `<Card><CardContent>...</CardContent></Card>`

### 2.2 クラスの結合（Class Merging）

コンポーネントにクラスを追加する場合、文字列結合は禁止です。  
必ず `cn()` ユーティリティを使用してください。

```ts
import { cn } from "@/lib/utils";

// ❌ 禁止
<div className={`flex ${className}`}>

// ✅ 必須（Tailwindの競合を自動解決）
<div className={cn("flex items-center gap-2", className)}>

## 3. ✍ Typography System

原則: サイズだけでなく、行間（leading）と字間（tracking）もセットで定義する。  
単なる `text-lg` 等の使用は避け、以下の階層ルールに従ってください。

| 要素 | 推奨クラス構成 | 備考 |
|------|----------------|------|
| H1 | `text-3xl font-bold tracking-tight lg:text-4xl` | タイトルは詰める |
| H2 | `text-2xl font-semibold tracking-tight` | セクション見出し |
| H3 | `text-xl font-semibold tracking-tight` | カードタイトル等 |
| Body | `text-base leading-7` | 重要: 行間を確保する |
| Small | `text-sm font-medium leading-none` | 補足情報 |
| Muted | `text-sm text-muted-foreground` | 色を薄くする場合 |

---

## 4. 📐 Layout & Spacing

原則: 4pxベースのスケールを使用し、PCファーストで記述する。

### 4.1 Desktop First Strategy (PC基準)

スタイルは「PC（デスクトップ）」を基準に書き、スマホ・タブレット向けを `max-md:` や `max-lg:` 等の max-width 修飾子で上書きします。

- ❌ `grid grid-cols-1 md:grid-cols-2`（スマホ基準：デフォルト1列 → PCで2列）
- ✅ `grid grid-cols-2 max-md:grid-cols-1`（PC基準：デフォルト2列 → スマホで1列）

### 4.2 Spacing Scale

| 用途 | クラス |
|------|--------|
| フォーム要素間の隙間 | `gap-4` |
| コンポーネント間の隙間 | `gap-6` |
| カード内パディング | `p-6` |
| セクション上下余白 | `py-16` (Desktop) / `max-md:py-12` (Mobile) |
| コンテンツ幅 | `max-w-screen-xl mx-auto px-8 max-md:px-4` |

---

## 5. 🖼 Icons (Lucide React)

原則: アイコンは lucide-react を使用し、サイズと色を統一する。

- **Import:**  
  `import { Search } from "lucide-react"`

- **Size**
  - UI内: `h-4 w-4`
  - 強調/ボタン内: `h-5 w-5`

- **Color**
  - 基本: 親の `text-foreground` を継承
  - 補足: `text-muted-foreground`

---

## 6. 🧭 Interactivity & Accessibility

原則: FocusリングとHoverステートを必ず実装する。

### Focus
- shadcn/ui は自動対応
- `focus-visible:ring-2 focus-visible:ring-ring`

### Hover
- ❌ `hover:bg-gray-100`
- ✅ `hover:bg-accent hover:text-accent-foreground`

---

## 7. 🧹 Developer Experience (DX)

原則: コードフォーマットを自動化し、属人性を排除する。

### 7.1 Class Sorting

`prettier-plugin-tailwindcss` を導入し、保存時にクラスを自動ソートしてください。  
手動での並び替えは禁止です。

### 7.2 Magic Numbers の禁止

デザインシステム外の数値を直接書くことを禁止します。

- ❌ `w-[350px]`, `mt-[13px]`, `z-[100]`
- ✅ `w-72`, `mt-4`, `z-50`
