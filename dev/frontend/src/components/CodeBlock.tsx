import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { vscDarkPlus } from "react-syntax-highlighter/dist/esm/styles/prism";

interface CodeBlockProps {
  children: string;
  className?: string;
  inline?: boolean;
}

/**
 * コードブロックのシンタックスハイライト表示コンポーネント
 * react-syntax-highlighterのPrismビルドを使用
 */
export function CodeBlock({ children, className, inline }: CodeBlockProps) {
  // 言語を className から抽出 (例: "language-javascript")
  const match = /language-(\w+)/.exec(className || "");
  const language = match ? match[1] : "";

  // インラインコードの場合は通常の code タグで表示
  if (inline) {
    return <code className={className}>{children}</code>;
  }

  // コードブロックの場合はシンタックスハイライトを適用
  return (
    <SyntaxHighlighter
      language={language || "text"}
      style={vscDarkPlus}
      showLineNumbers={true}
      customStyle={{
        margin: "1em 0",
        borderRadius: "8px",
        fontSize: "14px",
      }}
      lineNumberStyle={{
        minWidth: "3em",
        paddingRight: "1em",
        color: "#858585",
        userSelect: "none",
      }}
    >
      {String(children).replace(/\n$/, "")}
    </SyntaxHighlighter>
  );
}
