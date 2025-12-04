import { apiClient } from "@/lib/api-client";
import type {
  Chat,
  CreateProjectRequest,
  CreateProjectResponse,
  ForkPreviewResponse,
  MergePreviewResponse,
  Message,
  Project,
  GetProjectResponse,
  ForkChatResponse,
  MergeChatRequest,
  MergeChatResponse,
  CloseChatResponse,
  OpenChatResponse,
} from "../types";

export const getProjects = async (): Promise<Project[]> => {
  return apiClient.get("/api/projects");
};

export const getProject = async (
  projectId: string
): Promise<GetProjectResponse> => {
  return apiClient.get(`/api/projects/${projectId}`);
};

export const createProject = async (
  data: CreateProjectRequest
): Promise<CreateProjectResponse> => {
  return apiClient.post("/api/projects", data);
};

export const getChat = async (chatId: string): Promise<Chat> => {
  return apiClient.get(`/api/chats/${chatId}`);
};

export const getMessages = async (chatId: string): Promise<Message[]> => {
  return apiClient.get(`/api/chats/${chatId}/messages`);
};

export const sendMessage = async (
  chatId: string,
  content: string
): Promise<Message> => {
  return apiClient.post(`/api/chats/${chatId}/message`, { content });
};

export const forkChat = async (
  chatId: string,
  data: {
    targetMessageId: string;
    parentChatId: string;
    selectedText: string;
    rangeStart: number;
    rangeEnd: number;
    title: string;
    contextSummary: string;
  }
): Promise<ForkChatResponse> => {
  return apiClient.post(`/api/chats/${chatId}/fork`, {
    target_message_uuid: data.targetMessageId,
    parent_chat_uuid: data.parentChatId,
    selected_text: data.selectedText,
    range_start: data.rangeStart,
    range_end: data.rangeEnd,
    title: data.title,
    context_summary: data.contextSummary,
  });
};

export const getForkPreview = async (
  chatId: string,
  data: {
    messageId: string;
    selectedText: string;
    rangeStart: number;
    rangeEnd: number;
  }
): Promise<ForkPreviewResponse> => {
  return apiClient.post(`/api/chats/${chatId}/fork/preview`, {
    target_message_uuid: data.messageId,
    selected_text: data.selectedText,
    range_start: data.rangeStart,
    range_end: data.rangeEnd,
  });
};

export const mergeChat = async (
  chatId: string,
  data: MergeChatRequest
): Promise<MergeChatResponse> => {
  return apiClient.post(`/api/chats/${chatId}/merge`, data);
};

export const getMergePreview = async (
  chatId: string
): Promise<MergePreviewResponse> => {
  return apiClient.post(`/api/chats/${chatId}/merge/preview`);
};

export const closeChat = async (chatId: string): Promise<CloseChatResponse> => {
  return apiClient.post(`/api/chats/${chatId}/close`);
};

export const openChat = async (chatId: string): Promise<OpenChatResponse> => {
  return apiClient.post(`/api/chats/${chatId}/open`);
};

export const getInitialChatStream = (
  chatId: string,
  callbacks: {
    onChunk: (chunk: string) => void;
    onDone: () => void;
    onError?: (error: Event) => void;
  }
): (() => void) => {
  const baseURL = apiClient.defaults.baseURL || "";
  const url = `${baseURL}/api/chats/${chatId}/stream`;

  const eventSource = new EventSource(url, {
    withCredentials: true,
  });

  eventSource.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      if (data.status === "processing" && data.chunk) {
        callbacks.onChunk(data.chunk);
      } else if (data.status === "done") {
        callbacks.onDone();
        eventSource.close();
      }
    } catch (e) {
      console.error("Failed to parse SSE message", e);
    }
  };

  eventSource.onerror = (error) => {
    if (callbacks.onError) {
      callbacks.onError(error);
    }
    eventSource.close();
  };

  return () => {
    eventSource.close();
  };
};

export const getMessagesStream = (
  chatId: string,
  callbacks: {
    onChunk: (chunk: string) => void;
    onDone: () => void;
    onError?: (error: Event) => void;
  }
): (() => void) => {
  const baseURL = apiClient.defaults.baseURL || "";
  const url = `${baseURL}/api/chats/${chatId}/messages/stream`;

  const eventSource = new EventSource(url, {
    withCredentials: true,
  });

  eventSource.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      if (data.status === "processing" && data.chunk) {
        callbacks.onChunk(data.chunk);
      } else if (data.status === "done") {
        callbacks.onDone();
        eventSource.close();
      }
    } catch (e) {
      console.error("Failed to parse SSE message", e);
    }
  };

  eventSource.onerror = (error) => {
    if (callbacks.onError) {
      callbacks.onError(error);
    }
    eventSource.close();
  };

  return () => {
    eventSource.close();
  };
};
