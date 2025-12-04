import { useQuery } from "@tanstack/react-query";
import { getChat } from "../api/chat";
import type { Chat } from "../types";

export const useChatLineage = (chatId: string) => {
  return useQuery({
    queryKey: ["chatLineage", chatId],
    queryFn: async () => {
      const lineage: Chat[] = [];
      let currentChatId: string | undefined = chatId;

      while (currentChatId) {
        const chat: Chat = await getChat(currentChatId);
        lineage.unshift(chat);
        currentChatId = chat.parent_uuid;
      }

      return lineage;
    },
    enabled: !!chatId,
  });
};
