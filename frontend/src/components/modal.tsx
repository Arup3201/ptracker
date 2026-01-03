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
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60">
      <div className="w-110 rounded-md bg-(--bg-elevated) shadow-[0_4px_12px_rgba(0,0,0,0.6)]">
        <div className="flex items-center justify-between px-4 py-3 border-b border-(--border-muted)">
          <h2 className="text-sm font-semibold text-(--text-primary)">
            {title}
          </h2>
        </div>
        {body}
      </div>
    </div>,
    document.getElementById("modal")!
  );
};
