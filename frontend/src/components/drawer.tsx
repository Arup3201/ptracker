import { X } from "lucide-react";
import { type ReactNode, useEffect } from "react";

type DrawerProps = {
  open: boolean;
  title?: string;
  onClose: () => void;
  footer?: ReactNode;
  children: ReactNode;
};

export function Drawer({
  open,
  title,
  onClose,
  footer,
  children,
}: DrawerProps) {
  useEffect(() => {
    if (!open) return;

    const onKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        onClose();
      }
    };

    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, [open, onClose]);

  if (!open) return null;

  return (
    <div className="fixed inset-0 z-40">
      {/* Backdrop */}
      <div className="absolute inset-0 bg-black/40" onClick={onClose} />

      {/* Panel */}
      <aside className="absolute right-0 top-0 h-full w-110 bg-(--bg-elevated) border-l border-(--border-default) shadow-lg flex flex-col">
        {/* Header */}
        <div className="flex items-center justify-between px-4 py-3 border-b border-(--border-default)">
          <h2 className="text-sm font-medium text-(--text-primary) truncate">
            {title}
          </h2>

          <button
            onClick={onClose}
            className="text-(--text-muted) hover:text-(--text-primary) cursor-pointer"
          >
            <X size={28} />
          </button>
        </div>

        {/* Content */}
        <div className="flex-1 overflow-y-auto px-4 py-4">{children}</div>

        {/* Footer */}
        {footer && (
          <div className="border-t border-(--border-muted) px-4 py-3">
            {footer}
          </div>
        )}
      </aside>
    </div>
  );
}
