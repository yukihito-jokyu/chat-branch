import { createFileRoute } from "@tanstack/react-router";
import { ChatArea } from "@/features/chat/components/ChatArea";

export const Route = createFileRoute("/chat/$chatId")({
  component: ChatDetailPage,
});

function ChatDetailPage() {
  const { chatId } = Route.useParams();
  return <ChatArea chatId={chatId} />;
}
