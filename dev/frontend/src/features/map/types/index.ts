export type NodeData = {
  user_message: string;
  assistant: string;
};

export type NodePosition = {
  x: number;
  y: number;
};

export type MapNode = {
  id: string;
  chat_uuid: string;
  data: NodeData;
  position: NodePosition;
};

export type MapEdge = {
  id: string;
  source: string;
  target: string;
};

export type GetProjectTreeResponse = {
  nodes: MapNode[];
  edges: MapEdge[];
};
