type Role = "Owner" | "Member";

export interface ProjectSummaryApi {
  id: string;
  name: string;
  description?: string;
  skills?: string;
  role: Role;
  unassigned_tasks: number;
  ongoing_tasks: number;
  completed_tasks: number;
  created_at: string;
  updated_at?: string;
}

export interface ProjectSummary {
  id: string;
  name: string;
  description?: string;
  skills?: string;
  role: Role;
  unassignedTasks: number;
  ongoingTasks: number;
  completedTasks: number;
  createdAt: string;
  updatedAt?: string;
}

export interface ProjectsApiResponse {
  projects: ProjectSummaryApi[];
  page: number;
  limit: number;
  hasNext: boolean;
}

export const MapProject = (p: ProjectSummaryApi): ProjectSummary => ({
  id: p.id,
  name: p.name,
  description: p.description,
  skills: p.skills,
  role: p.role,
  unassignedTasks: p.unassigned_tasks,
  ongoingTasks: p.ongoing_tasks,
  completedTasks: p.completed_tasks,
  createdAt: p.created_at,
  updatedAt: p.updated_at,
});

export interface CreateProjectApi {
  id: string;
  name: string;
  description?: string;
  skills?: string;
  created_at: string;
  updated_at?: string;
}
