import { useEffect, useRef } from "react";
import { useQuery } from "@tanstack/react-query";
import { Loader2 } from "lucide-react";
import { MessageList } from "./MessageList";
import { InputArea } from "./InputArea";
import { useChatStream } from "../hooks/useChatStream";
import { getChat } from "../api/chat";
import { useChatStore } from "../stores/chatStore";
import { DeepDiveMenu } from "./DeepDiveMenu";
import { MergeChatModal } from "./MergeChatModal";
import { ChatHeader } from "./ChatHeader";
import { DeepDiveModal } from "./DeepDiveModal";
import { MapFlow } from "../../map/components/MapFlow";

type ChatAreaProps = {
  chatId: string;
};

export function ChatArea({ chatId }: ChatAreaProps) {
  const {
    messages,
    sendMessage,
    isStreaming,
    isLoading,
    connectToInitialStream,
  } = useChatStream(chatId);

  const { data: chat } = useQuery({
    queryKey: ["chat", chatId],
    queryFn: () => getChat(chatId),
  });

  const pendingStreamChatId = useChatStore(
    (state) => state.pendingStreamChatId
  );
  const setPendingStreamChatId = useChatStore(
    (state) => state.setPendingStreamChatId
  );
  const setCurrentChatId = useChatStore((state) => state.setCurrentChatId);
  const isMergeModalOpen = useChatStore((state) => state.isMergeModalOpen);
  const setMergeModalOpen = useChatStore((state) => state.setMergeModalOpen);

  const hasConnectedRef = useRef(false);

  useEffect(() => {
    setCurrentChatId(chatId);
    return () => setCurrentChatId(null);
  }, [chatId, setCurrentChatId]);

  useEffect(() => {
    if (pendingStreamChatId === chatId && !hasConnectedRef.current) {
      hasConnectedRef.current = true;
      connectToInitialStream();
      setPendingStreamChatId(null);
    }
  }, [
    chatId,
    pendingStreamChatId,
    connectToInitialStream,
    setPendingStreamChatId,
  ]);

  const viewMode = useChatStore((state) => state.viewMode);
  const setViewMode = useChatStore((state) => state.setViewMode);

  const handleNodeClick = (messageId: string) => {
    setViewMode("chat");
    setTimeout(() => {
      const element = document.getElementById(`message-${messageId}`);
      if (element) {
        element.scrollIntoView({ behavior: "smooth", block: "center" });
      }
    }, 100);
  };

  if (isLoading) {
    return (
      <div className="flex-1 flex items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full relative">
      <div className="flex-1 flex flex-col min-h-0">
        <ChatHeader
          chat={chat || { status: "open", parent_uuid: null }}
          chatId={chatId}
        />
        {viewMode === "chat" ? (
          <MessageList messages={messages} />
        ) : (
          <MapFlow messages={messages} onNodeClick={handleNodeClick} />
        )}
      </div>
      {viewMode === "chat" && (
        <InputArea
          onSend={sendMessage}
          disabled={isStreaming || chat?.status === "closed"}
        />
      )}

      <DeepDiveMenu />

      {chat && chat.parent_uuid && (
        <MergeChatModal
          isOpen={isMergeModalOpen}
          onClose={() => setMergeModalOpen(false)}
          chatId={chatId}
          parentChatId={chat.parent_uuid}
        />
      )}

      <DeepDiveModal />
    </div>
  );
}
