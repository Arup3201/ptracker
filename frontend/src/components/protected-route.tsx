import { Navigate, Outlet } from "react-router";
import { useAuth } from "../context/auth";
import { useEffect } from "react";

export default function ProtectedRoute() {
  const { user, loading, refreshToken } = useAuth();

  useEffect(() => {
    refreshToken().catch((err) => console.error(err));
  }, []);

  if (loading) return null; // or loader

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  return <Outlet />;
}
