import { Card, CardContent } from "../components/card.tsx";
import { Button } from "../components/button.tsx";
import { Logo } from "../components/logo.tsx";

export default function LoginPage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-(--bg-root) text-(--text-primary)">
      <div className="w-90">
        <Card>
          <CardContent>
            <div className="flex flex-col gap-5">
              <Logo />

              <Button
                onClick={() => {
                  // Redirect to Keycloak
                }}
              >
                Login with Keycloak
              </Button>

              <p className="text-[11px] text-(--text-muted) text-center leading-relaxed">
                By continuing, you agree to your organization’s authentication
                policy.
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
