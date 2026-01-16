import { type Role } from "./project";

export interface Member {
  id: string;
  name: string;
  role: Role;
  joinedAt: string;
}
