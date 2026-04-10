import { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { Button } from "../components/button";
import { TopBar } from "../components/topbar";
import { Input } from "../components/input";
import { ApiRequest } from "../api/request";
import {
  MapExploreProject,
  type ExploreProject,
  type ExploreProjectsApiResponse,
} from "../types/explore";
import { ROLES } from "../types/project";

function ProjectCard({ project }: { project: ExploreProject }) {
  const navigate = useNavigate();

  const handleView = () => {
    if (project.role === ROLES.OWNER || project.role == ROLES.MEMBER) {
      navigate(`/projects/${project.id}`);
    } else {
      navigate(`/explore/${project.id}`);
    }
  };

  return (
    <div className="flex flex-col rounded-lg border border-border bg-bg-surface p-4 shadow-sm hover:border-border-strong transition duration-fast">
      <div className="flex items-start justify-between gap-2">
        <h3 className="text-base font-semibold text-text-primary leading-snug">
          {project.name}
        </h3>
      </div>

      <p className="mt-2 text-sm text-text-secondary leading-relaxed line-clamp-3">
        {project.description}
      </p>

      {project.skills && (
        <span className="mt-2 block text-xs text-text-muted">
          Skills: {project.skills}
        </span>
      )}

      <div className="mt-4 pt-4 border-t border-border-muted">
        <Button variant="secondary" onClick={handleView}>
          View Project
        </Button>
      </div>
    </div>
  );
}

export default function ExploreProjectsPage() {
  const [projects, setProjects] = useState<ExploreProject[]>([]);
  const [query, setQuery] = useState("");

  const getProjects = async () => {
    try {
      const data = await ApiRequest<ExploreProjectsApiResponse>(
        "/public/projects",
        "GET",
        null,
      );
      if (data?.projects) {
        const projects = data.projects.map(MapExploreProject);
        setProjects(projects || []);
      }
    } catch (err) {
      console.error(err);
    }
  };

  useEffect(() => {
    getProjects();
  }, []);

  const filteredProjects = projects.filter(
    (project) =>
      project.name.toLowerCase().includes(query.toLowerCase()) ||
      project.description.toLowerCase().includes(query.toLowerCase()),
  );

  return (
    <>
      <TopBar title="Explore Projects" />

      <div className="flex-1 overflow-y-auto p-6 space-y-6">
        <div className="space-y-4">
          <p className="text-sm text-text-secondary">
            Browse active projects and decide where to collaborate.
          </p>
          <div className="w-72">
            <Input
              placeholder="Search projects..."
              onChange={(text) => setQuery(text)}
            />
          </div>
        </div>

        {filteredProjects.length === 0 ? (
          <p className="text-sm text-text-muted text-center py-12">
            No projects found.
          </p>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {filteredProjects.map((project) => (
              <ProjectCard key={project.id} project={project} />
            ))}
          </div>
        )}
      </div>
    </>
  );
}
