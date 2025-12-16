export function Logo() {
  return (
    <div className="flex flex-col items-center gap-2">
      <div className="h-10 w-10 rounded-sm bg-(--primary) flex items-center justify-center text-sm font-semibold text-white">
        PM
      </div>
      <h1 className="text-base font-semibold">Sign in to your workspace</h1>
      <p className="text-xs text-(--text-secondary) text-center leading-snug">
        Authenticate securely using your organization account
      </p>
    </div>
  );
}
