import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
  type Dispatch,
  type SetStateAction,
} from "react";
import { ApiFetch, getValidToken } from "../utils/api";
import { tokenStore } from "../utils/token";
import type { Avatar } from "../types/avatar";
import { userStore } from "../utils/user";

interface AuthContextValue {
  loading: boolean;
  setLoading: Dispatch<SetStateAction<boolean>>;
  refreshToken(): Promise<void>;
  user: Avatar | null;
  logout(): Promise<void>;
}

interface AuthProviderInterface {
  children: React.ReactNode;
}

const authContext = createContext<AuthContextValue>({
  loading: true,
  setLoading: () => {},
  refreshToken: () => Promise.resolve(),
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
    const token = tokenStore.get();
    if (token === null) {
      setLoading(true);
      try {
        await getValidToken();
      } catch (err) {
        console.error(err);
        throw new Error("Token refresh failed");
      } finally {
        setLoading(false);
      }
    }
  }

  const user = userStore.get();

  return (
    <authContext.Provider
      value={{
        loading,
        setLoading,
        refreshToken,
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
