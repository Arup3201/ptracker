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
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm animate-fade-in">
      <div className="w-110 rounded-xl bg-bg-elevated border border-border shadow-lg animate-scale-in">
        <div className="flex items-center justify-between px-4 py-3.5 border-b border-border-muted">
          <h2 className="text-base font-semibold text-text-primary tracking-snug">
            {title}
          </h2>
        </div>
        {body}
      </div>
    </div>,
    document.getElementById("modal")!,
  );
};
