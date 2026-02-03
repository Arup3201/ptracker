import { MapMember, type Member, type MemberApi } from "./project";

export const TASK_STATUS: Record<string, TaskStatus> = {
  UNASSIGNED: "Unassigned",
  ONGOING: "Ongoing",
  COMPLETED: "Completed",
  ABANDONED: "Abandoned",
};

export type TaskStatus = "Unassigned" | "Ongoing" | "Completed" | "Abandoned";

export interface Task {
  id: string;
  title: string;
  description: string;
  assignees: Member[];
  status: TaskStatus;
  createdAt: string;
  updatedAt?: string;
}

export interface TaskApi {
  id: string;
  title: string;
  description: string;
  assignees: MemberApi[];
  status: TaskStatus;
  created_at: string;
  updated_at?: string;
}

export interface TasksResponseApi {
  tasks: TaskApi[];
  page: number;
  limit: number;
  has_next: boolean;
}

export const MapTask = (t: TaskApi): Task => ({
  id: t.id,
  title: t.title,
  description: t.description,
  assignees: t.assignees.map(MapMember),
  status: t.status,
  createdAt: t.created_at,
  updatedAt: t.updated_at,
});

export interface TaskDetailApi {
  id: string;
  title: string;
  description: string;
  assignees: MemberApi[];
  status: TaskStatus;
  created_at: string;
  updated_at?: string;
}

export interface TaskDetails {
  id: string;
  title: string;
  description: string;
  assignees: Member[];
  status: TaskStatus;
  createdAt: string;
  updatedAt?: string;
}

export const MapTaskDetails = (t: TaskDetailApi): TaskDetails => ({
  id: t.id,
  title: t.title,
  description: t.description,
  assignees: t.assignees.map(MapMember),
  status: t.status,
  createdAt: t.created_at,
  updatedAt: t.updated_at,
});

export interface TaskComment {
  id: string;
  projectId: string;
  taskId: string;
  user: Member;
  content: string;
  createdAt: string;
  updatedAt?: string;
}

export interface TaskCommentApi {
  id: string;
  project_id: string;
  task_id: string;
  user: MemberApi;
  content: string;
  created_at: string;
  updated_at?: string;
}

export const MapTaskComment = (c: TaskCommentApi): TaskComment => ({
  id: c.id,
  projectId: c.project_id,
  taskId: c.task_id,
  user: MapMember(c.user),
  content: c.content,
  createdAt: c.created_at,
  updatedAt: c.updated_at,
});

export interface TaskCommentsResponseApi {
  comments: TaskCommentApi[];
  page: number;
  limit: number;
  has_next: boolean;
}
