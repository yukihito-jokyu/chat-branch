export interface UserSlice {
  user: { id: string; name: string } | null;
  isLoggedIn: boolean;
  login: (name: string) => void;
  logout: () => void;
}

export interface ChatMessage {
  id: string;
  text: string;
  sender: string;
  timestamp: number;
}

export interface ChatSlice {
  messages: ChatMessage[];
  addMessage: (text: string, sender: string) => void;
}

export type StoreState = UserSlice & ChatSlice;
