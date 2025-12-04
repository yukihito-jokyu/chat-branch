import { Link } from "@tanstack/react-router";
import { GitMerge, ArrowRight } from "lucide-react";
import { useTranslation } from "react-i18next";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { CodeBlock } from "./CodeBlock";
import type { Message } from "../types";

type MergedReportProps = {
  message: Message;
  chatId?: string;
};

export function MergedReport({ message, chatId }: MergedReportProps) {
  const { t } = useTranslation("chat");
  const targetChatId = chatId || message.source_chat_uuid;

  return (
    <Card className="my-4 border-l-4 border-l-purple-500 bg-purple-50/50 dark:bg-purple-900/10 merged-report-content">
      <CardHeader className="pb-2">
        <CardTitle className="flex items-center gap-2 text-base font-medium text-purple-700 dark:text-purple-300">
          <GitMerge className="h-4 w-4" />
          {targetChatId ? (
            <Link
              to="/chat/$chatId"
              params={{ chatId: targetChatId }}
              className="hover:underline underline decoration-purple-500/30 underline-offset-4"
            >
              {t("merged_report_title")}
            </Link>
          ) : (
            t("merged_report_title")
          )}
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="text-sm text-muted-foreground mb-3">
          <ReactMarkdown
            remarkPlugins={[remarkGfm]}
            components={{
              pre: ({ children }: any) => (
                <div className="not-prose">{children}</div>
              ),
              code: CodeBlock,
              a({ node, href, children, ...props }: any) {
                if (href && href.startsWith("/chat/")) {
                  const chatId = href.split("/chat/")[1];
                  return (
                    <Link
                      to="/chat/$chatId"
                      params={{ chatId }}
                      className="text-purple-600 hover:underline font-semibold dark:text-purple-400"
                      {...props}
                    >
                      {children}
                    </Link>
                  );
                }
                return (
                  <a
                    href={href}
                    className="text-purple-600 hover:underline dark:text-purple-400"
                    target="_blank"
                    rel="noopener noreferrer"
                    {...props}
                  >
                    {children}
                  </a>
                );
              },
              h1: ({ children }: any) => (
                <h1 className="text-lg font-bold mt-4 mb-2">{children}</h1>
              ),
              h2: ({ children }: any) => (
                <h2 className="text-base font-bold mt-3 mb-2">{children}</h2>
              ),
              h3: ({ children }: any) => (
                <h3 className="text-sm font-bold mt-2 mb-1">{children}</h3>
              ),
              ul: ({ children }: any) => (
                <ul className="list-disc list-outside ml-4 mb-2 space-y-1">
                  {children}
                </ul>
              ),
              ol: ({ children }: any) => (
                <ol className="list-decimal list-outside ml-4 mb-2 space-y-1">
                  {children}
                </ol>
              ),
              li: ({ children }: any) => (
                <li className="leading-relaxed">{children}</li>
              ),
              p: ({ children }: any) => (
                <p className="mb-2 leading-relaxed last:mb-0">{children}</p>
              ),
              blockquote: ({ children }: any) => (
                <blockquote className="border-l-4 border-purple-300 pl-4 py-1 my-2 bg-purple-100/50 dark:bg-purple-900/20 italic">
                  {children}
                </blockquote>
              ),
            }}
          >
            {message.content}
          </ReactMarkdown>
        </div>
        {targetChatId && (
          <Link
            to="/chat/$chatId"
            params={{ chatId: targetChatId }}
            className="inline-flex items-center text-xs font-medium text-purple-600 hover:underline dark:text-purple-400"
          >
            {t("check_original_chat")} <ArrowRight className="ml-1 h-3 w-3" />
          </Link>
        )}
      </CardContent>
    </Card>
  );
}
