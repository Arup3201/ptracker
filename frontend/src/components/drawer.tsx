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
    <div className="fixed inset-0 z-40 animate-fade-in">
      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-black/60 backdrop-blur-sm"
        onClick={onClose}
      />

      {/* Panel */}
      <aside className="absolute right-0 top-0 h-full w-110 bg-bg-elevated border-l border-border shadow-lg flex flex-col animate-slide-down">
        {/* Header */}
        <div className="flex items-center justify-between px-5 py-4 border-b border-border-muted shrink-0">
          <h2 className="text-base font-semibold text-text-primary tracking-snug truncate">
            {title}
          </h2>
          <button
            onClick={onClose}
            className="h-7 w-7 flex items-center justify-center rounded-md text-text-muted hover:text-text-primary hover:bg-bg-overlay transition duration-fast cursor-pointer"
          >
            <X size={16} />
          </button>
        </div>

        {/* Content */}
        <div className="flex-1 overflow-y-auto px-5 py-4">{children}</div>

        {/* Footer */}
        {footer && (
          <div className="shrink-0 border-t border-border-muted px-5 py-3">
            {footer}
          </div>
        )}
      </aside>
    </div>
  );
}
