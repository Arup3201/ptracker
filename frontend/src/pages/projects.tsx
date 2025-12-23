import { useMemo, useState } from "react";
import { Sidebar } from "../components/sidebar.tsx";
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

/**
 * Temporary shape — replace with real data later
 */
interface Project {
  id: string;
  name: string;
  role: "Owner" | "Member";
  ongoingCount: number;
  unassignedCount: number;
  updatedAt: string;
}

const PROJECTS: Project[] = [
  {
    id: "1",
    name: "Payments Revamp",
    role: "Owner",
    ongoingCount: 3,
    unassignedCount: 1,
    updatedAt: "2h ago",
  },
  {
    id: "2",
    name: "Auth Refactor",
    role: "Member",
    ongoingCount: 1,
    unassignedCount: 0,
    updatedAt: "1d ago",
  },
  {
    id: "3",
    name: "Internal Tooling",
    role: "Owner",
    ongoingCount: 0,
    unassignedCount: 2,
    updatedAt: "4d ago",
  },
];

export function ProjectsPage() {
  const [query, setQuery] = useState("");

  const filteredProjects = useMemo(() => {
    const q = query.trim().toLowerCase();

    if (!q) return PROJECTS;

    return PROJECTS.filter((project) => project.name.toLowerCase().includes(q));
  }, [query]);

  return (
    <>
      <Sidebar />

      <main className="flex flex-1 flex-col">
        <TopBar title="Projects" actions={<Button>New Project</Button>} />

        <div className="flex-1 overflow-y-auto p-4 space-y-4">
          <div className="w-90">
            <input
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder="Search projects"
              className="h-8 w-full rounded-xs border border-(--border-default) bg-(--bg-surface) px-3 text-sm placeholder:text-(--text-muted) focus:border-(--primary) focus:outline-none"
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
      </main>
    </>
  );
}

function renderTaskSignal(project: {
  ongoingCount: number;
  unassignedCount: number;
}) {
  const parts: string[] = [];

  if (project.ongoingCount > 0) {
    parts.push(`${project.ongoingCount} ongoing`);
  }

  if (project.unassignedCount > 0) {
    parts.push(`${project.unassignedCount} unassigned`);
  }

  return parts.length > 0 ? parts.join(" · ") : "—";
}
