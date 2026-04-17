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
import { AuthProvider } from "./context/auth";

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
            <Route path="/register" element={<RegisterPage />} />
            <Route element={<AppLayout />}>
              <Route index element={<Dashboard />} />
              <Route path="/projects" element={<ProjectsPage />} />
              <Route path="/projects/:id" element={<ProjectDetailsPage />} />
              <Route path="/explore" element={<ExploreProjectsPage />} />
              <Route path="/explore/:id" element={<ProjectExplorePage />} />
            </Route>
          </Routes>
        </AuthProvider>
      </Router>
    </>
  );
}

export default App;
