import { BrowserRouter as Router, Routes, Route } from "react-router";

import LoginPage from "./pages/auth";
import { Dashboard } from "./pages/dashboard";
import { AppLayout } from "./layout/app-layout";
import { ProjectsPage } from "./pages/projects";
import ExploreProjectsPage from "./pages/explore";

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route element={<AppLayout />}>
          <Route index element={<Dashboard />} />
          <Route path="/projects" element={<ProjectsPage />} />
          <Route path="/explore" element={<ExploreProjectsPage />} />
        </Route>
      </Routes>
    </Router>
  );
}

export default App;
