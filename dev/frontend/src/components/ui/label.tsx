"use client";

import * as React from "react";
import * as LabelPrimitive from "@radix-ui/react-label";
import { cva, type VariantProps } from "class-variance-authority";

import { cn } from "@/lib/utils";

/**
 * labelVariants
 *
 * ラベルのスタイル定義です。
 * - text-sm: フォントサイズを小さく設定
 * - font-medium: フォントの太さを中程度に設定
 * - leading-none: 行送りをなしに設定
 * - peer-disabled: 関連付けられたコントロールが無効な場合のスタイル（カーソル変更、不透明度低下）
 */
const labelVariants = cva(
  "text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
);

/**
 * Labelコンポーネント
 *
 * Radix UIのLabelプリミティブをベースにしたラベルコンポーネントです。
 * フォームコントロールのラベルとして使用し、アクセシビリティを向上させます。
 *
 * 使用方法:
 * ```tsx
 * import { Label } from "@/components/ui/label"
 * import { Input } from "@/components/ui/input"
 *
 * // 基本的な使用法（Inputと関連付ける場合）
 * <div className="grid w-full max-w-sm items-center gap-1.5">
 *   <Label htmlFor="email">Email</Label>
 *   <Input type="email" id="email" placeholder="Email" />
 * </div>
 * ```
 */

const Label = React.forwardRef<
  React.ComponentRef<typeof LabelPrimitive.Root>,
  React.ComponentPropsWithoutRef<typeof LabelPrimitive.Root> &
    VariantProps<typeof labelVariants>
>(({ className, ...props }, ref) => (
  <LabelPrimitive.Root
    ref={ref}
    className={cn(labelVariants(), className)}
    {...props}
  />
));
Label.displayName = LabelPrimitive.Root.displayName;

export { Label };
