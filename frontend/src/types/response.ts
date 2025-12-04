import type { ProjectRole } from "./project";
import type { TaskStatus } from "./task";

export interface MemberResponse {
    user_id: string;
    name: string;
    email: string;
    role: ProjectRole, 
    joined_at: string
}

export interface TaskResponse {
    id: string;
    name: string;
    description: string;
    status: TaskStatus;
    assignee: string;
    assignee_name: string;
    assignee_email: string;
}