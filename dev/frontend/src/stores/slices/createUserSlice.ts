import type { StateCreator } from "zustand";
import type { StoreState, UserSlice } from "../types";

export const createUserSlice: StateCreator<
  StoreState,
  [["zustand/devtools", never]],
  [],
  UserSlice
> = (set) => ({
  user: null,
  isLoggedIn: false,
  login: (name) =>
    set(
      { user: { id: crypto.randomUUID(), name }, isLoggedIn: true },
      false,
      "user/login"
    ),
  logout: () => set({ user: null, isLoggedIn: false }, false, "user/logout"),
});
