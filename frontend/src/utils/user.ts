import { MapAvatar, type Avatar, type AvatarApi } from "../types/avatar";

var user: Avatar | null = null;

const userStore = {
  set: (u: AvatarApi) => {
    if (u) {
      user = MapAvatar(u);
      return;
    }

    throw new Error("user is null");
  },
  get: () => user,
  clear: () => (user = null),
};

export { userStore };
