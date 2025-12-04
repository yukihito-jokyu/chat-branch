import { useCallback, useMemo, useEffect } from "react";
import {
  ReactFlow,
  MiniMap,
  Controls,
  Background,
  useNodesState,
  useEdgesState,
  addEdge,
  type Connection,
  type Edge,
  type Node,
} from "@xyflow/react";
import "@xyflow/react/dist/style.css";
import CustomNode from "./CustomNode";
import ChildLinkNode from "./ChildLinkNode";
import type { Message } from "../../chat/types";

type MapFlowProps = {
  messages: Message[];
  onNodeClick: (messageId: string) => void;
};

import { useQueries } from "@tanstack/react-query";
import { getChat } from "../../chat/api/chat";

export function MapFlow({ messages, onNodeClick }: MapFlowProps) {
  const [nodes, setNodes, onNodesChange] = useNodesState<Node>([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState<Edge>([]);

  const nodeTypes = useMemo(
    () => ({ custom: CustomNode, childLink: ChildLinkNode }),
    []
  );

  const forkChatIds = useMemo(() => {
    if (!messages) return [];
    return messages.flatMap((msg) => msg.forks?.map((f) => f.chat_uuid) || []);
  }, [messages]);

  const chatQueries = useQueries({
    queries: forkChatIds.map((chatId) => ({
      queryKey: ["chat", chatId],
      queryFn: () => getChat(chatId),
      staleTime: 1000 * 60 * 5,
    })),
  });

  const chatTitles = useMemo(() => {
    const titles: Record<string, string> = {};
    chatQueries.forEach((query) => {
      if (query.data) {
        titles[query.data.uuid] = query.data.title;
      }
    });
    return titles;
  }, [chatQueries]);

  useEffect(() => {
    if (!messages || messages.length === 0) return;

    const newNodes: Node[] = [];
    const newEdges: Edge[] = [];

    messages.forEach((msg, index) => {
      newNodes.push({
        id: msg.uuid,
        type: "custom",
        position: { x: 0, y: index * 200 },
        data: {
          label: msg.content,
          role: msg.role,
          isLatest: index === messages.length - 1,
        },
      });

      if (index < messages.length - 1) {
        newEdges.push({
          id: `e-${msg.uuid}-${messages[index + 1].uuid}`,
          source: msg.uuid,
          target: messages[index + 1].uuid,
          type: "smoothstep",
          animated: true,
          style: { stroke: "#94a3b8" },
        });
      }

      if (msg.forks && msg.forks.length > 0) {
        msg.forks.forEach((fork, forkIndex) => {
          const forkNodeId = `fork-${fork.chat_uuid}`;
          newNodes.push({
            id: forkNodeId,
            type: "childLink",
            position: { x: 300, y: index * 200 + forkIndex * 100 },
            data: {
              childChatId: fork.chat_uuid,
              title: chatTitles[fork.chat_uuid],
            },
          });

          newEdges.push({
            id: `e-${msg.uuid}-${forkNodeId}`,
            source: msg.uuid,
            target: forkNodeId,
            type: "smoothstep",
            animated: true,
            style: { stroke: "#60a5fa", strokeDasharray: "5,5" },
          });
        });
      }
    });

    setNodes(newNodes);
    setEdges(newEdges);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [messages, setNodes, setEdges, JSON.stringify(chatTitles)]);

  const onConnect = useCallback(
    (params: Connection) => setEdges((eds) => addEdge(params, eds)),
    [setEdges]
  );

  const handleNodeClick = useCallback(
    (_: React.MouseEvent, node: Node) => {
      onNodeClick(node.id);
    },
    [onNodeClick]
  );

  return (
    <div style={{ width: "100%", height: "100%" }} className="bg-gray-50">
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        nodeTypes={nodeTypes}
        onNodeClick={handleNodeClick}
        fitView
        minZoom={0.5}
        maxZoom={1.5}
      >
        <Controls />
        <MiniMap />
        <Background gap={20} size={1} color="#e2e8f0" />
      </ReactFlow>
    </div>
  );
}
