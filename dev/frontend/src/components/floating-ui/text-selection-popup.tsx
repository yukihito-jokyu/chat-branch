import * as React from "react";
import {
  useFloating,
  autoUpdate,
  offset,
  flip,
  shift,
  useDismiss,
  useRole,
  useInteractions,
  FloatingPortal,
  type Placement,
} from "@floating-ui/react";
import { cn } from "../../lib/utils";

interface TextSelectionPopupOptions {
  placement?: Placement;
  onAction?: (action: string, selectedText: string) => void;
}

export function useTextSelectionPopup({
  placement = "top",
  onAction,
}: TextSelectionPopupOptions = {}) {
  const [isOpen, setIsOpen] = React.useState(false);
  const [selectedText, setSelectedText] = React.useState("");
  const [virtualElement, setVirtualElement] = React.useState<{
    getBoundingClientRect: () => DOMRect;
  } | null>(null);

  const data = useFloating({
    placement,
    open: isOpen,
    onOpenChange: setIsOpen,
    whileElementsMounted: autoUpdate,
    middleware: [offset(10), flip(), shift({ padding: 5 })],
    elements: {
      reference: virtualElement as any,
    },
  });

  const context = data.context;

  const dismiss = useDismiss(context);
  const role = useRole(context, { role: "menu" });

  const interactions = useInteractions([dismiss, role]);

  const handleSelectionChange = React.useCallback(() => {
    const selection = window.getSelection();
    const text = selection?.toString().trim();

    if (text && text.length > 0) {
      const range = selection?.getRangeAt(0);
      if (range) {
        const rect = range.getBoundingClientRect();
        setVirtualElement({
          getBoundingClientRect: () => rect,
        });
        setSelectedText(text);
        setIsOpen(true);
      }
    } else {
      setIsOpen(false);
      setSelectedText("");
    }
  }, []);

  const handleAction = React.useCallback(
    (action: string) => {
      if (onAction && selectedText) {
        onAction(action, selectedText);
      }
      setIsOpen(false);
      window.getSelection()?.removeAllRanges();
    },
    [onAction, selectedText]
  );

  return React.useMemo(
    () => ({
      isOpen,
      setIsOpen,
      selectedText,
      handleSelectionChange,
      handleAction,
      ...interactions,
      ...data,
    }),
    [
      isOpen,
      selectedText,
      handleSelectionChange,
      handleAction,
      interactions,
      data,
    ]
  );
}

type ContextType = ReturnType<typeof useTextSelectionPopup> | null;

const TextSelectionPopupContext = React.createContext<ContextType>(null);

export const useTextSelectionPopupContext = () => {
  const context = React.useContext(TextSelectionPopupContext);

  if (context == null) {
    throw new Error(
      "TextSelectionPopup components must be wrapped in <TextSelectionPopup />"
    );
  }

  return context;
};

export function TextSelectionPopup({
  children,
  ...options
}: {
  children: React.ReactNode;
} & TextSelectionPopupOptions) {
  const popup = useTextSelectionPopup(options);

  React.useEffect(() => {
    document.addEventListener("selectionchange", popup.handleSelectionChange);
    return () => {
      document.removeEventListener(
        "selectionchange",
        popup.handleSelectionChange
      );
    };
  }, [popup.handleSelectionChange]);

  return (
    <TextSelectionPopupContext.Provider value={popup}>
      {children}
    </TextSelectionPopupContext.Provider>
  );
}

export const TextSelectionPopupContent = React.forwardRef<
  HTMLDivElement,
  React.HTMLProps<HTMLDivElement>
>(function TextSelectionPopupContent({ style, className, children, ...props }) {
  const context = useTextSelectionPopupContext();

  if (!context.isOpen) return null;

  return (
    <FloatingPortal>
      <div
        ref={context.refs.setFloating}
        style={{ ...context.floatingStyles, ...style }}
        className={cn(
          "z-50 min-w-[8rem] overflow-hidden rounded-md border bg-popover p-1 text-popover-foreground shadow-md animate-in fade-in-0 zoom-in-95",
          className
        )}
        {...(context.getFloatingProps(props) as any)}
      >
        {children}
      </div>
    </FloatingPortal>
  );
});

export const TextSelectionPopupItem = React.forwardRef<
  HTMLButtonElement,
  React.ButtonHTMLAttributes<HTMLButtonElement> & { action: string }
>(function TextSelectionPopupItem(
  { className, action, children, ...props },
  ref
) {
  const { handleAction } = useTextSelectionPopupContext();

  return (
    <button
      ref={ref}
      type="button"
      className={cn(
        "relative flex w-full cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none transition-colors hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground",
        className
      )}
      onClick={() => handleAction(action)}
      {...props}
    >
      {children}
    </button>
  );
});
