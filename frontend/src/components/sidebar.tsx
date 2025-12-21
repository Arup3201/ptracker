import clsx from "clsx";
import { Logo } from "./logo";

export function Sidebar() {
  return (
    <aside className="flex w-60 flex-col border-r border-(--border-default) bg-(--bg-surface)">
      <div className="px-4 py-3 text-sm font-semibold self-center">
        <Logo />
      </div>

      <nav className="flex-1 px-2">
        <NavItem active>Dashboard</NavItem>
        <NavItem>Projects</NavItem>
      </nav>

      <div className="border-t border-(--border-default) p-2">
        <button className="w-full rounded-xs px-3 py-2 text-left text-sm text-(--danger) hover:bg-(--danger)/10">
          Logout
        </button>
      </div>
    </aside>
  );
}

function NavItem({ children, active }: { children: string; active?: boolean }) {
  return (
    <div
      className={clsx(
        "rounded-xs px-3 py-2 text-sm cursor-pointer",
        active
          ? "bg-(--bg-elevated) text-(--text-primary)"
          : "text-(--text-secondary) hover:bg-(--bg-elevated)"
      )}
    >
      {children}
    </div>
  );
}
