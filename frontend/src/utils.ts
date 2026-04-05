import type { DashboardProject } from "./types/dashboard";
import type { ProjectSummary } from "./types/project";

export function renderTaskSignal(project: ProjectSummary | DashboardProject) {
  const parts: string[] = [];

  if (project.ongoingTasks > 0) {
    parts.push(`${project.ongoingTasks} ongoing`);
  }

  if (project.unassignedTasks > 0) {
    parts.push(`${project.unassignedTasks} unassigned`);
  }

  return parts.length > 0 ? parts.join(" · ") : "—";
}

export function renderLocalTime(utcTimestamp: string) {
  const local = new Date(utcTimestamp);
  return local.toLocaleString();
}
