import * as React from "react";

import { cn } from "@/lib/utils";

/**
 * Inputコンポーネント
 *
 * 標準的なHTMLのinput要素をラップし、プロジェクトの標準スタイルを適用したコンポーネントです。
 * フォーム入力フィールドとして使用します。
 *
 * @example
 * // 基本的な使用方法
 * <Input type="email" placeholder="Email" />
 *
 * // ラベルと組み合わせた使用方法
 * <div className="grid w-full max-w-sm items-center gap-1.5">
 *   <Label htmlFor="email">Email</Label>
 *   <Input type="email" id="email" placeholder="Email" />
 * </div>
 *
 * // ファイルアップロード
 * <Input id="picture" type="file" />
 *
 * // 無効化状態
 * <Input disabled />
 *
 * @param {string} className - 追加のスタイルクラス
 * @param {string} type - inputのタイプ (text, email, password, fileなど)
 * @param {React.ComponentProps<"input">} props - その他の標準的なinput属性
 */
const Input = React.forwardRef<HTMLInputElement, React.ComponentProps<"input">>(
  ({ className, type, ...props }, ref) => {
    return (
      <input
        type={type}
        className={cn(
          "flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-base ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 md:text-sm",
          className
        )}
        ref={ref}
        {...props}
      />
    );
  }
);
Input.displayName = "Input";

export { Input };
