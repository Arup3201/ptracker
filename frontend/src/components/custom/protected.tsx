import { Navigate, Outlet } from "react-router";
import useAuth from "@/hooks/auth";
import { Loader2 } from "lucide-react";

const Protected = () => {
  const { isLoading, isAuthenticated } = useAuth();

  if (isLoading) {
    return <Loader2 className="mx-auto mt-4 animate-spin" size={48} />
  }

  return isAuthenticated ? <Outlet /> : <Navigate to="/login" />;
};

export default Protected;
