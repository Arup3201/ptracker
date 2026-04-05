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

      {/* Desktop Drawer (Right Side) */}
      <aside className="hidden sm:flex absolute right-0 top-0 h-full w-110 max-w-full bg-bg-elevated border-l border-border shadow-lg flex-col animate-slide-in-right">
        {/* Header */}
        <div className="flex items-center justify-between px-5 py-4 border-b border-border-muted shrink-0">
          <h2 className="text-base font-semibold text-text-primary truncate">
            {title}
          </h2>
          <button
            onClick={onClose}
            className="h-8 w-8 flex items-center justify-center rounded-md text-text-muted hover:text-text-primary hover:bg-bg-overlay transition"
          >
            <X size={18} />
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

      {/* Mobile Bottom Sheet */}
      <aside className="sm:hidden fixed inset-x-0 bottom-0 bg-bg-elevated rounded-t-2xl border-t border-border shadow-lg max-h-[90vh] flex flex-col animate-slide-up">
        {/* Drag Handle */}
        <div className="flex justify-center py-2">
          <div className="h-1.5 w-10 rounded-full bg-border-muted" />
        </div>

        {/* Header */}
        <div className="flex items-center justify-between px-4 py-3 border-b border-border-muted shrink-0">
          <h2 className="text-base font-semibold text-text-primary truncate">
            {title}
          </h2>
          <button
            onClick={onClose}
            className="h-9 w-9 flex items-center justify-center rounded-md text-text-muted hover:text-text-primary hover:bg-bg-overlay transition"
          >
            <X size={18} />
          </button>
        </div>

        {/* Content */}
        <div className="flex-1 overflow-y-auto px-4 py-3">{children}</div>

        {/* Footer */}
        {footer && (
          <div className="shrink-0 border-t border-border-muted px-4 py-3">
            {footer}
          </div>
        )}
      </aside>
    </div>
  );
}
