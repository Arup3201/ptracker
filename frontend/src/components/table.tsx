import clsx from "clsx";
import type { ReactNode } from "react";

export function Table({ children }: { children: ReactNode }) {
  return (
    <div className="overflow-hidden rounded-sm border border-(--border-default) bg-(--bg-surface) shadow-[0_1px_2px_rgba(0,0,0,0.4)]">
      <div className="w-full">{children}</div>
    </div>
  );
}

export function TableHeader({ children }: { children: ReactNode }) {
  return (
    <table className="w-full table-fixed border-collapse text-sm">
      <thead className="border-b border-(--border-muted) text-(--text-muted)">
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
        "border-b border-(--border-muted) last:border-0 hover:bg-(--bg-elevated)",
        className
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
  return <th className={`px-3 py-2 text-${align} font-medium`}>{children}</th>;
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
      className={`px-3 py-2 text-${align} ${
        muted ? "text-(--text-secondary)" : ""
      }`}
    >
      {children}
    </td>
  );
}
