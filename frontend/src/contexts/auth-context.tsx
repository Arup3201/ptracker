import { createContext, useEffect, useState } from "react";

import { HttpGet } from "@/utils/http";
import type { User } from "@/types/user";

interface AuthData {
  isLoading: boolean;
  isAuthenticated: boolean;
  setIsAuthenticated: (_: boolean) => void;
  user: User | undefined;
  setUser: (_: User) => void;
}

const AuthContext = createContext<AuthData>({
  isLoading: false,
  isAuthenticated: false,
  setIsAuthenticated: (_: boolean) => {},
  user: {} as User,
  setUser: (_: User) => {},
});

const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const [user, setUser] = useState<User | undefined>();
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  const getUser = async () => {
    try {
      const data = await HttpGet("/auth/me");
      setUser({
        id: data.id,
        name: data.name,
        email: data.email,
      });
      setIsAuthenticated(true);
    } catch (err) {
      if (err instanceof Error) {
        console.error(`getUser failed with error: ${err.message}`);
        setUser(undefined);
        setIsAuthenticated(false);
      }
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    getUser();
  }, []);

  return (
    <AuthContext.Provider
      value={{ user, setUser, isAuthenticated, setIsAuthenticated, isLoading }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export { AuthContext, AuthProvider };
