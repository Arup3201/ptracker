import { Toaster } from "react-hot-toast";
import { BrowserRouter as Router, Routes, Route } from "react-router";

import LoginPage from "./pages/login";
import RegisterPage from "./pages/register";
import { Dashboard } from "./pages/dashboard";
import { AppLayout } from "./layout/app-layout";
import { ProjectsPage } from "./pages/projects";
import ExploreProjectsPage from "./pages/explore";
import ProjectDetailsPage from "./pages/project-details";
import ProjectExplorePage from "./pages/explore-project";
import CheckEmailPage from "./pages/check-email";
import VerifyEmailPage from "./pages/verify-email";
import { AuthProvider } from "./context/auth";
import ProtectedRoute from "./components/protected-route";
import ForgotPasswordPage from "./pages/forgot-password";
import ResetPasswordPage from "./pages/reset-password";
import TermsPage from "./pages/terms";
import PrivacyPage from "./pages/privacy";
import CookiePage from "./pages/cookie";
import GDPRPage from "./pages/gdpr";
import { MessagesPage } from "./pages/messages";

function App() {
  return (
    <>
      <Toaster
        position="bottom-right"
        toastOptions={{
          duration: 7000,
        }}
      />
      <Router>
        <AuthProvider>
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route path="/forgot-password" element={<ForgotPasswordPage />} />
            <Route path="/reset-password" element={<ResetPasswordPage />} />
            <Route path="/register" element={<RegisterPage />} />
            <Route path="/check-email" element={<CheckEmailPage />} />
            <Route path="/verify-email" element={<VerifyEmailPage />} />

            <Route path="/terms" element={<TermsPage />} />
            <Route path="/privacy" element={<PrivacyPage />} />
            <Route path="/cookies" element={<CookiePage />} />
            <Route path="/gdpr" element={<GDPRPage />} />

            <Route element={<ProtectedRoute />}>
              <Route element={<AppLayout />}>
                <Route index element={<Dashboard />} />
                <Route path="/projects" element={<ProjectsPage />} />
                <Route path="/projects/:id" element={<ProjectDetailsPage />} />
                <Route path="/explore" element={<ExploreProjectsPage />} />
                <Route path="/explore/:id" element={<ProjectExplorePage />} />
                <Route path="/messages" element={<MessagesPage />} />
              </Route>
            </Route>
          </Routes>
        </AuthProvider>
      </Router>
    </>
  );
}

export default App;
