import { useEffect, useMemo, useState } from "react";
import { useNavigate } from "react-router";

import {
  MapProject,
  type ProjectsApiResponse,
  type ProjectSummary,
} from "../types/project.ts";
import { ApiFetch } from "../utils/api";
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
import { CreateProjectModal } from "../components/create-project.tsx";
import { renderLocalTime, renderTaskSignal } from "../utils.ts";

export function ProjectsPage() {
  const navigate = useNavigate();

  const [showModal, setShowModal] = useState(false);
  const [projects, setProjects] = useState<ProjectSummary[]>([]);
  const [query, setQuery] = useState("");

  const getProjects = async () => {
    try {
      const response = await ApiFetch("/projects");
      if (response.ok) {
        const respondeData = await response.json();
        const data: ProjectsApiResponse = respondeData.data;
        if (data?.projects) {
          setProjects(data.projects.map(MapProject) || []);
        }
      } else {
        throw new Error("Failed to get projects.");
      }
    } catch (err) {
      console.log(err);
    }
  };

  useEffect(() => {
    getProjects();
  }, []);

  const filteredProjects = useMemo(() => {
    const q = query.trim().toLowerCase();

    if (!q) return projects;

    return projects.filter((project) => project.name.toLowerCase().includes(q));
  }, [query, projects]);

  return (
    <>
      <TopBar
        title="Projects"
        actions={
          <Button onClick={() => setShowModal(true)}>New Project</Button>
        }
      />

      <div className="flex-1 overflow-y-auto p-6 space-y-4">
        <div className="w-72">
          <Input
            placeholder="Search projects..."
            onChange={(e) => setQuery(e.target.value)}
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
            {filteredProjects.length === 0 ? (
              <tr>
                <td
                  colSpan={4}
                  className="px-4 py-8 text-center text-sm text-text-muted"
                >
                  No projects found
                </td>
              </tr>
            ) : (
              filteredProjects.map((project) => (
                <TableRow
                  key={project.id}
                  onClick={() => navigate(`/projects/${project.id}`)}
                  className="cursor-pointer"
                >
                  <TableCell>
                    <span className="font-medium text-text-primary hover:text-primary transition duration-fast">
                      {project.name}
                    </span>
                  </TableCell>
                  <TableCell muted>{project.role}</TableCell>
                  <TableCell muted>{renderTaskSignal(project)}</TableCell>
                  <TableCell align="right" muted>
                    {project.updatedAt && renderLocalTime(project.updatedAt)}
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      <CreateProjectModal
        open={showModal}
        onClose={() => setShowModal(false)}
      />
    </>
  );
}
