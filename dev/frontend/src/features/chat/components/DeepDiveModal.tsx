import { useState } from "react";
import { useNavigate } from "@tanstack/react-router";
import { useMutation } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { Loader2 } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useChatStore } from "../stores/chatStore";
import { getForkPreview, forkChat } from "../api/chat";

export function DeepDiveModal() {
  const { t } = useTranslation("chat");
  const navigate = useNavigate();
  const deepDive = useChatStore((state) => state.deepDive);
  const setDeepDive = useChatStore((state) => state.setDeepDive);
  const resetDeepDive = useChatStore((state) => state.resetDeepDive);
  const currentChatId = useChatStore((state) => state.currentChatId);

  // 深掘りチャットの状態
  const [previewTitle, setPreviewTitle] = useState("");
  const [previewContext, setPreviewContext] = useState("");

  const previewMutation = useMutation({
    mutationFn: (args: {
      chatId: string;
      data: {
        messageId: string;
        selectedText: string;
        rangeStart: number;
        rangeEnd: number;
      };
    }) => getForkPreview(args.chatId, args.data),
    onSuccess: (data) => {
      setPreviewTitle(data.suggested_title);
      setPreviewContext(data.generated_context);
    },
  });

  const forkMutation = useMutation({
    mutationFn: (args: {
      chatId: string;
      data: {
        targetMessageId: string;
        parentChatId: string;
        selectedText: string;
        rangeStart: number;
        rangeEnd: number;
        title: string;
        contextSummary: string;
      };
    }) => forkChat(args.chatId, args.data),
    onSuccess: (data) => {
      setDeepDive({ ...deepDive, isOpen: false });
      navigate({ to: "/chat/$chatId", params: { chatId: data.new_chat_id } });

      // アニメーションの完了（200ms）を待ってから状態をリセット
      setTimeout(() => {
        resetDeepDive();
        setPreviewTitle("");
        setPreviewContext("");
      }, 300);
    },
  });

  const handleGeneratePreview = () => {
    if (!currentChatId) return;
    previewMutation.mutate({
      chatId: currentChatId,
      data: {
        messageId: deepDive.selectedMessageId,
        selectedText: deepDive.selectedText,
        rangeStart: deepDive.rangeStart,
        rangeEnd: deepDive.rangeEnd,
      },
    });
  };

  const handleConfirmFork = () => {
    if (!currentChatId) return;
    forkMutation.mutate({
      chatId: currentChatId,
      data: {
        targetMessageId: deepDive.selectedMessageId,
        parentChatId: currentChatId,
        selectedText: deepDive.selectedText,
        rangeStart: deepDive.rangeStart,
        rangeEnd: deepDive.rangeEnd,
        title: previewTitle,
        contextSummary: previewContext,
      },
    });
  };

  const handleOpenChange = (open: boolean) => {
    if (!open) {
      setDeepDive({ ...deepDive, isOpen: false });
      // アニメーションの完了（200ms）を待ってから状態をリセット
      setTimeout(() => {
        if (!useChatStore.getState().deepDive.isOpen) {
          resetDeepDive();
          setPreviewTitle("");
          setPreviewContext("");
        }
      }, 300);
    } else {
      setDeepDive({ ...deepDive, isOpen: open });
    }
  };

  return (
    <Dialog open={deepDive.isOpen} onOpenChange={handleOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>{t("create_deep_dive_title")}</DialogTitle>
          <DialogDescription>
            {t("create_deep_dive_description")}
          </DialogDescription>
        </DialogHeader>

        <div className="grid gap-4 py-4">
          <div className="grid gap-2">
            <Label>{t("selected_text_label")}</Label>
            <div className="max-h-[100px] overflow-y-auto rounded-md border bg-muted p-2 text-sm text-muted-foreground">
              {deepDive.selectedText}
            </div>
          </div>

          <div className="flex justify-end">
            <Button
              size="sm"
              variant="outline"
              onClick={handleGeneratePreview}
              disabled={previewMutation.isPending || !currentChatId}
            >
              {previewMutation.isPending && (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              )}
              {t("generate_preview")}
            </Button>
          </div>

          <div className="grid gap-2">
            <Label htmlFor="title">{t("chat_title_label")}</Label>
            <Input
              id="title"
              value={previewTitle}
              onChange={(e) => setPreviewTitle(e.target.value)}
              placeholder={t("auto_generated_placeholder")}
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="context">{t("context_summary_label")}</Label>
            <Textarea
              id="context"
              value={previewContext}
              onChange={(e) => setPreviewContext(e.target.value)}
              rows={4}
              placeholder={t("auto_generated_placeholder")}
            />
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => handleOpenChange(false)}>
            {t("cancel")}
          </Button>
          <Button
            onClick={handleConfirmFork}
            disabled={
              previewMutation.isPending ||
              forkMutation.isPending ||
              !previewTitle ||
              !currentChatId
            }
          >
            {forkMutation.isPending && (
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            )}
            {t("create_action")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
