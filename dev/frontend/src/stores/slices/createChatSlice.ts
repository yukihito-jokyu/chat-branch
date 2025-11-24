import type { StateCreator } from "zustand";
import type { StoreState, ChatSlice } from "../types";

export const createChatSlice: StateCreator<
  StoreState,
  [["zustand/devtools", never]],
  [],
  ChatSlice
> = (set) => ({
  messages: [],
  addMessage: (text, sender) =>
    set(
      (state) => ({
        messages: [
          ...state.messages,
          {
            id: crypto.randomUUID(),
            text,
            sender,
            timestamp: Date.now(),
          },
        ],
      }),
      false,
      "chat/addMessage"
    ),
});
