export function Section({
  title,
  children,
}: {
  title: string;
  children: React.ReactNode;
}) {
  return (
    <section className="space-y-2">
      <h2 className="text-lg font-semibold text-text-primary">{title}</h2>
      <div className="text-text-secondary leading-relaxed space-y-2">
        {children}
      </div>
    </section>
  );
}

export function Paragraph({ children }: { children: React.ReactNode }) {
  return <p>{children}</p>;
}

export function List({ items }: { items: React.ReactNode[] }) {
  return (
    <ul className="list-disc pl-5 space-y-1">
      {items.map((item, idx) => (
        <li key={idx}>{item}</li>
      ))}
    </ul>
  );
}
