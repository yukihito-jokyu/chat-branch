import { Handle, Position, type NodeProps } from "@xyflow/react";
import { memo, useCallback } from "react";
import { useNavigate } from "@tanstack/react-router";
import { useChatStore } from "../../chat/stores/chatStore";
import { ArrowRight } from "lucide-react";

const ChildLinkNode = ({ data }: NodeProps) => {
  const navigate = useNavigate();
  const setViewMode = useChatStore((state) => state.setViewMode);
  const childChatId = data.childChatId as string;

  const handleClick = useCallback(
    (e: React.MouseEvent) => {
      e.stopPropagation();
      setViewMode("map");
      navigate({ to: "/chat/$chatId", params: { chatId: childChatId } });
    },
    [navigate, setViewMode, childChatId]
  );

  return (
    <div
      className="w-[200px] p-3 bg-white border border-blue-200 rounded-lg shadow-sm flex items-center gap-2 cursor-pointer hover:bg-blue-50 transition-colors group"
      onClick={handleClick}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400 !bg-blue-400"
      />

      <div className="flex-1 min-w-0">
        <div className="text-xs font-medium text-blue-600 mb-0.5 truncate">
          {data.title ? (data.title as string) : "Loading..."}
        </div>
        <div className="text-xs text-gray-500 truncate">Click to view map</div>
      </div>

      <ArrowRight className="w-4 h-4 text-blue-400 group-hover:translate-x-0.5 transition-transform" />
    </div>
  );
};

export default memo(ChildLinkNode);
