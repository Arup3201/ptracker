import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
} from "react";
import { ApiFetch, getValidToken } from "../utils/api";
import { tokenStore } from "../utils/token";
import type { Avatar } from "../types/avatar";
import { userStore } from "../utils/user";

interface AuthContextValue {
  loading: boolean;
  user: Avatar | null;
  logout(): Promise<void>;
}

interface AuthProviderInterface {
  children: React.ReactNode;
}

const authContext = createContext<AuthContextValue>({
  loading: true,
  user: null,
  logout: () => Promise.resolve(),
});

const AuthProvider: React.FC<AuthProviderInterface> = ({ children }) => {
  const [loading, setLoading] = useState(true);

  const logout = useCallback(async () => {
    const response = await ApiFetch("/auth/logout", {
      method: "POST",
      credentials: "include",
    });
    if (response.status === 200) {
      console.log("logging out...");
    } else {
      console.error("Something went wrong during logout.");
    }
  }, []);

  useEffect(() => {
    window.addEventListener("auth:logout", logout);
    return () => window.removeEventListener("auth:logout", logout);
  }, [logout]);

  async function refreshToken() {
    setLoading(true);
    try {
      await getValidToken();
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    const token = tokenStore.get();
    if (token === null) {
      refreshToken();
    }
  }, []);

  const user = userStore.get();

  return (
    <authContext.Provider
      value={{
        loading,
        user,
        logout,
      }}
    >
      {children}
    </authContext.Provider>
  );
};

const useAuth = () => useContext(authContext);

export { useAuth, AuthProvider };
