import { BrowserRouter as Router, Routes, Route } from "react-router";
import LoginPage from "./pages/auth";

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
      </Routes>
    </Router>
  );
}

export default App;
