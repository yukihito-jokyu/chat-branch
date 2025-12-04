import { createFileRoute, Outlet } from "@tanstack/react-router";
import { ChatLayout } from "@/features/chat/components/ChatLayout";

export const Route = createFileRoute("/chat")({
  component: ChatPage,
});

function ChatPage() {
  return (
    <ChatLayout>
      <Outlet />
    </ChatLayout>
  );
}
