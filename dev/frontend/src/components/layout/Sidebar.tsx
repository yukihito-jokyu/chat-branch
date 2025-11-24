import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";

interface SidebarProps {
  currentView: string;
  onNavigate: (view: string) => void;
}

const components = [
  { name: "Button", id: "button" },
  { name: "Card", id: "card" },
  { name: "Input", id: "input" },
  { name: "Select", id: "select" },
  { name: "Dialog", id: "dialog" },
  { name: "Avatar", id: "avatar" },
  { name: "Toast", id: "toast" },
];

export function Sidebar({ currentView, onNavigate }: SidebarProps) {
  return (
    <aside className="fixed top-14 z-30 -ml-2 hidden h-[calc(100vh-3.5rem)] w-full shrink-0 md:sticky md:block">
      <ScrollArea className="h-full py-6 pr-6 lg:py-8">
        <div className="w-full">
          <div className="mb-4">
            <h4 className="mb-1 rounded-md px-2 py-1 text-sm font-semibold">
              Components
            </h4>
            <div className="grid grid-flow-row auto-rows-max text-sm">
              {components.map((item) => (
                <Button
                  key={item.id}
                  variant={currentView === item.id ? "secondary" : "ghost"}
                  className="w-full justify-start font-normal"
                  onClick={() => onNavigate(item.id)}
                >
                  {item.name}
                </Button>
              ))}
            </div>
          </div>
        </div>
      </ScrollArea>
    </aside>
  );
}
