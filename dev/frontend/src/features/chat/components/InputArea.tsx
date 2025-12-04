import { useState, type KeyboardEvent, useRef, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { Send } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { cn } from "@/lib/utils";

type InputAreaProps = {
  onSend: (text: string) => void;
  disabled?: boolean;
  className?: string;
};

export function InputArea({ onSend, disabled, className }: InputAreaProps) {
  const { t } = useTranslation("chat");
  const [text, setText] = useState("");
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  const handleSend = () => {
    if (text.trim() && !disabled) {
      onSend(text);
      setText("");
    }
  };

  const handleKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.nativeEvent.isComposing) return;
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  useEffect(() => {
    if (textareaRef.current) {
      textareaRef.current.style.height = "auto";
      textareaRef.current.style.height = `${textareaRef.current.scrollHeight}px`;
    }
  }, [text]);

  return (
    <div className={cn("p-4 border-t bg-background", className)}>
      <div className="relative flex items-end gap-2 max-w-3xl mx-auto">
        <Textarea
          ref={textareaRef}
          value={text}
          onChange={(e) => setText(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder={t("input_placeholder")}
          className="min-h-[50px] max-h-[200px] resize-none pr-12"
          disabled={disabled}
          rows={1}
        />
        <Button
          size="icon"
          className="absolute right-2 bottom-2 h-8 w-8"
          onClick={handleSend}
          disabled={!text.trim() || disabled}
        >
          <Send className="h-4 w-4" />
        </Button>
      </div>
    </div>
  );
}
