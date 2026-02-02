export const ROLES: Record<string, TaskRole> = {
  OWNER: "Owner",
  MEMBER: "Member",
  ASSIGNEE: "Assignee",
};

export type TaskRole = "Owner" | "Member" | "Assignee";

export type ProjectRole = "Owner" | "Member";

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
  role: ProjectRole;
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
  role: ProjectRole;
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

export interface ProjectDetailsApi {
  id: string;
  name: string;
  description?: string;
  skills?: string;
  role: ProjectRole;
  owner: OwnerApi;
  unassigned_tasks: number;
  ongoing_tasks: number;
  completed_tasks: number;
  abandoned_tasks: number;
  members_count: number;
  created_at: string;
  updated_at?: string;
}

export interface ProjectDetails {
  id: string;
  name: string;
  description?: string;
  skills?: string;
  role: ProjectRole;
  owner: Owner;
  unassignedTasks: number;
  ongoingTasks: number;
  completedTasks: number;
  abandonedTasks: number;
  membersCount: number;
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
  membersCount: p.members_count,
  createdAt: p.created_at,
  updatedAt: p.updated_at,
});

type JoinStatus = "Pending" | "Accepted" | "Rejected";

export interface JoinRequest {
  projectId: string;
  userId: string;
  username: string;
  displayName: string;
  email: string;
  avatarURL: string;
  isActive: boolean;
  status: JoinStatus;
  createdAt: string;
}

export interface JoinRequestApi {
  project_id: string;
  user_id: string;
  username: string;
  display_name: string;
  email: string;
  avatar_url: string;
  is_active: boolean;
  status: JoinStatus;
  created_at: string;
}

export const MapJoinRequest = (request: JoinRequestApi): JoinRequest => ({
  projectId: request.project_id,
  userId: request.user_id,
  username: request.username,
  displayName: request.display_name,
  email: request.email,
  avatarURL: request.avatar_url,
  isActive: request.is_active,
  status: request.status,
  createdAt: request.created_at,
});

export interface JoinRequestsResponseApi {
  join_requests: JoinRequestApi[];
}

export interface Member {
  projectId: string;
  id: string;
  username: string;
  displayName: string;
  email: string;
  avatarUrl: string;
  isActive: boolean;
  role: ProjectRole;
  createdAt: string;
  updatedAt: string;
}

export interface MemberApi {
  project_id: string;
  id: string;
  username: string;
  display_name: string;
  email: string;
  avatar_url: string;
  is_active: boolean;
  role: ProjectRole;
  created_at: string;
  updated_at: string;
}

export interface MembersResponse {
  members: MemberApi[];
}

export const MapMember = (m: MemberApi): Member => ({
  projectId: m.project_id,
  id: m.id,
  username: m.username,
  displayName: m.display_name,
  email: m.email,
  avatarUrl: m.avatar_url,
  isActive: m.is_active,
  role: m.role,
  createdAt: m.created_at,
  updatedAt: m.updated_at,
});
