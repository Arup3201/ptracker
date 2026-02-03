import { useEffect, useState } from "react";
import { MapUser, type User, type UserApi } from "../types/user";
import { ApiRequest } from "../api/request";

export function useCurrentUser() {
  const [user, setUser] = useState<User>();

  async function getUser() {
    try {
      const data = await ApiRequest<UserApi>("/auth/me", "GET", null);
      if (data) {
        setUser(MapUser(data));
      }
    } catch (err) {
      console.error(err);
    }
  }

  useEffect(() => {
    getUser();
  }, []);

  return user;
}
