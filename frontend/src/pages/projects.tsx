import type { Project } from "../types/project.ts";
import { useMemo, useState } from "react";
import { TopBar } from "../components/topbar.tsx";
import {
  Table,
  TableHeader,
  TableBody,
  TableRow,
  TableHead,
  TableCell,
} from "../components/table.tsx";
import { Button } from "../components/button";
import { Input } from "../components/input.tsx";

export function ProjectsPage() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [query, setQuery] = useState("");

  const filteredProjects = useMemo(() => {
    const q = query.trim().toLowerCase();

    if (!q) return projects;

    return projects.filter((project) => project.name.toLowerCase().includes(q));
  }, [query]);

  return (
    <>
      <TopBar title="Projects" actions={<Button>New Project</Button>} />

      <div className="flex-1 p-4 space-y-4">
        <div className="w-90">
          <Input
            placeholder="Search Projects"
            onChange={(text) => setQuery(text)}
          />
        </div>

        <Table>
          <TableHeader>
            <TableHead>Project</TableHead>
            <TableHead>Role</TableHead>
            <TableHead>Tasks</TableHead>
            <TableHead align="right">Updated</TableHead>
          </TableHeader>

          <TableBody>
            {filteredProjects.length === 0 && (
              <tr>
                <td
                  colSpan={4}
                  className="px-3 py-6 text-center text-sm text-(--text-muted)"
                >
                  No projects found
                </td>
              </tr>
            )}

            {filteredProjects.map((project) => (
              <TableRow
                key={project.id}
                onClick={() => {
                  // navigate(`/projects/${project.id}`)
                }}
                className="cursor-pointer"
              >
                <TableCell>
                  <span className="font-medium">{project.name}</span>
                </TableCell>

                <TableCell muted>{project.role}</TableCell>

                <TableCell muted>{renderTaskSignal(project)}</TableCell>

                <TableCell align="right" muted>
                  {project.updatedAt}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>
    </>
  );
}

function renderTaskSignal(project: Project) {
  const parts: string[] = [];

  if (project.ongoingTasks > 0) {
    parts.push(`${project.ongoingTasks} ongoing`);
  }

  if (project.unassignedTasks > 0) {
    parts.push(`${project.unassignedTasks} unassigned`);
  }

  return parts.length > 0 ? parts.join(" · ") : "—";
}
