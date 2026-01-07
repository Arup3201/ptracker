type TaskStatus = "Unassigned" | "Ongoing" | "Completed" | "Abandoned";

export interface Task {
  id: string;
  title: string;
  description: string;
  assignee?: string;
  status: TaskStatus;
  createdAt: string;
  updatedAt?: string;
}

export interface TaskApi {
  id: string;
  title: string;
  description: string;
  assignee?: string;
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
  assignee: t.assignee,
  status: t.status,
  createdAt: t.created_at,
  updatedAt: t.updated_at,
});
