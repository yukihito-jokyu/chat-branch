import { useState, useEffect } from "react";
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
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { getMergePreview, mergeChat } from "../api/chat";

type MergeChatModalProps = {
  isOpen: boolean;
  onClose: () => void;
  chatId: string;
  parentChatId: string;
};

export function MergeChatModal({
  isOpen,
  onClose,
  chatId,
  parentChatId,
}: MergeChatModalProps) {
  const { t } = useTranslation("chat");
  const [summary, setSummary] = useState("");

  const previewMutation = useMutation({
    mutationFn: () => getMergePreview(chatId),
    onSuccess: (data) => {
      setSummary(data.suggested_summary);
    },
  });

  const mergeMutation = useMutation({
    mutationFn: (args: {
      chatId: string;
      summary: string;
      parentChatId: string;
    }) =>
      mergeChat(args.chatId, {
        parent_chat_uuid: args.parentChatId,
        summary_content: args.summary,
      }),
    onSuccess: () => {
      onClose();
    },
  });

  useEffect(() => {
    if (isOpen) {
      previewMutation.mutate();
    } else {
      setSummary("");
      previewMutation.reset();
      mergeMutation.reset();
    }
  }, [isOpen, chatId]);

  const handleMerge = () => {
    mergeMutation.mutate({ chatId, summary, parentChatId });
  };

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>{t("merge_to_parent_title")}</DialogTitle>
          <DialogDescription>
            {t("merge_to_parent_description")}
          </DialogDescription>
        </DialogHeader>

        <div className="grid gap-4 py-4">
          <div className="grid gap-2">
            <Label htmlFor="summary">{t("merge_summary_label")}</Label>
            {previewMutation.isPending ? (
              <div className="flex items-center justify-center h-[100px] border rounded-md bg-muted/50">
                <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
              </div>
            ) : (
              <Textarea
                id="summary"
                value={summary}
                onChange={(e) => setSummary(e.target.value)}
                rows={6}
                placeholder={t("merge_summary_placeholder")}
              />
            )}
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={onClose}>
            {t("cancel")}
          </Button>
          <Button
            onClick={handleMerge}
            disabled={
              previewMutation.isPending || mergeMutation.isPending || !summary
            }
          >
            {mergeMutation.isPending && (
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            )}
            {t("merge_action")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
