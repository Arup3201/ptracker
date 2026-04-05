import clsx from "clsx";
import type { ReactNode } from "react";

export function Table({ children }: { children: ReactNode }) {
  return (
    <div className="overflow-hidden rounded-lg border border-border bg-bg-surface shadow-sm">
      <div className="w-full overflow-x-auto">
        <table className="w-full min-w-150 table-fixed border-collapse text-sm">
          {children}
        </table>
      </div>
    </div>
  );
}

export function TableHeader({ children }: { children: ReactNode }) {
  return (
    <thead className="bg-bg-elevated sticky top-0 z-10 border-b border-border-muted">
      <tr>{children}</tr>
    </thead>
  );
}

export function TableBody({ children }: { children: ReactNode }) {
  return <tbody>{children}</tbody>;
}

export function TableRow({
  onClick = () => {},
  className = "",
  children,
}: {
  onClick?: () => void;
  className?: string;
  children: ReactNode;
}) {
  return (
    <tr
      className={clsx(
        "border-b border-border-muted last:border-0 transition duration-fast hover:bg-bg-overlay cursor-default",
        className,
      )}
      onClick={onClick}
    >
      {children}
    </tr>
  );
}

export function TableHead({
  children,
  align = "left",
  hideOnMobile = false,
}: {
  children: ReactNode;
  align?: "left" | "right" | "center";
  hideOnMobile?: boolean;
}) {
  return (
    <th
      className={clsx(
        "px-4 py-2.5 text-xs font-medium text-text-muted tracking-wide",
        align === "left" && "text-left",
        align === "right" && "text-right",
        align === "center" && "text-center",
        hideOnMobile && "hidden sm:table-cell",
      )}
    >
      {children}
    </th>
  );
}

export function TableCell({
  children,
  align = "left",
  muted,
  hideOnMobile = false,
}: {
  children: ReactNode;
  align?: "left" | "right" | "center";
  muted?: boolean;
  hideOnMobile?: boolean;
}) {
  return (
    <td
      className={clsx(
        "px-4 py-2.5 text-sm",
        align === "left" && "text-left",
        align === "right" && "text-right",
        align === "center" && "text-center",
        muted ? "text-text-secondary" : "text-text-primary",
        hideOnMobile && "hidden sm:table-cell",
      )}
    >
      {children}
    </td>
  );
}
