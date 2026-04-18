import { useState } from "react";
import { Button } from "../components/button";
import { Card, CardContent } from "../components/card";
import { Input } from "../components/input";
import { Logo } from "../components/logo";
import { API_ROOT } from "../utils/api";

export default function ForgotPasswordPage() {
  const [loading, setLoading] = useState(false);
  const [email, setEmail] = useState("");
  const [sent, setSent] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setSent(false);

    try {
      if (!email) throw new Error("Missing email in URL.");

      const res = await fetch(API_ROOT + "/auth/password-reset-email", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          email: email,
        }),
      });

      if (res.ok) {
        setSent(true);
      } else {
        const data = await res.json().catch(() => null);
        throw new Error(data?.message || "Failed to send reset email.");
      }
    } catch (err: any) {
      setError(err.message || "Something went wrong.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-(--bg-root) text-(--text-primary) px-4">
      <div className="w-full max-w-90">
        <Card>
          <CardContent className="p-6">
            <form onSubmit={handleSubmit} className="flex flex-col gap-5">
              <div className="flex flex-col items-center gap-2">
                <Logo />
                <h1 className="text-lg font-semibold tracking-tight">
                  Forgot Password?
                </h1>
                <p className="text-xs text-(--text-secondary) text-center leading-snug">
                  You will receive a password reset link on this email address.
                </p>
              </div>

              {error && (
                <div className="rounded-md bg-(--danger-muted) border border-(--danger)/20 px-3 py-2 text-xs text-(--danger) leading-snug">
                  {error}
                </div>
              )}
              {sent && (
                <p className="mt-3 text-sm text-success">
                  Password reset email sent successfully.
                </p>
              )}

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
                  />
                </div>
              </div>

              <Button type="submit" disabled={loading} className="w-full">
                {loading ? "Sending…" : "Send Link"}
              </Button>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
