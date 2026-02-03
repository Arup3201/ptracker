import type { TaskStatus } from "./task";

export interface DashboardProjectApi {
  id: string;
  name: string;
  unassigned_tasks: number;
  ongoing_tasks: number;
  completed_tasks: number;
  created_at: string;
  updated_at?: string;
}

export interface DashboardProject {
  id: string;
  name: string;
  unassignedTasks: number;
  ongoingTasks: number;
  completedTasks: number;
  createdAt: string;
  updatedAt?: string;
}

export const MapDashboardProject = (
  p: DashboardProjectApi,
): DashboardProject => ({
  id: p.id,
  name: p.name,
  unassignedTasks: p.unassigned_tasks,
  ongoingTasks: p.ongoing_tasks,
  completedTasks: p.completed_tasks,
  createdAt: p.created_at,
  updatedAt: p.updated_at,
});

export interface DashboardTask {
  id: string;
  projectName: string;
  title: string;
  status: TaskStatus;
  createdAt: string;
  updatedAt?: string;
}

export interface DashboardTaskApi {
  id: string;
  project_name: string;
  title: string;
  status: TaskStatus;
  created_at: string;
  updated_at?: string;
}

export const MapDashboardTask = (t: DashboardTaskApi): DashboardTask => ({
  id: t.id,
  projectName: t.project_name,
  title: t.title,
  status: t.status,
  createdAt: t.created_at,
  updatedAt: t.updated_at,
});

export interface DashboardTasksResponse {
  tasks: DashboardTaskApi[];
}

export interface DashboardProjectsResponse {
  projects: DashboardProjectApi[];
}
