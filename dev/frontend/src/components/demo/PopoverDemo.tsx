import {
  Popover,
  PopoverTrigger,
  PopoverContent,
  PopoverClose,
} from "../floating-ui/popover";
import { Settings } from "lucide-react";

export function PopoverDemo() {
  return (
    <div className="flex flex-col gap-4 rounded-lg border p-4 shadow-sm">
      <h3 className="text-lg font-semibold">Popover</h3>
      <p className="text-sm text-muted-foreground">
        Click the button to toggle the popover.
      </p>
      <div className="flex items-center justify-center p-8">
        <Popover>
          <PopoverTrigger className="flex items-center gap-2 rounded-md bg-primary px-4 py-2 text-primary-foreground hover:bg-primary/90">
            <Settings className="h-4 w-4" />
            <span>Settings</span>
          </PopoverTrigger>
          <PopoverContent className="flex flex-col gap-4">
            <div className="space-y-2">
              <h4 className="font-medium leading-none">Dimensions</h4>
              <p className="text-sm text-muted-foreground">
                Set the dimensions for the layer.
              </p>
            </div>
            <div className="grid gap-2">
              <div className="grid grid-cols-3 items-center gap-4">
                <label htmlFor="width" className="text-sm font-medium">
                  Width
                </label>
                <input
                  id="width"
                  defaultValue="100%"
                  className="col-span-2 h-8 rounded-md border px-2 text-sm"
                />
              </div>
              <div className="grid grid-cols-3 items-center gap-4">
                <label htmlFor="maxWidth" className="text-sm font-medium">
                  Max. width
                </label>
                <input
                  id="maxWidth"
                  defaultValue="300px"
                  className="col-span-2 h-8 rounded-md border px-2 text-sm"
                />
              </div>
            </div>
            <PopoverClose className="rounded-md bg-secondary px-2 py-1 text-xs hover:bg-secondary/80">
              Close
            </PopoverClose>
          </PopoverContent>
        </Popover>
      </div>
    </div>
  );
}
