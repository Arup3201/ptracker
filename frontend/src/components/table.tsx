import clsx from "clsx";
import type { ReactNode } from "react";

export function Table({ children }: { children: ReactNode }) {
  return (
    <div className="overflow-hidden rounded-sm border border-(--border-default) bg-(--bg-surface) shadow-[0_1px_2px_rgba(0,0,0,0.4)]">
      <table className="w-full border-collapse text-sm">{children}</table>
    </div>
  );
}

export function TableHeader({ children }: { children: ReactNode }) {
  return (
    <thead className="border-b border-(--border-muted) text-(--text-muted)">
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
