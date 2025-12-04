import type { ReactNode } from "react";
import { Sidebar } from "./Sidebar";

type ChatLayoutProps = {
  children: ReactNode;
};

export function ChatLayout({ children }: ChatLayoutProps) {
  return (
    <div className="flex h-screen overflow-hidden bg-background">
      <aside className="w-64 border-r bg-muted/40 hidden md:block">
        <Sidebar />
      </aside>
      <main className="flex-1 flex flex-col min-w-0">{children}</main>
    </div>
  );
}
