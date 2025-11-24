import {
  Dialog,
  DialogTrigger,
  DialogContent,
  DialogTitle,
  DialogDescription,
  DialogClose,
} from "../floating-ui/dialog";

export function DialogDemo() {
  return (
    <div className="flex flex-col gap-4 rounded-lg border p-4 shadow-sm">
      <h3 className="text-lg font-semibold">Dialog</h3>
      <p className="text-sm text-muted-foreground">
        Click the button to open the modal dialog.
      </p>
      <div className="flex items-center justify-center p-8">
        <Dialog>
          <DialogTrigger className="rounded-md bg-destructive px-4 py-2 text-destructive-foreground hover:bg-destructive/90">
            Delete Account
          </DialogTrigger>
          <DialogContent>
            <div className="flex flex-col gap-4">
              <div className="space-y-2">
                <DialogTitle>Are you absolutely sure?</DialogTitle>
                <DialogDescription>
                  This action cannot be undone. This will permanently delete
                  your account and remove your data from our servers.
                </DialogDescription>
              </div>
              <div className="flex justify-end gap-2">
                <DialogClose className="rounded-md border bg-background px-4 py-2 hover:bg-accent hover:text-accent-foreground">
                  Cancel
                </DialogClose>
                <DialogClose
                  className="rounded-md bg-destructive px-4 py-2 text-destructive-foreground hover:bg-destructive/90"
                  onClick={() => alert("Account deleted!")}
                >
                  Delete
                </DialogClose>
              </div>
            </div>
          </DialogContent>
        </Dialog>
      </div>
    </div>
  );
}
