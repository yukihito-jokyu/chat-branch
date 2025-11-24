import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
} from "../floating-ui/dropdown-menu";
import { Menu, User, CreditCard, Settings, LogOut } from "lucide-react";

export function DropdownMenuDemo() {
  return (
    <div className="flex flex-col gap-4 rounded-lg border p-4 shadow-sm">
      <h3 className="text-lg font-semibold">Dropdown Menu</h3>
      <p className="text-sm text-muted-foreground">
        Click the button to open the dropdown menu.
      </p>
      <div className="flex items-center justify-center p-8">
        <DropdownMenu>
          <DropdownMenuTrigger className="flex items-center gap-2 rounded-md border px-4 py-2 hover:bg-accent hover:text-accent-foreground">
            <Menu className="h-4 w-4" />
            <span>Open Menu</span>
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuItem onClick={() => console.log("Profile")}>
              <User className="mr-2 h-4 w-4" />
              <span>Profile</span>
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => console.log("Billing")}>
              <CreditCard className="mr-2 h-4 w-4" />
              <span>Billing</span>
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => console.log("Settings")}>
              <Settings className="mr-2 h-4 w-4" />
              <span>Settings</span>
            </DropdownMenuItem>
            <div className="my-1 h-px bg-muted" />
            <DropdownMenuItem
              className="text-destructive focus:text-destructive"
              onClick={() => console.log("Logout")}
            >
              <LogOut className="mr-2 h-4 w-4" />
              <span>Log out</span>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </div>
  );
}
