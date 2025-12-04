import { useState, useEffect, useCallback, useMemo } from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import {
  getMessages,
  sendMessage as sendMessageApi,
  getMessagesStream,
  getInitialChatStream,
} from "../api/chat";
import type { Message } from "../types";

export function useChatStream(chatId: string) {
  const queryClient = useQueryClient();
  const [streamedMessage, setStreamedMessage] = useState<string>("");
  const [isStreaming, setIsStreaming] = useState(false);

  const { data: initialMessages, isLoading } = useQuery({
    queryKey: ["messages", chatId],
    queryFn: () => getMessages(chatId),
  });

  const [messages, setMessages] = useState<Message[]>([]);

  useEffect(() => {
    if (initialMessages) {
      setMessages(initialMessages);
    }
  }, [initialMessages]);

  const connectToInitialStream = useCallback(() => {
    setIsStreaming(true);
    setStreamedMessage("");

    getInitialChatStream(chatId, {
      onChunk: (chunk) => {
        setStreamedMessage((prev) => prev + chunk);
      },
      onDone: () => {
        setIsStreaming(false);
        queryClient.invalidateQueries({ queryKey: ["messages", chatId] });
        setStreamedMessage("");
      },
      onError: (error) => {
        console.error("EventSource failed:", error);
        setIsStreaming(false);
      },
    });
  }, [chatId, queryClient]);

  const sendMessage = useCallback(
    async (content: string) => {
      const tempUserMessage: Message = {
        uuid: Date.now().toString(),
        role: "user",
        content,
        forks: [],
        merge_reports: [],
      };
      setMessages((prev) => [...prev, tempUserMessage]);

      try {
        await sendMessageApi(chatId, content);
        setIsStreaming(true);
        setStreamedMessage("");

        getMessagesStream(chatId, {
          onChunk: (chunk) => {
            setStreamedMessage((prev) => prev + chunk);
          },
          onDone: () => {
            setIsStreaming(false);
            queryClient.invalidateQueries({ queryKey: ["messages", chatId] });
            setStreamedMessage("");
          },
          onError: (error) => {
            console.error("EventSource failed:", error);
            setIsStreaming(false);
          },
        });
      } catch (error) {
        console.error("Failed to send message:", error);
        setIsStreaming(false);
      }
    },
    [chatId, queryClient]
  );

  const displayMessages = useMemo(() => {
    const msgs = [...messages];
    if (isStreaming && streamedMessage) {
      msgs.push({
        uuid: "streaming-ai",
        role: "assistant",
        content: streamedMessage,
        forks: [],
        merge_reports: [],
      });
    }
    return msgs;
  }, [messages, isStreaming, streamedMessage]);

  return {
    messages: displayMessages,
    sendMessage,
    connectToInitialStream,
    isStreaming,
    isLoading,
  };
}
