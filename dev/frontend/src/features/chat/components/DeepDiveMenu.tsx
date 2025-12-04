import { GitFork } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useTranslation } from "react-i18next";
import { useChatStore } from "../stores/chatStore";
import { useEffect } from "react";

export function DeepDiveMenu() {
  const { t } = useTranslation("chat");
  const selection = useChatStore((state) => state.selection);
  const setSelection = useChatStore((state) => state.setSelection);
  const setDeepDive = useChatStore((state) => state.setDeepDive);

  // スクロール時にメニューを閉じる
  useEffect(() => {
    const handleScroll = () => {
      setSelection(null);
    };

    // 子要素（ScrollAreaなど）のスクロールイベントを検出するために capture: true を使用
    window.addEventListener("scroll", handleScroll, { capture: true });
    return () => {
      window.removeEventListener("scroll", handleScroll, { capture: true });
    };
  }, [setSelection]);

  if (!selection) return null;

  const handleDeepDive = () => {
    setDeepDive({
      isOpen: true,
      selectedText: selection.text,
      selectedMessageId: selection.messageId,
      rangeStart: selection.rangeStart,
      rangeEnd: selection.rangeEnd,
    });
    setSelection(null);
    window.getSelection()?.removeAllRanges();
  };

  return (
    <div
      className="fixed z-50 animate-in fade-in zoom-in duration-200 deep-dive-menu"
      style={{
        left: selection.x,
        top: selection.y - 40,
        transform: "translateX(-50%)",
      }}
    >
      <Button
        size="sm"
        variant="secondary"
        className="shadow-lg gap-2 bg-background border hover:bg-accent"
        onClick={(e) => {
          e.stopPropagation();
          handleDeepDive();
        }}
        onMouseDown={(e) => e.preventDefault()}
      >
        <GitFork className="h-3 w-3" />
        {t("deep_dive")}
      </Button>
    </div>
  );
}
