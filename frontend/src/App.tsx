import { BrowserRouter as Router, Routes, Route } from "react-router";

import LoginPage from "./pages/auth";
import { Dashboard } from "./pages/dashboard";
import { AppLayout } from "./layout/app-layout";

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route element={<AppLayout />}>
          <Route index element={<Dashboard />} />
        </Route>
      </Routes>
    </Router>
  );
}

export default App;
