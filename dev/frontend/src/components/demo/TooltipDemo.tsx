import {
  Tooltip,
  TooltipTrigger,
  TooltipContent,
} from "../floating-ui/tooltip";
import { Info } from "lucide-react";

export function TooltipDemo() {
  return (
    <div className="flex flex-col gap-4 rounded-lg border p-4 shadow-sm">
      <h3 className="text-lg font-semibold">Tooltip</h3>
      <p className="text-sm text-muted-foreground">
        Hover over the icon to see the tooltip.
      </p>
      <div className="flex items-center justify-center p-8">
        <Tooltip>
          <TooltipTrigger className="rounded-full p-2 hover:bg-accent">
            <Info className="h-5 w-5" />
          </TooltipTrigger>
          <TooltipContent>
            <p>This is a tooltip built with floating-ui/react</p>
          </TooltipContent>
        </Tooltip>
      </div>
    </div>
  );
}
