import clsx from "clsx";
import { useNavigate } from "react-router";
import { Logo } from "./logo";
import { Button } from "./button";
import { useCurrentUser } from "../hooks/current_user";
import { ApiRequest } from "../api/request";

import type { LucideIcon } from "lucide-react";
import {
  LogOut,
  ChevronsUpDown,
  MenuIcon,
  SquareCheckBig,
  Telescope,
} from "lucide-react";

const NavItems = [
  {
    name: "Dashboard",
    path: "/",
    icon: MenuIcon,
  },
  {
    name: "Projects",
    path: "/projects",
    icon: SquareCheckBig,
  },
  {
    name: "Explore",
    path: "/explore",
    icon: Telescope,
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

  const currentUser = useCurrentUser();

  const initials = currentUser
    ? currentUser?.displayName[0].toUpperCase()
    : "U";

  async function handleLogout() {
    try {
      await ApiRequest("/auth/logout", "POST", null);
      navigate("/login");
    } catch (err) {
      console.error(err);
    }
  }

  return (
    <aside className="flex w-56 flex-col bg-bg-surface border-r border-border">
      {/* Logo */}
      <div className="h-14 flex items-center px-4 border-b border-border shrink-0">
        <Logo />
      </div>

      {/* Nav */}
      <nav className="flex-1 overflow-y-auto px-2 py-3 space-y-0.5">
        {NavItems.map((navitem, i) => (
          <NavItem
            key={`${navitem.name}-${i}`}
            active={active === i}
            icon={navitem.icon}
            onClick={() => navigate(navitem.path)}
          >
            {navitem.name}
          </NavItem>
        ))}
      </nav>

      {/* User footer */}
      <div className="shrink-0 border-t border-zinc-800 p-2 space-y-1">
        <div className="flex items-center gap-2.5 rounded-lg px-2 py-1.5 hover:bg-zinc-900 transition-colors duration-150 cursor-default">
          <div className="h-7 w-7 rounded-full bg-zinc-800 flex items-center justify-center text-xs font-semibold text-emerald-400 shrink-0 ring-1 ring-zinc-700">
            {initials}
          </div>
          <div className="flex-1 min-w-0">
            <p className="text-sm font-semibold text-text-primary truncate leading-tight">
              {currentUser?.displayName}
            </p>
            <p className="text-xs text-text-muted truncate leading-tight">
              {currentUser?.email}
            </p>
          </div>
          <ChevronsUpDown size={13} className="text-text-muted shrink-0" />
        </div>

        <Button
          variant="danger"
          className="w-full gap-1.5"
          onClick={handleLogout}
        >
          <LogOut size={13} />
          Sign out
        </Button>
      </div>
    </aside>
  );
}

function NavItem({
  children,
  active,
  icon: Icon,
  onClick = () => {},
}: {
  children: string;
  active?: boolean;
  icon?: LucideIcon;
  onClick?: () => void;
}) {
  return (
    <div
      onClick={onClick}
      className={clsx(
        "group flex items-center gap-2.5 rounded-md px-2.5 py-1.5 text-sm cursor-pointer transition-all duration-fast select-none",
        active
          ? "bg-zinc-800 text-zinc-100 font-medium"
          : "text-zinc-500 hover:bg-zinc-900 hover:text-zinc-300",
      )}
    >
      {Icon && (
        <Icon
          size={15}
          className={clsx(
            "shrink-0 transition-colors duration-150",
            active
              ? "text-emerald-400"
              : "text-zinc-600 group-hover:text-zinc-400",
          )}
        />
      )}
      <span className="truncate">{children}</span>
      {active && (
        <span className="ml-auto h-1.5 w-1.5 rounded-full bg-emerald-400 shrink-0" />
      )}
    </div>
  );
}
