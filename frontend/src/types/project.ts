export const ROLES: Record<string, Role> = {
  OWNER: "Owner",
  ASSIGNEE: "Assignee",
  MEMBER: "Member",
};

export type Role = "Owner" | "Assignee" | "Member";

export type OwnerApi = {
  id: string;
  username: string;
  display_name: string;
};

export type Owner = {
  id: string;
  username: string;
  displayName: string;
};

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

export interface ProjectDetailsApi {
  id: string;
  name: string;
  description?: string;
  skills?: string;
  role: Role;
  owner: OwnerApi;
  unassigned_tasks: number;
  ongoing_tasks: number;
  completed_tasks: number;
  abandoned_tasks: number;
  member_count: number;
  created_at: string;
  updated_at?: string;
}

export interface ProjectDetails {
  id: string;
  name: string;
  description?: string;
  skills?: string;
  role: Role;
  owner: Owner;
  unassignedTasks: number;
  ongoingTasks: number;
  completedTasks: number;
  abandonedTasks: number;
  memberCount: number;
  createdAt: string;
  updatedAt?: string;
}

export const MapProjectDetails = (p: ProjectDetailsApi): ProjectDetails => ({
  id: p.id,
  name: p.name,
  description: p.description,
  skills: p.skills,
  owner: {
    id: p.owner.id,
    username: p.owner.username,
    displayName: p.owner.display_name,
  },
  role: p.role,
  unassignedTasks: p.unassigned_tasks,
  ongoingTasks: p.ongoing_tasks,
  completedTasks: p.completed_tasks,
  abandonedTasks: p.abandoned_tasks,
  memberCount: p.member_count,
  createdAt: p.created_at,
  updatedAt: p.updated_at,
});
