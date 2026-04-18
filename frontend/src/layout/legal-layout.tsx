import { useNavigate } from "react-router";

interface LegalLayoutProps {
  title: string;
  lastUpdated?: string;
  children: React.ReactNode;
}

export default function LegalLayout({
  title,
  lastUpdated,
  children,
}: LegalLayoutProps) {
  const navigate = useNavigate();

  return (
    <div className="min-h-screen bg-bg-root text-text-primary">
      <div className="max-w-3xl mx-auto px-4 py-10">
        <button
          onClick={() => navigate(-1)}
          className="mb-6 inline-flex items-center text-text-muted hover:text-text-primary transition-colors"
        >
          <svg
            className="w-5 h-5 mr-2"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M15 19l-7-7 7-7"
            />
          </svg>
          Back
        </button>

        <div className="mb-8">
          <h1 className="text-3xl font-semibold tracking-tight">{title}</h1>

          {lastUpdated && (
            <p className="text-sm text-text-muted mt-2">
              Last updated: {lastUpdated}
            </p>
          )}
        </div>

        <div className="bg-bg-surface border border-border rounded-xl shadow-sm p-6 space-y-6">
          {children}
        </div>
      </div>
    </div>
  );
}
