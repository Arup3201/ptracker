import { useState } from "react";
import { API_ROOT } from "../utils/api";
import { useSearchParams } from "react-router";

export default function CheckEmailPage() {
  const [searchParams] = useSearchParams();
  const [loading, setLoading] = useState(false);
  const [sent, setSent] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleResend = async () => {
    setLoading(true);
    setError(null);
    setSent(false);

    try {
      const email = searchParams.get("email");
      if (!email) throw new Error("Missing email in URL.");

      const res = await fetch(API_ROOT + "/auth/resend-verification", {
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
        throw new Error(
          data?.message || "Failed to resend verification email.",
        );
      }
    } catch (err: any) {
      setError(err.message || "Something went wrong.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-bg-root px-4">
      <div className="w-full max-w-md bg-bg-surface border border-border rounded-xl shadow-sm p-6 animate-fade-in">
        <h1 className="text-2xl font-semibold text-text-primary tracking-snug mb-2">
          Check your email
        </h1>
        <p className="text-text-secondary text-sm mb-6">
          We have sent you a verification link. Please check your inbox and
          click the link to verify your account.
        </p>

        <div className="text-sm text-text-secondary">
          Didn't receive the email?{" "}
          <button
            onClick={handleResend}
            disabled={loading}
            className="cursor-pointer text-primary hover:text-primary-hover font-medium transition"
          >
            {loading ? "Sending..." : "Resend email"}
          </button>
        </div>

        {sent && (
          <p className="mt-3 text-sm text-success">
            Verification email resent successfully.
          </p>
        )}

        {error && <p className="mt-3 text-sm text-danger">{error}</p>}
      </div>
    </div>
  );
}
