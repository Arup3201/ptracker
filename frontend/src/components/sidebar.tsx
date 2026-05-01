import clsx from "clsx";
import { useNavigate } from "react-router";
import { Logo } from "./logo";
import { Button } from "./button";

import type { LucideIcon } from "lucide-react";
import {
  LogOut,
  ChevronsUpDown,
  MenuIcon,
  SquareCheckBig,
  Telescope,
  X,
  BellRing,
} from "lucide-react";
import { useState } from "react";
import { useAuth } from "../context/auth";

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
  {
    name: "Messages",
    path: "/messages",
    icon: BellRing,
  },
];

export function Sidebar() {
  const navigate = useNavigate();
  const { user, logout } = useAuth();

  const [showUserMenu, setShowUserMenu] = useState(false);

  let active: number;
  NavItems.forEach((navItem, i) => {
    if (window.location.pathname === navItem.path) {
      active = i;
    }
  });

  const initials = user
    ? user.displayName
      ? user.displayName[0].toUpperCase()
      : user.username[0].toUpperCase()
    : "U";

  async function handleLogout() {
    try {
      await logout();
      navigate("/login");
    } catch (err) {
      console.error(err);
    }
  }

  return (
    <>
      {/* ── Desktop sidebar ── */}
      <aside className="hidden md:flex w-56 flex-col bg-bg-surface border-r border-border shrink-0">
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
        <div className="shrink-0 border-t border-border p-2 space-y-1">
          <div className="flex items-center gap-2.5 rounded-md px-2 py-1.5 hover:bg-bg-elevated transition duration-fast cursor-default">
            <div className="h-7 w-7 rounded-full bg-bg-elevated flex items-center justify-center text-xs font-semibold text-primary shrink-0 ring-1 ring-border">
              {initials}
            </div>
            <div className="flex-1 min-w-0">
              <p className="text-sm font-semibold text-text-primary truncate leading-tight">
                {user?.displayName}
              </p>
              <p className="text-xs text-text-muted truncate leading-tight">
                {user?.email}
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

      {/* ── Mobile bottom nav ── */}
      <nav className="md:hidden fixed bottom-0 left-0 right-0 z-40 bg-zinc-900 border-t border-zinc-800 flex items-center justify-between px-3 h-16">
        {NavItems.map((navitem, i) => (
          <MobileNavItem
            key={`${navitem.name}-${i}`}
            active={active === i}
            icon={navitem.icon}
            onClick={() => navigate(navitem.path)}
          >
            {navitem.name}
          </MobileNavItem>
        ))}

        {/* Avatar as last item — same width as nav items */}
        <button
          onClick={() => setShowUserMenu(true)}
          className="relative flex flex-col items-center justify-center gap-1 py-1.5 rounded-md w-full transition duration-fast cursor-pointer hover:bg-zinc-800 group"
        >
          <div className="h-4.5 w-4.5 rounded-full bg-zinc-800 flex items-center justify-center text-[10px] font-semibold text-emerald-400 ring-1 ring-zinc-700 group-hover:ring-emerald-500 transition duration-fast">
            {initials}
          </div>
          <span className="text-[10px] font-medium leading-none text-zinc-500 group-hover:text-zinc-300 transition duration-fast">
            Account
          </span>
        </button>
      </nav>

      {/* ── Mobile user menu overlay ── */}
      {showUserMenu && (
        <div className="md:hidden fixed inset-0 z-50 flex flex-col justify-end">
          {/* Backdrop */}
          <div
            className="absolute inset-0 bg-black/60 backdrop-blur-sm"
            onClick={() => setShowUserMenu(false)}
          />

          {/* Sheet */}
          <div className="relative bg-bg-elevated border-t border-border rounded-t-xl p-4 space-y-4 animate-slide-up">
            {/* Close */}
            <button
              onClick={() => setShowUserMenu(false)}
              className="absolute top-3.5 right-4 h-7 w-7 flex items-center justify-center rounded-md text-text-muted hover:text-text-primary hover:bg-bg-overlay transition duration-fast"
            >
              <X size={15} />
            </button>

            {/* User info */}
            <div className="flex items-center gap-3 pr-8">
              <div className="h-10 w-10 rounded-full bg-bg-overlay flex items-center justify-center text-sm font-semibold text-primary shrink-0 ring-1 ring-border">
                {initials}
              </div>
              <div className="min-w-0">
                <p className="text-sm font-semibold text-text-primary truncate">
                  {user?.displayName}
                </p>
                <p className="text-xs text-text-muted truncate">
                  {user?.email}
                </p>
              </div>
            </div>

            {/* Sign out */}
            <Button
              variant="danger"
              className="w-full gap-1.5"
              onClick={handleLogout}
            >
              <LogOut size={13} />
              Sign out
            </Button>
          </div>
        </div>
      )}
    </>
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

function MobileNavItem({
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
    <button
      onClick={onClick}
      className={clsx(
        "relative flex flex-col items-center justify-center gap-1 py-1.5 rounded-md transition duration-fast select-none cursor-pointer w-full",
        active ? "text-emerald-400" : "text-zinc-500 hover:text-zinc-300",
      )}
    >
      {active && <span className="absolute inset-0 rounded-md bg-zinc-800" />}

      {/* Icon wrapper — dot anchored to icon */}
      <span className="relative">
        {Icon && (
          <Icon
            size={18}
            className={clsx(
              "relative shrink-0 transition duration-fast",
              active ? "text-emerald-400" : "text-zinc-500",
            )}
          />
        )}
        {active && (
          <span className="absolute -top-0.5 -right-0.5 h-1.5 w-1.5 rounded-full bg-emerald-400" />
        )}
      </span>

      <span className="relative text-[10px] font-medium leading-none truncate">
        {children}
      </span>
    </button>
  );
}
