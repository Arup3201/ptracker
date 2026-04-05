export function Card({ children }: { children: React.ReactNode }) {
  return (
    <div className="bg-bg-surface border border-border rounded-lg shadow-sm">
      {children}
    </div>
  );
}

export function CardContent({ children }: { children: React.ReactNode }) {
  return <div className="p-4">{children}</div>;
}
