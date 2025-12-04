import { BrowserRouter as Router, Routes, Route } from "react-router";

import { AuthProvider } from "./contexts/auth-context";

import Login from "./pages/Login";
import Register from "./pages/Register";
import Protected from "./components/custom/protected";
import Dashboard from "./pages/Dashboard";
import Project from "./pages/Project";

function App() {
  return (
    <Router>
      <AuthProvider>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />

          <Route element={<Protected />}>
            <Route index element={<Dashboard />} />
            <Route path="/projects/:id" element={<Project />} />
          </Route>
        </Routes>
      </AuthProvider>
    </Router>
  );
}

export default App;
