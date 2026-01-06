export interface Member {
  id: string;
  name: string;
  role: "Owner" | "Member";
  joinedAt: string;
}
