export type Project = {
  uuid: string;
  title: string;
  updated_at: string;
};

export type MessageRole = "user" | "assistant" | "system";

export type ForkResponse = {
  chat_uuid: string;
  selected_text: string;
  range_start: number;
  range_end: number;
};

export type Message = {
  uuid: string;
  role: MessageRole;
  content: string;
  forks: ForkResponse[];
  source_chat_uuid?: string;
  merge_reports: Message[];
};

export type Chat = {
  uuid: string;
  project_uuid: string;
  parent_uuid?: string;
  title: string;
  status: string;
  context_summary: string;
};

export type CreateProjectResponse = {
  project_uuid: string;
  chat_uuid: string;
  message_info: MessageInfo;
  updated_at: string;
};

export type MessageInfo = {
  message_uuid: string;
  message: string;
};

export type CreateProjectRequest = {
  initial_message: string;
};

export type ForkPreviewResponse = {
  suggested_title: string;
  generated_context: string;
};

export type MergePreviewResponse = {
  suggested_summary: string;
};

export type GetProjectResponse = {
  chat_uuid: string;
};

export type ForkChatResponse = {
  new_chat_id: string;
  message: string;
};

export type MergeChatRequest = {
  parent_chat_uuid: string;
  summary_content: string;
};

export type MergeChatResponse = {
  report_message_id: string;
  summary_content: string;
};

export type CloseChatResponse = {
  chat_uuid: string;
};

export type OpenChatResponse = {
  chat_uuid: string;
};
