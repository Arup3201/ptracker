import { BrowserRouter as Router, Routes, Route } from "react-router";

import LoginPage from "./pages/auth";
import { Dashboard } from "./pages/dashboard";
import { AppLayout } from "./layout/app-layout";
import { ProjectsPage } from "./pages/projects";

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route element={<AppLayout />}>
          <Route index element={<Dashboard />} />
          <Route path="/projects" element={<ProjectsPage />} />
        </Route>
      </Routes>
    </Router>
  );
}

export default App;
