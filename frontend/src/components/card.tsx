export function Card({ children }: { children: React.ReactNode }) {
  return (
    <div className="bg-(--bg-surface) border border-(--border-default) rounded-sm shadow-sm">
      {children}
    </div>
  );
}

export function CardContent({ children }: { children: React.ReactNode }) {
  return <div className="p-6">{children}</div>;
}
