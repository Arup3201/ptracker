import { useState } from "react";
import { Link, useNavigate } from "react-router";
import { Button } from "../components/button";
import { Card, CardContent } from "../components/card";
import { Input } from "../components/input";
import { Logo } from "../components/logo";
import { API_ROOT } from "../utils/api";

function PasswordStrength({ password }: { password: string }) {
  const checks = [
    password.length >= 8,
    /[A-Z]/.test(password),
    /[0-9]/.test(password),
    /[^A-Za-z0-9]/.test(password),
  ];
  const score = checks.filter(Boolean).length;
  const labels = ["", "Weak", "Fair", "Good", "Strong"];
  const segmentColors = [
    "", // 0 - unused
    "bg-danger",
    "bg-warning",
    "bg-primary",
    "bg-success",
  ];

  if (!password) return null;

  return (
    <div className="flex flex-col gap-1.5 mt-1">
      <div className="flex gap-1">
        {[1, 2, 3, 4].map((i) => (
          <div
            key={i}
            className={`h-0.5 flex-1 rounded-full transition-colors duration-300 ${
              i <= score ? segmentColors[score] : "bg-border-muted"
            }`}
          />
        ))}
      </div>
      <span className="text-[11px] text-text-muted">{labels[score]}</span>
    </div>
  );
}

export default function RegisterPage() {
  const navigate = useNavigate();

  const [username, setUsername] = useState("");
  const [displayName, setDisplayName] = useState("");
  const [email, setEmail] = useState("");
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

    setLoading(true);
    try {
      const res = await fetch(API_ROOT + "/auth/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          email: email,
          username: username,
          display_name: displayName,
          password: password,
        }),
      });
      if (res.status === 201) {
        console.log("User created!");
        navigate(`/verify`);
      } else {
        throw new Error("User registration failed.");
      }
    } catch (err: any) {
      setError(err.message ?? "Something went wrong.");
    } finally {
      setLoading(false);
    }
  }

  const passwordMismatch = confirm.length > 0 && password !== confirm;

  return (
    <div className="min-h-screen flex items-center justify-center bg-bg-root text-text-primary px-4">
      <div className="w-full max-w-90">
        <Card>
          <CardContent className="p-6">
            <form onSubmit={handleSubmit} className="flex flex-col gap-5">
              {/* Header */}
              <div className="flex flex-col items-center gap-2">
                <Logo />
                <h1 className="text-lg font-semibold tracking-tight">
                  Create an account
                </h1>
                <p className="text-xs text-text-secondary text-center leading-snug">
                  Get started with us in seconds
                </p>
              </div>

              {/* Error banner */}
              {error && (
                <div className="rounded-md bg-danger-muted border border-danger/20 px-3 py-2 text-xs text-danger leading-snug">
                  {error}
                </div>
              )}

              {/* Fields */}
              <div className="flex flex-col gap-3">
                <div className="flex flex-col gap-1.5">
                  <label className="text-xs font-medium text-text-secondary">
                    Username
                  </label>
                  <Input
                    type="text"
                    placeholder="doejane"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    required
                    autoComplete="name"
                  />
                </div>

                <div className="flex flex-col gap-1.5">
                  <label className="text-xs font-medium text-text-secondary">
                    Display Name
                  </label>
                  <Input
                    type="text"
                    placeholder="Jane Doe"
                    value={displayName}
                    onChange={(e) => setDisplayName(e.target.value)}
                    required
                    autoComplete="display_name"
                  />
                </div>

                <div className="flex flex-col gap-1.5">
                  <label className="text-xs font-medium text-text-secondary">
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
              </div>

              {/* Submit */}
              <Button
                type="submit"
                size="lg"
                disabled={loading || passwordMismatch}
                className="w-full"
              >
                {loading ? "Creating account…" : "Create account"}
              </Button>

              {/* Login link */}
              <p className="text-[11px] text-text-muted text-center">
                Already have an account?{" "}
                <Link
                  to="/login"
                  className="text-text-secondary hover:text-text-primary font-medium transition-colors"
                >
                  Sign in
                </Link>
              </p>
            </form>
          </CardContent>
        </Card>

        <p className="mt-4 text-[11px] text-text-muted text-center leading-relaxed px-4">
          By creating an account, you agree to our{" "}
          <a
            href="/terms"
            className="underline underline-offset-2 hover:text-text-secondary transition-colors"
          >
            Terms of Service
          </a>{" "}
          and{" "}
          <a
            href="/privacy"
            className="underline underline-offset-2 hover:text-text-secondary transition-colors"
          >
            Privacy Policy
          </a>
          .
        </p>
      </div>
    </div>
  );
}
