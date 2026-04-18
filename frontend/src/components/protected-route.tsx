import { Navigate, Outlet } from "react-router";
import { useAuth } from "../context/auth";

export default function ProtectedRoute() {
  const { user, loading } = useAuth();

  if (loading) return null; // or loader

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  return <Outlet />;
}
