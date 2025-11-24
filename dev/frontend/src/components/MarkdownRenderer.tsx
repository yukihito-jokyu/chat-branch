import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { CodeBlock } from "./CodeBlock";
import type { Components } from "react-markdown";

interface MarkdownRendererProps {
  content: string;
  className?: string;
}

/**
 * Markdownコンテンツをレンダリングする再利用可能なコンポーネント
 * react-markdownとremark-gfmを統合し、カスタムコンポーネントマッピングを提供
 */
export function MarkdownRenderer({
  content,
  className,
}: MarkdownRendererProps) {
  // カスタムコンポーネントマッピング
  const components: Components = {
    // コードブロックをCodeBlockコンポーネントに委譲
    code({ className, children }) {
      const childrenString = String(children);
      // インラインコードかどうかを判定（改行を含まない場合はインライン）
      const isInline = !childrenString.includes("\n");

      return (
        <CodeBlock className={className} inline={isInline}>
          {childrenString}
        </CodeBlock>
      );
    },
  };

  return (
    <div className={className}>
      <ReactMarkdown remarkPlugins={[remarkGfm]} components={components}>
        {content}
      </ReactMarkdown>
    </div>
  );
}
