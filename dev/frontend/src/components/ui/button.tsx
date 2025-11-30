import * as React from "react";
import { Slot } from "@radix-ui/react-slot";
import { cva, type VariantProps } from "class-variance-authority";

import { cn } from "@/lib/utils";

// ボタンのスタイル定義 (cvaを使用)
const buttonVariants = cva(
  "inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0",
  {
    variants: {
      variant: {
        default: "bg-primary text-primary-foreground hover:bg-primary/90",
        destructive:
          "bg-destructive text-destructive-foreground hover:bg-destructive/90",
        outline:
          "border border-input bg-background hover:bg-accent hover:text-accent-foreground",
        secondary:
          "bg-secondary text-secondary-foreground hover:bg-secondary/80",
        ghost: "hover:bg-accent hover:text-accent-foreground",
        link: "text-primary underline-offset-4 hover:underline",
      },
      size: {
        default: "h-10 px-4 py-2",
        sm: "h-9 rounded-md px-3",
        lg: "h-11 rounded-md px-8",
        icon: "h-10 w-10",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  }
);

/**
 * Button Component
 *
 * Shadcn UIのButtonコンポーネントです。
 *
 * ## 使用方法
 * ```tsx
 * import { Button } from "@/components/ui/button"
 *
 * // 基本的な使用法
 * <Button>Click me</Button>
 *
 * // バリアントの指定
 * <Button variant="destructive">Delete</Button>
 *
 * // サイズの指定
 * <Button size="sm">Small</Button>
 *
 * // リンクとして振る舞う場合 (asChild)
 * <Button asChild>
 *   <a href="/login">Login</a>
 * </Button>
 * ```
 *
 * ## Props (設定)
 *
 * @prop {string} variant - ボタンのスタイルバリアントを指定します。
 *   - `default`: プライマリカラーの背景 (通常のアクション)
 *   - `destructive`: 破壊的なアクション (削除など) 用の赤色背景
 *   - `outline`: 枠線付き、背景なし (サブアクション)
 *   - `secondary`: セカンダリカラーの背景 (重要度の低いアクション)
 *   - `ghost`: 背景なし、ホバー時のみ背景表示 (ツールバーなど)
 *   - `link`: リンクのようなスタイル (アンダーライン)
 *
 * @prop {string} size - ボタンのサイズを指定します。
 *   - `default`: 標準サイズ (h-10 px-4 py-2)
 *   - `sm`: 小サイズ (h-9 rounded-md px-3)
 *   - `lg`: 大サイズ (h-11 rounded-md px-8)
 *   - `icon`: アイコン用正方形サイズ (h-10 w-10)
 *
 * @prop {boolean} asChild - trueの場合、デフォルトのbutton要素の代わりに子要素をレンダリングします。
 *   Radix UIのSlotを使用しており、子要素にスタイルやpropsを渡します。
 *   リンク(<a>タグ)や他のコンポーネントをボタンの見た目にしたい場合に使用します。
 */
export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  asChild?: boolean;
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, asChild = false, ...props }, ref) => {
    const Comp = asChild ? Slot : "button";
    return (
      <Comp
        className={cn(buttonVariants({ variant, size, className }))}
        ref={ref}
        {...props}
      />
    );
  }
);
Button.displayName = "Button";

export { Button, buttonVariants };
