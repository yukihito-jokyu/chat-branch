import { useCallback, useEffect } from "react";
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
import { useQuery } from "@tanstack/react-query";
import { getProjectTree } from "../api/map";

type ProjectMapFlowProps = {
  projectId: string;
};

export function ProjectMapFlow({ projectId }: ProjectMapFlowProps) {
  const [nodes, setNodes, onNodesChange] = useNodesState<Node>([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState<Edge>([]);

  const { data } = useQuery({
    queryKey: ["projectTree", projectId],
    queryFn: () => getProjectTree(projectId),
  });

  useEffect(() => {
    if (data) {
      const flowNodes: Node[] = data.nodes.map((node) => ({
        id: node.id,
        position: node.position,
        data: { label: node.data.user_message },
        type: "default",
      }));

      const flowEdges: Edge[] = data.edges.map((edge) => ({
        id: edge.id,
        source: edge.source,
        target: edge.target,
        type: "smoothstep",
        animated: true,
      }));

      setNodes(flowNodes);
      setEdges(flowEdges);
    }
  }, [data, setNodes, setEdges]);

  const onConnect = useCallback(
    (params: Connection) => setEdges((eds) => addEdge(params, eds)),
    [setEdges]
  );

  return (
    <div style={{ width: "100%", height: "100%" }} className="bg-gray-50">
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        fitView
      >
        <Controls />
        <MiniMap />
        <Background gap={20} size={1} color="#e2e8f0" />
      </ReactFlow>
    </div>
  );
}
