import { useState } from "react";
import { useNavigate, useSearchParams } from "react-router";
import { Card, CardContent } from "../components/card";
import { Input } from "../components/input";
import { Logo } from "../components/logo";
import { PasswordStrength } from "./register";
import { Button } from "../components/button";
import { API_ROOT } from "../utils/api";

export default function ResetPasswordPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);

    if (password !== confirm) {
      setError("Passwords do not match.");
      return;
    }
    if (password.length < 8) {
      setError("Password must be at least 8 characters.");
      return;
    }

    const token = searchParams.get("token");

    if (!token) {
      setError("Reset password token is missing.");
      return;
    }

    setLoading(true);
    try {
      const res = await fetch(API_ROOT + "/auth/password-reset", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          password: password,
          token: token,
        }),
      });

      if (res.ok) {
        navigate("/login");
      } else {
        const data = await res.json().catch(() => null);
        throw new Error(
          data?.message || "Failed to resend verification email.",
        );
      }
    } catch (err) {
      console.error(err);
      setError((err as Error).message);
    } finally {
      setLoading(false);
    }
  }

  const passwordMismatch = confirm.length > 0 && password !== confirm;

  return (
    <div className="min-h-screen flex items-center justify-center bg-(--bg-root) text-(--text-primary) px-4">
      <div className="w-full max-w-90">
        <Card>
          <CardContent className="p-6">
            <form onSubmit={handleSubmit} className="flex flex-col gap-5">
              <div className="flex flex-col items-center gap-2">
                <Logo />
                <h1 className="text-lg font-semibold tracking-tight">
                  Reset Password
                </h1>
                <p className="text-xs text-(--text-secondary) text-center leading-snug">
                  You will set the new password here.
                </p>
              </div>

              {error && (
                <div className="rounded-md bg-danger-muted border border-danger/20 px-3 py-2 text-xs text-danger leading-snug">
                  {error}
                </div>
              )}

              <div className="flex flex-col gap-1.5">
                <label className="text-xs font-medium text-text-secondary">
                  Password
                </label>
                <Input
                  type="password"
                  placeholder="••••••••"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                  autoComplete="new-password"
                />
                <PasswordStrength password={password} />
              </div>

              <div className="flex flex-col gap-1.5">
                <label className="text-xs font-medium text-text-secondary">
                  Confirm password
                </label>
                <Input
                  type="password"
                  placeholder="••••••••"
                  value={confirm}
                  onChange={(e) => setConfirm(e.target.value)}
                  required
                  autoComplete="new-password"
                  className={
                    passwordMismatch
                      ? "border-danger focus:shadow-focus-danger"
                      : undefined
                  }
                />
                {passwordMismatch && (
                  <p className="text-[11px] text-danger">
                    Passwords do not match.
                  </p>
                )}
              </div>

              <Button type="submit" disabled={loading} className="w-full">
                {loading ? "Resetting…" : "Reset"}
              </Button>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
