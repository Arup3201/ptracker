import { useState } from "react";
import { Button } from "../components/button";
import { Card, CardContent } from "../components/card";
import { Input } from "../components/input";
import { Logo } from "../components/logo";
import { Link } from "react-router";
import { useAuth } from "../context/auth";

export default function LoginPage() {
  const { login } = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      await login({ email: email, password: password });
    } catch (err: any) {
      setError(err.message ?? "Something went wrong.");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-(--bg-root) text-(--text-primary) px-4">
      <div className="w-full max-w-90">
        <Card>
          <CardContent className="p-6">
            <form onSubmit={handleSubmit} className="flex flex-col gap-5">
              {/* Header */}
              <div className="flex flex-col items-center gap-2">
                <Logo />
                <h1 className="text-lg font-semibold tracking-tight">
                  Welcome back
                </h1>
                <p className="text-xs text-(--text-secondary) text-center leading-snug">
                  Sign in to your account
                </p>
              </div>

              {/* Error banner */}
              {error && (
                <div className="rounded-md bg-(--danger-muted) border border-(--danger)/20 px-3 py-2 text-xs text-(--danger) leading-snug">
                  {error}
                </div>
              )}

              {/* Fields */}
              <div className="flex flex-col gap-3">
                <div className="flex flex-col gap-1.5">
                  <label className="text-xs font-medium text-(--text-secondary)">
                    Email
                  </label>
                  <Input
                    type="email"
                    placeholder="you@example.com"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                    autoComplete="email"
                  />
                </div>

                <div className="flex flex-col gap-1.5">
                  <div className="flex items-center justify-between">
                    <label className="text-xs font-medium text-(--text-secondary)">
                      Password
                    </label>
                    <Link
                      to="/forgot-password"
                      className="text-[11px] text-(--text-muted) hover:text-(--text-secondary) transition-colors"
                    >
                      Forgot password?
                    </Link>
                  </div>
                  <Input
                    type="password"
                    placeholder="••••••••"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                    autoComplete="current-password"
                  />
                </div>
              </div>

              {/* Submit */}
              <Button type="submit" disabled={loading} className="w-full">
                {loading ? "Signing in…" : "Sign in"}
              </Button>

              {/* Divider */}
              <div className="relative flex items-center gap-3">
                <div className="flex-1 h-px bg-(--border-muted)" />
                <span className="text-[11px] text-(--text-muted)">or</span>
                <div className="flex-1 h-px bg-(--border-muted)" />
              </div>

              {/* Register link */}
              <p className="text-[11px] text-(--text-muted) text-center">
                Don't have an account?{" "}
                <Link
                  to="/register"
                  className="text-(--text-secondary) hover:text-(--text-primary) font-medium transition-colors"
                >
                  Create one
                </Link>
              </p>
            </form>
          </CardContent>
        </Card>

        <p className="mt-4 text-[11px] text-(--text-muted) text-center leading-relaxed px-4">
          By signing in, you agree to our{" "}
          <a
            href="/terms"
            className="underline underline-offset-2 hover:text-(--text-secondary) transition-colors"
          >
            Terms of Service
          </a>{" "}
          and{" "}
          <a
            href="/privacy"
            className="underline underline-offset-2 hover:text-(--text-secondary) transition-colors"
          >
            Privacy Policy
          </a>
          .
        </p>
      </div>
    </div>
  );
}
