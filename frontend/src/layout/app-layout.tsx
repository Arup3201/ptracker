import { Outlet, Link } from "react-router";
import { Sidebar } from "../components/sidebar";

export function AppLayout() {
  return (
    <div className="flex h-screen overflow-hidden bg-bg-root text-text-primary">
      <Sidebar />
      <main className="flex flex-1 flex-col overflow-hidden pb-16 md:pb-0">
        <Outlet />

        <footer className="border-t border-border bg-bg-root">
          <div className="max-w-5xl mx-auto px-4 py-4 flex flex-col sm:flex-row items-center justify-between gap-3">
            <p className="text-sm text-text-muted">
              © {new Date().getFullYear()} ProjectMate
            </p>

            <div className="flex items-center gap-4 text-sm">
              <Link
                to="/terms"
                className="text-text-muted hover:text-text-primary transition-colors"
              >
                Terms
              </Link>
              <Link
                to="/privacy"
                className="text-text-muted hover:text-text-primary transition-colors"
              >
                Privacy
              </Link>
              <Link
                to="/cookies"
                className="text-text-muted hover:text-text-primary transition-colors"
              >
                Cookies
              </Link>
              <Link
                to="/gdpr"
                className="text-text-muted hover:text-text-primary transition-colors"
              >
                GDPR
              </Link>
            </div>
          </div>
        </footer>
      </main>
    </div>
  );
}
