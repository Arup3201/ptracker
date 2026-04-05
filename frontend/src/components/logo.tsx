export function Logo() {
  return (
    <div className="flex items-center gap-2.5">
      <div className="h-7 w-7 rounded-lg bg-root flex items-center justify-center shadow-sm">
        <span className="text-xs font-semibold text-primary-foreground tracking-tight">
          PM
        </span>
      </div>
      <span className="text-base font-semibold text-text-primary tracking-snug">
        ProjectMate
      </span>
    </div>
  );
}
