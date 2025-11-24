import * as React from "react";
import {
  useFloating,
  autoUpdate,
  offset,
  flip,
  shift,
  useClick,
  useDismiss,
  useRole,
  useInteractions,
  useMergeRefs,
  useListNavigation,
  useTypeahead,
  FloatingPortal,
  FloatingFocusManager,
  type Placement,
} from "@floating-ui/react";
import { cn } from "../../lib/utils";

interface DropdownMenuOptions {
  initialOpen?: boolean;
  placement?: Placement;
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
}

export function useDropdownMenu({
  initialOpen = false,
  placement = "bottom-start",
  open: controlledOpen,
  onOpenChange: setControlledOpen,
}: DropdownMenuOptions = {}) {
  const [uncontrolledOpen, setUncontrolledOpen] = React.useState(initialOpen);
  const [activeIndex, setActiveIndex] = React.useState<number | null>(null);

  const open = controlledOpen ?? uncontrolledOpen;
  const setOpen = setControlledOpen ?? setUncontrolledOpen;

  const listItemsRef = React.useRef<Array<HTMLButtonElement | null>>([]);
  const listContentRef = React.useRef<Array<string | null>>([]);

  const data = useFloating({
    placement,
    open,
    onOpenChange: setOpen,
    whileElementsMounted: autoUpdate,
    middleware: [
      offset(5),
      flip({
        padding: 5,
      }),
      shift({ padding: 5 }),
    ],
  });

  const context = data.context;

  const click = useClick(context, {
    event: "mousedown",
  });
  const dismiss = useDismiss(context);
  const role = useRole(context, { role: "menu" });

  const listNavigation = useListNavigation(context, {
    listRef: listItemsRef,
    activeIndex,
    onNavigate: setActiveIndex,
    loop: true,
  });

  const typeahead = useTypeahead(context, {
    listRef: listContentRef,
    activeIndex,
    onMatch: (index) => {
      if (open) {
        setActiveIndex(index);
      }
    },
  });

  const interactions = useInteractions([
    click,
    dismiss,
    role,
    listNavigation,
    typeahead,
  ]);

  return React.useMemo(
    () => ({
      open,
      setOpen,
      activeIndex,
      setActiveIndex,
      listItemsRef,
      listContentRef,
      ...interactions,
      ...data,
    }),
    [open, setOpen, activeIndex, interactions, data]
  );
}

type ContextType = ReturnType<typeof useDropdownMenu> | null;

const DropdownMenuContext = React.createContext<ContextType>(null);

export const useDropdownMenuContext = () => {
  const context = React.useContext(DropdownMenuContext);

  if (context == null) {
    throw new Error(
      "DropdownMenu components must be wrapped in <DropdownMenu />"
    );
  }

  return context;
};

export function DropdownMenu({
  children,
  ...options
}: {
  children: React.ReactNode;
} & DropdownMenuOptions) {
  const dropdownMenu = useDropdownMenu(options);
  return (
    <DropdownMenuContext.Provider value={dropdownMenu}>
      {children}
    </DropdownMenuContext.Provider>
  );
}

export const DropdownMenuTrigger = React.forwardRef<
  HTMLElement,
  React.HTMLProps<HTMLElement> & { asChild?: boolean }
>(function DropdownMenuTrigger(
  { children, asChild = false, ...props },
  propRef
) {
  const context = useDropdownMenuContext();
  const childrenRef = (children as any).ref;
  const ref = useMergeRefs([context.refs.setReference, propRef, childrenRef]);

  if (asChild && React.isValidElement(children)) {
    return React.cloneElement(
      children,
      context.getReferenceProps({
        ref,
        ...props,
        ...(children.props as any),
        "data-state": context.open ? "open" : "closed",
      })
    );
  }

  return (
    <button
      ref={ref}
      type="button"
      data-state={context.open ? "open" : "closed"}
      {...context.getReferenceProps(props)}
    >
      {children}
    </button>
  );
});

export const DropdownMenuContent = React.forwardRef<
  HTMLDivElement,
  React.HTMLProps<HTMLDivElement>
>(function DropdownMenuContent({ style, className, ...props }, propRef) {
  const context = useDropdownMenuContext();
  const ref = useMergeRefs([context.refs.setFloating, propRef]);

  if (!context.open) return null;

  return (
    <FloatingPortal>
      <FloatingFocusManager context={context.context} modal={false}>
        <div
          ref={ref}
          style={{ ...context.floatingStyles, ...style }}
          className={cn(
            "z-50 min-w-[8rem] overflow-hidden rounded-md border bg-popover p-1 text-popover-foreground shadow-md animate-in fade-in-0 zoom-in-95 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2",
            className
          )}
          {...(context.getFloatingProps(props) as any)}
        >
          {props.children}
        </div>
      </FloatingFocusManager>
    </FloatingPortal>
  );
});

export const DropdownMenuItem = React.forwardRef<
  HTMLButtonElement,
  React.ButtonHTMLAttributes<HTMLButtonElement>
>(function DropdownMenuItem({ className, ...props }, propRef) {
  const { activeIndex, setActiveIndex, listItemsRef, listContentRef, setOpen } =
    useDropdownMenuContext();

  const ref = React.useRef<HTMLButtonElement>(null);
  const mergedRef = useMergeRefs([ref, propRef]);

  // Register item for list navigation
  React.useLayoutEffect(() => {
    const index = listItemsRef.current.indexOf(null);
    if (index !== -1) {
      listItemsRef.current[index] = ref.current;
      listContentRef.current[index] = ref.current?.textContent ?? null;
      return () => {
        listItemsRef.current[index] = null;
        listContentRef.current[index] = null;
      };
    } else {
      listItemsRef.current.push(ref.current);
      listContentRef.current.push(ref.current?.textContent ?? null);
      return () => {
        const idx = listItemsRef.current.indexOf(ref.current);
        if (idx !== -1) {
          listItemsRef.current.splice(idx, 1);
          listContentRef.current.splice(idx, 1);
        }
      };
    }
  }, [listItemsRef, listContentRef]);

  const isActive =
    activeIndex !== null && listItemsRef.current[activeIndex] === ref.current;

  return (
    <button
      ref={mergedRef}
      type="button"
      role="menuitem"
      tabIndex={isActive ? 0 : -1}
      className={cn(
        "relative flex w-full cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none transition-colors hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50",
        isActive && "bg-accent text-accent-foreground",
        className
      )}
      {...props}
      onClick={(event) => {
        props.onClick?.(event);
        setOpen(false);
      }}
      onMouseEnter={() => {
        const index = listItemsRef.current.indexOf(ref.current);
        if (index !== -1) {
          setActiveIndex(index);
        }
      }}
      onMouseLeave={() => setActiveIndex(null)}
      onFocus={() => {
        const index = listItemsRef.current.indexOf(ref.current);
        if (index !== -1) {
          setActiveIndex(index);
        }
      }}
    >
      {props.children}
    </button>
  );
});
