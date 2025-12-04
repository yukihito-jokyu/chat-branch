import { Handle, Position, type NodeProps } from "@xyflow/react";
import { memo } from "react";

const CustomNode = ({ data }: NodeProps) => {
  const isUser = data.role === "user";
  const isLatest = data.isLatest as boolean;

  return (
    <div
      className={`w-[250px] p-3 bg-white border rounded-lg shadow-sm flex flex-col gap-2 transition-colors hover:border-blue-400 cursor-pointer ${
        isLatest ? "border-blue-500 ring-2 ring-blue-200" : "border-gray-200"
      }`}
    >
      <Handle
        type="target"
        position={Position.Top}
        className="w-2 h-2 bg-gray-400 !bg-gray-400"
      />

      <div className="flex items-center gap-2 border-b pb-2 mb-1">
        <div
          className={`w-2 h-2 rounded-full ${
            isUser ? "bg-blue-500" : "bg-green-500"
          }`}
        />
        <span className="text-xs font-semibold text-gray-500 uppercase">
          {isUser ? "User" : "AI"}
        </span>
      </div>

      <div className="text-xs text-gray-700 line-clamp-4 leading-relaxed">
        {data.label as string}
      </div>

      <Handle
        type="source"
        position={Position.Bottom}
        className="w-2 h-2 bg-gray-400 !bg-gray-400"
      />
    </div>
  );
};

export default memo(CustomNode);
