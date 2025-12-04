import { createFileRoute } from "@tanstack/react-router";
import { InitChatArea } from "@/features/chat/components/InitChatArea";

export const Route = createFileRoute("/chat/")({
  component: ChatIndexPage,
});

function ChatIndexPage() {
  return (
    <div className="flex flex-col items-center justify-center h-full w-full p-4">
      <InitChatArea />
    </div>
  );
}
