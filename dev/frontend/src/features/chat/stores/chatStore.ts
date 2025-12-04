import { create } from "zustand";

type SelectionState = {
  x: number;
  y: number;
  text: string;
  messageId: string;
  rangeStart: number;
  rangeEnd: number;
} | null;

type DeepDiveState = {
  isOpen: boolean;
  selectedText: string;
  selectedMessageId: string;
  rangeStart: number;
  rangeEnd: number;
};

type ChatStore = {
  pendingStreamChatId: string | null;
  setPendingStreamChatId: (chatId: string | null) => void;
  selection: SelectionState;
  setSelection: (selection: SelectionState) => void;
  deepDive: DeepDiveState;
  setDeepDive: (deepDive: DeepDiveState) => void;
  resetDeepDive: () => void;
  currentChatId: string | null;
  setCurrentChatId: (chatId: string | null) => void;
  isMergeModalOpen: boolean;
  setMergeModalOpen: (isOpen: boolean) => void;
};

const initialDeepDiveState: DeepDiveState = {
  isOpen: false,
  selectedText: "",
  selectedMessageId: "",
  rangeStart: 0,
  rangeEnd: 0,
};

export const useChatStore = create<ChatStore>((set) => ({
  pendingStreamChatId: null,
  setPendingStreamChatId: (chatId) => set({ pendingStreamChatId: chatId }),
  selection: null,
  setSelection: (selection) => set({ selection }),
  deepDive: initialDeepDiveState,
  setDeepDive: (deepDive) => set({ deepDive }),
  resetDeepDive: () => set({ deepDive: initialDeepDiveState }),
  currentChatId: null,
  setCurrentChatId: (chatId) => set({ currentChatId: chatId }),
  isMergeModalOpen: false,
  setMergeModalOpen: (isOpen) => set({ isMergeModalOpen: isOpen }),
}));
