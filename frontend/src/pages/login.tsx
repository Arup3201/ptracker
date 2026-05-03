import { useState } from "react";
import { Link, useNavigate } from "react-router";
import { Button } from "../components/button";
import { Card, CardContent } from "../components/card";
import { Input } from "../components/input";
import { Logo } from "../components/logo";
import { API_ROOT } from "../utils/api";
import { tokenStore } from "../utils/token";
import { userStore } from "../utils/user";

export default function LoginPage() {
  const navigate = useNavigate();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      const res = await fetch(API_ROOT + "/auth/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include", // for refresh token cookie
        body: JSON.stringify({
          email: email,
          password: password,
        }),
      });
      if (!res.ok) throw new Error("Login failed");

      const respondData = await res.json();
      if (respondData.data) {
        tokenStore.set(respondData.data.access_token);
        userStore.set(respondData.data.user);
      } else {
        throw new Error("Empty response data from login.");
      }

      navigate("/");
    } catch (err: any) {
      setError(err.message ?? "Something went wrong.");
    } finally {
      setLoading(false);
    }
  }

  async function handleGoogleLogin() {
    try {
      const response = await fetch(API_ROOT + "/auth/google/redirect");
      const json = await response.json();
      if (json.data) {
        window.open(json.data, "_parent");
      } else {
        throw new Error("No redirect URL in the response!");
      }
    } catch (err) {
      console.error(err);
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

              <Button
                type="button"
                variant="secondary"
                className="w-full flex items-center justify-center gap-2"
                onClick={handleGoogleLogin}
              >
                <svg
                  className="w-4 h-4"
                  viewBox="0 0 24 24"
                  fill="currentColor"
                >
                  <path
                    d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                    fill="#4285F4"
                  />
                  <path
                    d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                    fill="#34A853"
                  />
                  <path
                    d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                    fill="#FBBC05"
                  />
                  <path
                    d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                    fill="#EA4335"
                  />
                </svg>
                Sign in with Google
              </Button>

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
