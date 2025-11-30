import { createRouter } from "@tanstack/react-router";
import { routeTree } from "../routeTree.gen";

// 新しいルーターインスタンスを作成
export const router = createRouter({ routeTree });

// 型安全性のためにルーターインスタンスを登録
declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}
