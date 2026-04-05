import { createPortal } from "react-dom";

export const Modal = ({
  open,
  title,
  body,
}: {
  open: boolean;
  title: string;
  body: React.ReactNode;
}) => {
  if (!open) return null;

  return createPortal(
    <div className="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm animate-fade-in">
      {/* Desktop Centered Modal */}
      <div className="hidden sm:flex h-full w-full items-center justify-center">
        <div className="w-full max-w-lg rounded-xl bg-bg-elevated border border-border shadow-lg animate-scale-in">
          <div className="flex items-center justify-between px-4 py-3.5 border-b border-border-muted">
            <h2 className="text-base font-semibold text-text-primary">
              {title}
            </h2>
          </div>
          {body}
        </div>
      </div>

      {/* Mobile Bottom Sheet */}
      <div className="md:hidden fixed inset-x-0 bottom-0 animate-slide-up">
        <div className="w-full rounded-t-2xl bg-bg-elevated border-t border-border shadow-lg max-h-[85vh] flex flex-col">
          {/* Drag Handle */}
          <div className="flex justify-center py-2">
            <div className="h-1.5 w-10 rounded-full bg-border-muted" />
          </div>

          {/* Header */}
          <div className="flex items-center justify-between px-4 py-3 border-b border-border-muted">
            <h2 className="text-base font-semibold text-text-primary">
              {title}
            </h2>
          </div>

          {/* Body (Scrollable) */}
          <div className="overflow-y-auto px-4 py-3">{body}</div>
        </div>
      </div>
    </div>,
    document.getElementById("modal")!,
  );
};
