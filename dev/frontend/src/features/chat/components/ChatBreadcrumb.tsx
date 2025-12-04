import { Link } from "@tanstack/react-router";
import { ChevronRight, Loader2 } from "lucide-react";
import { useChatLineage } from "../hooks/useChatLineage";
import { Fragment } from "react";

type ChatBreadcrumbProps = {
  chatId: string;
};

export function ChatBreadcrumb({ chatId }: ChatBreadcrumbProps) {
  const { data: lineage, isLoading } = useChatLineage(chatId);

  if (isLoading) {
    return <Loader2 className="h-4 w-4 animate-spin text-muted-foreground" />;
  }

  if (!lineage || lineage.length === 0) {
    return null;
  }

  return (
    <div className="flex items-center text-sm text-muted-foreground overflow-hidden whitespace-nowrap">
      {lineage.map((chat, index) => {
        const isLast = index === lineage.length - 1;

        return (
          <Fragment key={chat.uuid}>
            {index > 0 && (
              <ChevronRight className="h-4 w-4 mx-1 flex-shrink-0" />
            )}
            {isLast ? (
              <span className="font-medium text-foreground truncate">
                {chat.title}
              </span>
            ) : (
              <Link
                to="/chat/$chatId"
                params={{ chatId: chat.uuid }}
                className="hover:text-foreground hover:underline truncate transition-colors"
              >
                {chat.title}
              </Link>
            )}
          </Fragment>
        );
      })}
    </div>
  );
}
