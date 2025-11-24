import {
  TextSelectionPopup,
  TextSelectionPopupContent,
  TextSelectionPopupItem,
} from "../floating-ui/text-selection-popup";
import { Copy, Search, Share2 } from "lucide-react";
import { toast } from "sonner";

export function TextSelectionPopupDemo() {
  const handleAction = (action: string, selectedText: string) => {
    switch (action) {
      case "copy":
        navigator.clipboard.writeText(selectedText);
        toast.success("コピーしました");
        break;
      case "search":
        window.open(
          `https://www.google.com/search?q=${encodeURIComponent(selectedText)}`,
          "_blank"
        );
        break;
      case "share":
        toast.info(`共有: ${selectedText}`);
        break;
    }
  };

  return (
    <div className="flex flex-col gap-4 rounded-lg border p-4 shadow-sm">
      <h3 className="text-lg font-semibold">Text Selection Popup</h3>
      <p className="text-sm text-muted-foreground">
        テキストを選択すると、ポップアップメニューが表示されます。
      </p>
      <TextSelectionPopup onAction={handleAction}>
        <div className="rounded-lg border p-6">
          <p className="mb-4 text-base leading-relaxed">
            このテキストを選択してみてください。選択すると、コピー、検索、共有などのアクションを実行できるポップアップメニューが表示されます。
          </p>
          <p className="text-base leading-relaxed">
            floating-ui/react
            を使用することで、テキスト選択範囲に対して正確に配置されたポップアップを実装できます。
            これは、エディタやドキュメントビューアなどで便利な機能です。
          </p>
        </div>
        <TextSelectionPopupContent>
          <TextSelectionPopupItem action="copy">
            <Copy className="mr-2 h-4 w-4" />
            <span>コピー</span>
          </TextSelectionPopupItem>
          <TextSelectionPopupItem action="search">
            <Search className="mr-2 h-4 w-4" />
            <span>検索</span>
          </TextSelectionPopupItem>
          <TextSelectionPopupItem action="share">
            <Share2 className="mr-2 h-4 w-4" />
            <span>共有</span>
          </TextSelectionPopupItem>
        </TextSelectionPopupContent>
      </TextSelectionPopup>
    </div>
  );
}
