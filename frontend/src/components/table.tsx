import clsx from "clsx";
import type { ReactNode } from "react";

export function Table({ children }: { children: ReactNode }) {
  return (
    <div className="overflow-hidden rounded-lg border border-border bg-bg-surface shadow-sm">
      <div className="w-full">{children}</div>
    </div>
  );
}

export function TableHeader({ children }: { children: ReactNode }) {
  return (
    <table className="w-full table-fixed border-collapse text-sm">
      <thead className="border-b border-border-muted bg-bg-elevated">
        <tr>{children}</tr>
      </thead>
    </table>
  );
}

export function TableBody({
  children,
  maxHeight = 500,
}: {
  children: ReactNode;
  maxHeight?: number | string;
}) {
  return (
    <div className="overflow-y-auto" style={{ maxHeight }}>
      <table className="w-full table-fixed border-collapse text-sm">
        <tbody>{children}</tbody>
      </table>
    </div>
  );
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
}: {
  children: ReactNode;
  align?: "left" | "right" | "center";
}) {
  return (
    <th
      className={clsx(
        "px-4 py-2.5 text-xs font-medium text-text-muted tracking-wide",
        align === "left" && "text-left",
        align === "right" && "text-right",
        align === "center" && "text-center",
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
}: {
  children: ReactNode;
  align?: "left" | "right" | "center";
  muted?: boolean;
}) {
  return (
    <td
      className={clsx(
        "px-4 py-2.5 text-sm",
        align === "left" && "text-left",
        align === "right" && "text-right",
        align === "center" && "text-center",
        muted ? "text-text-secondary" : "text-text-primary",
      )}
    >
      {children}
    </td>
  );
}
