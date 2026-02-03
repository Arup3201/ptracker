export interface User {
  id: string;
  username: string;
  displayName: string;
  email: string;
  createdAt: string;
  updatedAt: string;
}

export interface UserApi {
  id: string;
  username: string;
  display_name: string;
  email: string;
  created_at: string;
  updated_at: string;
}

export const MapUser = (u: UserApi): User => ({
  id: u.id,
  username: u.username,
  displayName: u.display_name,
  email: u.display_name,
  createdAt: u.created_at,
  updatedAt: u.updated_at,
});
