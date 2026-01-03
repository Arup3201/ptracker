export interface Project {
  id: string;
  name: string;
  description: string;
  skills: string;
  role: "Owner" | "Member";
  unassignedTasks: number;
  ongoingTasks: number;
  completedTasks: number;
  createdAt: string;
  updatedAt: string;
}
