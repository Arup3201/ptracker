type TaskStatus = "Unassigned" | "Ongoing" | "Completed" | "Abandoned";

export interface Task {
  id: string;
  title: string;
  assignee?: string;
  status: TaskStatus;
  createdAt: string;
  updatedAt?: string;
}
