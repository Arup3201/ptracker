import type { Owner, OwnerApi } from "./project";

export type JoinStatus = "Not Requested" | "Pending" | "Accepted";

export const JOIN_STATUS: Record<string, JoinStatus> = {
  NOT_REQUESTED: "Not Requested",
  PENDING: "Pending",
  ACCEPTED: "Accepted",
};

export interface ExploreProject {
  id: string;
  name: string;
  description: string;
  skills: string;
  role: string;
  createdAt: string;
  updatedAt: string;
}

export interface ExploreProjectApi {
  id: string;
  name: string;
  description: string;
  skills: string;
  role: string;
  created_at: string;
  updated_at: string;
}

export const MapExploreProject = (p: ExploreProjectApi): ExploreProject => ({
  id: p.id,
  name: p.name,
  description: p.description,
  skills: p.skills,
  role: p.role,
  createdAt: p.created_at,
  updatedAt: p.updated_at,
});

export interface ExploreProjectsApiResponse {
  projects: ExploreProjectApi[];
  page: number;
  limit: number;
  has_next: boolean;
}

export interface ExploredProjectDetailsApi {
  id: string;
  name: string;
  description?: string;
  skills?: string;
  owner: OwnerApi;
  unassigned_tasks: number;
  ongoing_tasks: number;
  completed_tasks: number;
  abandoned_tasks: number;
  join_status: JoinStatus;
  created_at: string;
  updated_at?: string;
}

export interface ExploredProjectDetails {
  id: string;
  name: string;
  description?: string;
  skills?: string;
  owner: Owner;
  unassignedTasks: number;
  ongoingTasks: number;
  completedTasks: number;
  abandonedTasks: number;
  joinStatus: JoinStatus;
  createdAt: string;
  updatedAt?: string;
}

export const MapExploredProjectDetails = (p: ExploredProjectDetailsApi) => ({
  id: p.id,
  name: p.name,
  description: p.description,
  skills: p.skills,
  owner: {
    id: p.owner.id,
    username: p.owner.username,
    displayName: p.owner.display_name,
  },
  unassignedTasks: p.unassigned_tasks,
  ongoingTasks: p.ongoing_tasks,
  completedTasks: p.completed_tasks,
  abandonedTasks: p.abandoned_tasks,
  joinStatus: p.join_status,
  createdAt: p.created_at,
  updatedAt: p.updated_at,
});
