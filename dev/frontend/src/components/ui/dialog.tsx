import * as React from "react";
import * as DialogPrimitive from "@radix-ui/react-dialog";
import { X } from "lucide-react";

import { cn } from "@/lib/utils";

/**
 * Dialogコンポーネント
 *
 * モーダルダイアログを表示するためのルートコンポーネントです。
 * 状態管理（開閉状態）を行い、コンテキストを子コンポーネントに提供します。
 *
 * @example
 * <Dialog>
 *   <DialogTrigger>Open</DialogTrigger>
 *   <DialogContent>
 *     <DialogHeader>
 *       <DialogTitle>Title</DialogTitle>
 *       <DialogDescription>Description</DialogDescription>
 *     </DialogHeader>
 *     Content goes here
 *   </DialogContent>
 * </Dialog>
 */
const Dialog = DialogPrimitive.Root;

/**
 * DialogTriggerコンポーネント
 *
 * ダイアログを開くためのトリガーとなる要素です。
 * デフォルトではボタンとしてレンダリングされますが、asChildプロパティを使用することで
 * 子要素をトリガーとして機能させることができます。
 */
const DialogTrigger = DialogPrimitive.Trigger;

/**
 * DialogPortalコンポーネント
 *
 * ダイアログのコンテンツをDOMの別の場所（通常はbody直下）にレンダリングするために使用されます。
 * これにより、親要素のスタイル（overflow: hiddenなど）の影響を受けずに表示できます。
 */
const DialogPortal = DialogPrimitive.Portal;

/**
 * DialogCloseコンポーネント
 *
 * ダイアログを閉じるためのボタンとして機能するコンポーネントです。
 * このコンポーネントでラップされた要素をクリックするとダイアログが閉じます。
 */
const DialogClose = DialogPrimitive.Close;

/**
 * DialogOverlayコンポーネント
 *
 * ダイアログが開いているときに背景を覆うオーバーレイ（バックドロップ）です。
 * 背景をクリックするとダイアログが閉じる動作を提供します。
 */
const DialogOverlay = React.forwardRef<
  React.ComponentRef<typeof DialogPrimitive.Overlay>,
  React.ComponentPropsWithoutRef<typeof DialogPrimitive.Overlay>
>(({ className, ...props }, ref) => (
  <DialogPrimitive.Overlay
    ref={ref}
    className={cn(
      "fixed inset-0 z-50 bg-black/80 data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0",
      className
    )}
    {...props}
  />
));
DialogOverlay.displayName = DialogPrimitive.Overlay.displayName;

/**
 * DialogContentコンポーネント
 *
 * ダイアログのメインコンテンツを表示するコンテナです。
 * 画面中央に配置され、アニメーション効果を持ちます。
 * 右上に閉じるボタン（Xアイコン）が自動的に配置されます。
 */
const DialogContent = React.forwardRef<
  React.ComponentRef<typeof DialogPrimitive.Content>,
  React.ComponentPropsWithoutRef<typeof DialogPrimitive.Content>
>(({ className, children, ...props }, ref) => (
  <DialogPortal>
    <DialogOverlay />
    <DialogPrimitive.Content
      ref={ref}
      className={cn(
        "fixed left-[50%] top-[50%] z-50 grid w-full max-w-lg translate-x-[-50%] translate-y-[-50%] gap-4 border bg-background p-6 shadow-lg duration-200 data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[state=closed]:slide-out-to-left-1/2 data-[state=closed]:slide-out-to-top-[48%] data-[state=open]:slide-in-from-left-1/2 data-[state=open]:slide-in-from-top-[48%] sm:rounded-lg",
        className
      )}
      {...props}
    >
      {children}
      <DialogPrimitive.Close className="absolute right-4 top-4 rounded-sm opacity-70 ring-offset-background transition-opacity hover:opacity-100 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:pointer-events-none data-[state=open]:bg-accent data-[state=open]:text-muted-foreground">
        <X className="h-4 w-4" />
        <span className="sr-only">Close</span>
      </DialogPrimitive.Close>
    </DialogPrimitive.Content>
  </DialogPortal>
));
DialogContent.displayName = DialogPrimitive.Content.displayName;

/**
 * DialogHeaderコンポーネント
 *
 * ダイアログのヘッダーセクションです。
 * タイトルや説明文を適切なレイアウトで配置するために使用します。
 */
const DialogHeader = ({
  className,
  ...props
}: React.HTMLAttributes<HTMLDivElement>) => (
  <div
    className={cn(
      "flex flex-col space-y-1.5 text-center sm:text-left",
      className
    )}
    {...props}
  />
);
DialogHeader.displayName = "DialogHeader";

/**
 * DialogFooterコンポーネント
 *
 * ダイアログのフッターセクションです。
 * アクションボタン（保存、キャンセルなど）を右寄せなどで配置するために使用します。
 */
const DialogFooter = ({
  className,
  ...props
}: React.HTMLAttributes<HTMLDivElement>) => (
  <div
    className={cn(
      "flex flex-col-reverse sm:flex-row sm:justify-end sm:space-x-2",
      className
    )}
    {...props}
  />
);
DialogFooter.displayName = "DialogFooter";

/**
 * DialogTitleコンポーネント
 *
 * ダイアログのタイトルを表示します。
 * アクセシビリティのために重要な要素であり、スクリーンリーダーによって読み上げられます。
 */
const DialogTitle = React.forwardRef<
  React.ComponentRef<typeof DialogPrimitive.Title>,
  React.ComponentPropsWithoutRef<typeof DialogPrimitive.Title>
>(({ className, ...props }, ref) => (
  <DialogPrimitive.Title
    ref={ref}
    className={cn(
      "text-lg font-semibold leading-none tracking-tight",
      className
    )}
    {...props}
  />
));
DialogTitle.displayName = DialogPrimitive.Title.displayName;

/**
 * DialogDescriptionコンポーネント
 *
 * ダイアログの説明文を表示します。
 * タイトルの下に配置し、ダイアログの目的や詳細をユーザーに伝えます。
 */
const DialogDescription = React.forwardRef<
  React.ComponentRef<typeof DialogPrimitive.Description>,
  React.ComponentPropsWithoutRef<typeof DialogPrimitive.Description>
>(({ className, ...props }, ref) => (
  <DialogPrimitive.Description
    ref={ref}
    className={cn("text-sm text-muted-foreground", className)}
    {...props}
  />
));
DialogDescription.displayName = DialogPrimitive.Description.displayName;

export {
  Dialog,
  DialogPortal,
  DialogOverlay,
  DialogClose,
  DialogTrigger,
  DialogContent,
  DialogHeader,
  DialogFooter,
  DialogTitle,
  DialogDescription,
};
