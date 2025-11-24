import { create } from "zustand";
import { devtools } from "zustand/middleware";
import type { StoreState } from "./types";
import { createUserSlice } from "./slices/createUserSlice";
import { createChatSlice } from "./slices/createChatSlice";

export const useStore = create<StoreState>()(
  devtools(
    (...a) => ({
      ...createUserSlice(...a),
      ...createChatSlice(...a),
    }),
    { name: "ChatAppStore" }
  )
);
