export interface JoinRequest {
  id: string;
  name: string;
  note: string;
}

export interface ExploreProject {
  id: string;
  title: string;
  description: string;
  skills: string;
  role: string;
  createdAt: string;
  updatedAt: string;
}

export interface ExploreProjectApi {
  id: string;
  title: string;
  description: string;
  skills: string;
  role: string;
  created_at: string;
  updated_at: string;
}

export const MapExploreProject = (p: ExploreProjectApi): ExploreProject => ({
  id: p.id,
  title: p.title,
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
