export type ProjectRole = "Owner" | "Memeber"

export interface Project {
  id: string;
  name: string;
  description: string;
  deadline: string;
  code: string;
}

export interface NewProjectData {
  name: string;
  description: string;
  deadline: string;
}