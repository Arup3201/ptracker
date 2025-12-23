import clsx from "clsx";
import { useNavigate } from "react-router";
import { Logo } from "./logo";
import { Button } from "./button";

const NavItems = [
  {
    name: "Dashboard",
    path: "/",
  },
  {
    name: "Projects",
    path: "/projects",
  },
  {
    name: "Explore",
    path: "/explore",
  },
];

export function Sidebar() {
  const navigate = useNavigate();
  let active: number;
  NavItems.forEach((navItem, i) => {
    if (window.location.pathname === navItem.path) {
      active = i;
    }
  });

  return (
    <aside className="flex w-60 flex-col border-r border-(--border-default) bg-(--bg-surface)">
      <div className="px-4 py-3 text-sm font-semibold self-center">
        <Logo />
      </div>

      <nav className="flex-1 px-2">
        {NavItems.map((navitem, i) => (
          <NavItem
            key={`${navitem}-${i}`}
            active={active === i}
            onClick={() => {
              navigate(navitem.path);
            }}
          >
            {navitem.name}
          </NavItem>
        ))}
      </nav>

      <div className="border-t border-(--border-default) p-2">
        <Button variant="danger" className="w-full">
          Logout
        </Button>
      </div>
    </aside>
  );
}

function NavItem({
  children,
  active,
  onClick = () => {},
}: {
  children: string;
  active?: boolean;
  onClick?: () => void;
}) {
  return (
    <div
      className={clsx(
        "rounded-xs px-3 py-2 text-sm cursor-pointer",
        active
          ? "bg-(--bg-elevated) text-(--text-primary)"
          : "text-(--text-secondary) hover:bg-(--bg-elevated)"
      )}
      onClick={onClick}
    >
      {children}
    </div>
  );
}
