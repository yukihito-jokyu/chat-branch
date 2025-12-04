import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";
import { createProject } from "../api/chat";
import { InputArea } from "./InputArea";

import { useChatStore } from "../stores/chatStore";

export function InitChatArea() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const setPendingStreamChatId = useChatStore(
    (state) => state.setPendingStreamChatId
  );

  const createProjectMutation = useMutation({
    mutationFn: createProject,
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ["projects"] });
      setPendingStreamChatId(data.chat_uuid);
      navigate({
        to: "/chat/$chatId",
        params: { chatId: data.chat_uuid },
      });
    },
  });

  const handleSend = (text: string) => {
    createProjectMutation.mutate({ initial_message: text });
  };

  return (
    <div className="w-full max-w-3xl space-y-8">
      <h1 className="text-3xl font-bold text-center text-foreground/80">
        どのようなお手伝いができますか？
      </h1>
      <InputArea
        onSend={handleSend}
        disabled={createProjectMutation.isPending}
        className="border-t-0 bg-transparent p-0"
      />
    </div>
  );
}
