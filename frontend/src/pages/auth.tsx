import { Card, CardContent } from "../components/card.tsx";
import { Button } from "../components/button.tsx";
import { Logo } from "../components/logo.tsx";

const API_ROOT = "http://localhost:8081";

export default function LoginPage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-(--bg-root) text-(--text-primary)">
      <div className="w-90">
        <Card>
          <CardContent>
            <div className="flex flex-col gap-5">
              <div className="flex flex-col items-center gap-2">
                <Logo />
                <h1 className="text-base font-semibold">Sign in to Collap</h1>
                <p className="text-xs text-(--text-secondary) text-center leading-snug">
                  Authenticate securely using your account
                </p>
              </div>

              <Button
                onClick={() => {
                  window.location.href = API_ROOT + "/api/auth/login";
                }}
              >
                Login with Keycloak
              </Button>

              <p className="text-[11px] text-(--text-muted) text-center leading-relaxed">
                By continuing, you agree to Keycloak's authentication policy.
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
