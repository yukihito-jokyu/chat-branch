import { useRef, useEffect, useMemo } from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { CodeBlock } from "./CodeBlock";
import { cn } from "@/lib/utils";
import type { Message } from "../types";
import { MergedReport } from "./MergedReport";
import { useChatStore } from "../stores/chatStore";
import { Link } from "@tanstack/react-router";

type MessageItemProps = {
  message: Message;
};

export function MessageItem({ message }: MessageItemProps) {
  const setSelection = useChatStore((state) => state.setSelection);
  const contentRef = useRef<HTMLDivElement>(null);

  const handleMouseUp = () => {
    const windowSelection = window.getSelection();
    if (!windowSelection || windowSelection.isCollapsed) {
      setSelection(null);
      return;
    }

    const text = windowSelection.toString().trim();
    if (!text) {
      return;
    }

    // 選択範囲がこのメッセージ内にあることを確認
    if (
      contentRef.current &&
      contentRef.current.contains(windowSelection.anchorNode)
    ) {
      // 選択範囲がマージレポート内にあるか確認
      const anchorNode = windowSelection.anchorNode;
      if (
        anchorNode &&
        anchorNode.parentElement &&
        anchorNode.parentElement.closest(".merged-report-content")
      ) {
        setSelection(null);
        return;
      }

      const range = windowSelection.getRangeAt(0);
      const rect = range.getBoundingClientRect();

      // コンテナのテキストコンテンツに対する相対的な範囲の開始と終了を計算
      const preSelectionRange = range.cloneRange();
      preSelectionRange.selectNodeContents(contentRef.current);
      preSelectionRange.setEnd(range.startContainer, range.startOffset);
      const start = preSelectionRange.toString().length;
      const end = start + text.length;

      setSelection({
        x: rect.left + rect.width / 2,
        y: rect.top,
        text,
        messageId: message.uuid,
        rangeStart: start,
        rangeEnd: end,
      });
    }
  };

  // 別の場所をクリックしたときに選択を解除
  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      const currentSelection = useChatStore.getState().selection;

      if (
        currentSelection &&
        currentSelection.messageId === message.uuid &&
        contentRef.current &&
        !contentRef.current.contains(e.target as Node)
      ) {
        const target = e.target as HTMLElement;
        if (target.closest(".deep-dive-menu")) return;

        setSelection(null);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, [message.uuid, setSelection]);

  const markdownComponents = useMemo(
    () => ({
      pre: ({ children }: any) => <div className="not-prose">{children}</div>,
      code: CodeBlock,
      a({ node, href, children, ...props }: any) {
        // リンクが内部チャットリンクかどうかを確認
        if (href && href.startsWith("/chat/")) {
          const chatId = href.split("/chat/")[1];
          return (
            <Link
              to="/chat/$chatId"
              params={{ chatId }}
              className="text-blue-500 hover:underline font-semibold"
              {...props}
            >
              {children}
            </Link>
          );
        }
        return (
          <a
            href={href}
            className="text-blue-500 hover:underline"
            target="_blank"
            rel="noopener noreferrer"
            {...props}
          >
            {children}
          </a>
        );
      },
      h1: ({ children }: any) => (
        <h1 className="text-2xl font-bold mt-6 mb-4">{children}</h1>
      ),
      h2: ({ children }: any) => (
        <h2 className="text-xl font-bold mt-5 mb-3">{children}</h2>
      ),
      h3: ({ children }: any) => (
        <h3 className="text-lg font-bold mt-4 mb-2">{children}</h3>
      ),
      h4: ({ children }: any) => (
        <h4 className="text-base font-bold mt-3 mb-2">{children}</h4>
      ),
      ul: ({ children }: any) => (
        <ul className="list-disc list-outside ml-6 mb-4 space-y-1">
          {children}
        </ul>
      ),
      ol: ({ children }: any) => (
        <ol className="list-decimal list-outside ml-6 mb-4 space-y-1">
          {children}
        </ol>
      ),
      li: ({ children }: any) => (
        <li className="leading-relaxed">{children}</li>
      ),
      p: ({ children }: any) => (
        <p className="mb-4 leading-relaxed">{children}</p>
      ),
      strong: ({ children }: any) => (
        <strong className="font-bold">{children}</strong>
      ),
      em: ({ children }: any) => <em className="italic">{children}</em>,
      blockquote: ({ children }: any) => (
        <blockquote className="border-l-4 border-gray-300 pl-4 py-1 my-4 bg-gray-50 dark:bg-gray-800/50 italic">
          {children}
        </blockquote>
      ),
    }),
    []
  );

  const remarkPlugins = useMemo(() => [remarkGfm], []);

  const isUser = message.role === "user";

  const displayContent = useMemo(() => {
    if (!message.forks || message.forks.length === 0) {
      return message.content;
    }

    // range_startでフォークをソートして順に処理
    const sortedForks = [...message.forks].sort(
      (a, b) => a.range_start - b.range_start
    );

    let result = "";
    let searchStartIndex = 0;
    const originalContent = message.content;

    for (const fork of sortedForks) {
      // 正確な位置決めのためにrange_startとrange_endを直接使用
      const start = fork.range_start;
      const end = fork.range_end;

      // フォーク範囲が重複しているか、現在の位置より前の場合はスキップ（ソート済みのフォークでは発生しないはず）
      if (start < searchStartIndex) continue;

      // 一致する前のテキストを追加
      result += originalContent.slice(searchStartIndex, start);

      // 完全一致を保証するために、元のコンテンツからリンクするテキストを取得
      const textToLink = originalContent.slice(start, end);

      // リンクを追加
      result += `[${textToLink}](/chat/${fork.chat_uuid})`;

      // searchStartIndexを更新
      searchStartIndex = end;
    }

    // 残りのテキストを追加
    result += originalContent.slice(searchStartIndex);

    return result;
  }, [message.content, message.forks]);

  return (
    <div
      id={`message-${message.uuid}`}
      className={cn(
        "flex gap-4 p-4 min-w-0",
        isUser ? "flex-row-reverse" : "flex-row"
      )}
    >
      <Avatar className="h-8 w-8">
        <AvatarImage src={isUser ? "/user-avatar.png" : "/ai-avatar.png"} />
        <AvatarFallback>{isUser ? "U" : "AI"}</AvatarFallback>
      </Avatar>

      <div
        className={cn(
          "flex flex-col max-w-[90%] md:max-w-[80%] lg:max-w-[70%] xl:max-w-4xl min-w-0",
          isUser ? "items-end" : "items-start"
        )}
      >
        <div className="flex items-center gap-2 mb-1">
          <span className="text-sm font-semibold">{isUser ? "You" : "AI"}</span>
        </div>

        <div
          ref={contentRef}
          onMouseUp={handleMouseUp}
          className={cn(
            "rounded-lg p-4 text-sm relative group selection:bg-yellow-300 selection:text-black overflow-hidden",
            isUser
              ? "bg-primary text-primary-foreground whitespace-pre-wrap"
              : "bg-muted/50"
          )}
        >
          <ReactMarkdown
            remarkPlugins={remarkPlugins}
            components={markdownComponents}
          >
            {displayContent}
          </ReactMarkdown>

          {message.merge_reports && message.merge_reports.length > 0 && (
            <div className="mt-4 space-y-2">
              {message.merge_reports.map((report) => (
                <MergedReport
                  key={report.uuid}
                  message={report}
                  chatId={message.forks?.[0]?.chat_uuid}
                />
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
