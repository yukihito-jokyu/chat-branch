import { Button } from "@/components/ui/button";
import { ArrowUpRight, Lock, Unlock, Loader2 } from "lucide-react";
import { useTranslation } from "react-i18next";
import { ChatBreadcrumb } from "./ChatBreadcrumb";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { closeChat, openChat } from "../api/chat";
import { useChatStore } from "../stores/chatStore";

interface ChatHeaderProps {
  chat: {
    status: string;
    parent_uuid?: string | null;
  };
  chatId: string;
}

export function ChatHeader({ chat, chatId }: ChatHeaderProps) {
  const { t } = useTranslation("chat");
  const queryClient = useQueryClient();
  const setMergeModalOpen = useChatStore((state) => state.setMergeModalOpen);

  const closeMutation = useMutation({
    mutationFn: () => closeChat(chatId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["chat", chatId] });
    },
  });

  const openMutation = useMutation({
    mutationFn: () => openChat(chatId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["chat", chatId] });
    },
  });

  return (
    <div
      className={`flex items-center justify-between p-2 border-b ${
        chat.status === "closed"
          ? "bg-red-50"
          : chat.status === "merged"
          ? "bg-purple-50"
          : chat.status === "open"
          ? "bg-green-50"
          : "bg-muted/30"
      }`}
    >
      <div className="flex-1 min-w-0 mr-4">
        <ChatBreadcrumb chatId={chatId} />
      </div>
      <div className="flex gap-2 flex-shrink-0">
        <Button
          variant="outline"
          size="sm"
          onClick={() => setMergeModalOpen(true)}
          disabled={chat.status === "closed" || chat.status === "merged"}
        >
          <ArrowUpRight className="w-4 h-4 mr-2" />
          {chat.status === "merged" ? t("merged") : t("merge")}
        </Button>
        {chat.status === "closed" ? (
          <Button
            variant="outline"
            size="sm"
            onClick={() => openMutation.mutate()}
            disabled={openMutation.isPending}
          >
            {openMutation.isPending ? (
              <Loader2 className="w-4 h-4 mr-2 animate-spin" />
            ) : (
              <Unlock className="w-4 h-4 mr-2" />
            )}
            {t("open")}
          </Button>
        ) : (
          <Button
            variant="outline"
            size="sm"
            onClick={() => closeMutation.mutate()}
            disabled={closeMutation.isPending || chat.status === "merged"}
          >
            {closeMutation.isPending ? (
              <Loader2 className="w-4 h-4 mr-2 animate-spin" />
            ) : (
              <Lock className="w-4 h-4 mr-2" />
            )}
            {t("close")}
          </Button>
        )}
      </div>
    </div>
  );
}
