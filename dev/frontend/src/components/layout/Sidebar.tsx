import { ScrollArea } from "@/components/ui/scroll-area";
import { cn } from "@/lib/utils";

interface SidebarProps {
  currentView: string;
  onNavigate: (view: string) => void;
}

const sidebarNavItems = [
  {
    title: "Components",
    items: [
      { id: "button", label: "Button" },
      { id: "card", label: "Card" },
      { id: "input", label: "Input" },
      { id: "select", label: "Select" },
      { id: "dialog", label: "Dialog" },
      { id: "avatar", label: "Avatar" },
      { id: "toast", label: "Toast" },
    ],
  },
  {
    title: "Floating UI",
    items: [
      { id: "tooltip", label: "Tooltip" },
      { id: "popover", label: "Popover" },
      { id: "floating-dialog", label: "Dialog" },
      { id: "dropdown", label: "Dropdown" },
      { id: "text-selection", label: "Text Selection" },
    ],
  },
];

export function Sidebar({ currentView, onNavigate }: SidebarProps) {
  return (
    <aside className="fixed top-14 z-30 -ml-2 hidden h-[calc(100vh-3.5rem)] w-full shrink-0 md:sticky md:block">
      <ScrollArea className="h-full py-6 pr-6 lg:py-8">
        <div className="w-full">
          {sidebarNavItems.map((group, index) => (
            <div key={index} className="pb-4">
              <h4 className="mb-1 rounded-md px-2 py-1 text-sm font-semibold">
                {group.title}
              </h4>
              {group.items.map((item) => (
                <button
                  key={item.id}
                  onClick={() => onNavigate(item.id)}
                  className={cn(
                    "group flex w-full items-center rounded-md border border-transparent px-2 py-1 hover:underline",
                    currentView === item.id
                      ? "font-medium text-foreground"
                      : "text-muted-foreground"
                  )}
                >
                  {item.label}
                </button>
              ))}
            </div>
          ))}
        </div>
      </ScrollArea>
    </aside>
  );
}
