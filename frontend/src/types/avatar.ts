export interface Avatar {
  userId: string;
  username: string;
  displayName?: string;
  email: string;
  avatarUrl?: string;
}

export interface AvatarApi {
  user_id: string;
  username: string;
  display_name?: string;
  email: string;
  avatar_url?: string;
}

export const MapAvatar = (a: AvatarApi): Avatar => {
  return {
    userId: a.user_id,
    username: a.username,
    displayName: a.display_name,
    email: a.email,
    avatarUrl: a.avatar_url,
  };
};
