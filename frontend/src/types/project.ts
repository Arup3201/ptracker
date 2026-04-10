export const ROLES: Record<string, Role> = {
  OWNER: "Owner",
  MEMBER: "Member",
};

export type Role = "Owner" | "Member";

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
  updated_at: string;
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
  updatedAt: string;
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
  role: Role;
  owner: AvatarApi;
  unassigned_tasks: number;
  ongoing_tasks: number;
  completed_tasks: number;
  abandoned_tasks: number;
  members_count: number;
  created_at: string;
  updated_at: string;
}

export interface ProjectDetails {
  id: string;
  name: string;
  description?: string;
  skills?: string;
  role: Role;
  owner: Avatar;
  unassignedTasks: number;
  ongoingTasks: number;
  completedTasks: number;
  abandonedTasks: number;
  membersCount: number;
  createdAt: string;
  updatedAt: string;
}

export const MapProjectDetails = (p: ProjectDetailsApi): ProjectDetails => ({
  id: p.id,
  name: p.name,
  description: p.description,
  skills: p.skills,
  owner: {
    userId: p.owner.user_id,
    username: p.owner.username,
    displayName: p.owner.display_name,
    email: p.owner.email,
    avatarUrl: p.owner.avatar_url,
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
  status: JoinStatus;
  createdAt: string;
  avatar: Avatar;
}

export interface JoinRequestApi {
  project_id: string;
  status: JoinStatus;
  created_at: string;
  avatar: AvatarApi;
}

export const MapJoinRequest = (request: JoinRequestApi): JoinRequest => ({
  projectId: request.project_id,
  status: request.status,
  createdAt: request.created_at,
  avatar: MapAvatar(request.avatar),
});

export interface JoinRequestsResponseApi {
  join_requests: JoinRequestApi[];
}

export interface Member {
  projectId: string;
  role: string;
  createdAt: string;
  updatedAt: string;
  avatar: Avatar;
}

export interface MemberApi {
  project_id: string;
  role: string;
  created_at: string;
  updated_at: string;
  avatar: AvatarApi;
}

export interface MembersResponse {
  members: MemberApi[];
}

export const MapMember = (m: MemberApi): Member => ({
  projectId: m.project_id,
  role: m.role,
  createdAt: m.created_at,
  updatedAt: m.updated_at,
  avatar: MapAvatar(m.avatar),
});
