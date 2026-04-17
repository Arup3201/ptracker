import { useEffect, useState } from "react";
import { useSearchParams, Link } from "react-router";
import { API_ROOT } from "../utils/api";

type Status = "loading" | "success" | "error";

export default function VerifyEmailPage() {
  const [searchParams] = useSearchParams();
  const [status, setStatus] = useState<Status>("loading");
  const [message, setMessage] = useState("");

  useEffect(() => {
    const token = searchParams.get("token");

    if (!token) {
      setStatus("error");
      setMessage("Invalid verification link.");
      return;
    }
    const verify = async () => {
      try {
        const res = await fetch(API_ROOT + "/auth/verify-email", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            token: token,
          }),
        });

        if (res.ok) {
          setStatus("success");
        } else {
          throw new Error("Token expired or invalid.");
        }
      } catch (err: any) {
        setStatus("error");
        setMessage(err.message || "Verification failed.");
      }
    };

    verify();
  }, [searchParams]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-bg-root px-4">
      <div className="w-full max-w-md bg-bg-surface border border-border rounded-xl shadow-sm p-6 animate-fade-in">
        {status === "loading" && (
          <>
            <h1 className="text-2xl font-semibold text-text-primary mb-2">
              Verifying...
            </h1>
            <p className="text-sm text-text-secondary">
              Please wait while we verify your email.
            </p>
          </>
        )}

        {status === "success" && (
          <>
            <h1 className="text-2xl font-semibold text-text-primary mb-2">
              Email verified
            </h1>
            <p className="text-sm text-text-secondary mb-6">
              Your email has been successfully verified.
            </p>

            <Link
              to="/login"
              className="inline-flex items-center justify-center w-full h-10 rounded-md bg-primary text-primary-foreground text-sm font-medium hover:bg-primary-hover transition"
            >
              Go to login
            </Link>
          </>
        )}

        {status === "error" && (
          <>
            <h1 className="text-2xl font-semibold text-text-primary mb-2">
              Verification failed
            </h1>
            <p className="text-sm text-danger mb-6">
              {message || "Invalid or expired token."}
            </p>

            <Link
              to="/resend-verification"
              className="text-primary hover:text-primary-hover text-sm font-medium"
            >
              Resend verification email
            </Link>
          </>
        )}
      </div>
    </div>
  );
}
