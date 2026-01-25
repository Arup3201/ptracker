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
    <div className="rounded-sm border border-(--border-default) bg-(--bg-surface) p-4 shadow-[0_1px_2px_rgba(0,0,0,0.4)]">
      <div className="flex items-start justify-between gap-2">
        <h3 className="text-[14px] font-semibold text-(--text-primary) leading-snug">
          {project.title}
        </h3>

        {project.role !== "User" && (
          <span className="text-[11px] font-medium text-(--text-muted) border border-(--border-muted) rounded-xs px-2 py-[2px]">
            {project.role}
          </span>
        )}
      </div>

      <p className="mt-2 text-[13px] text-(--text-secondary) leading-relaxed line-clamp-3">
        {project.description}
      </p>

      <span className="mt-2 block text-[12px] text-(--text-muted)">
        Skills: {project.skills}
      </span>

      <div className="mt-4">
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
        "/explore/projects",
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
      project.title.toLowerCase().includes(query.toLowerCase()) ||
      project.description.toLowerCase().includes(query.toLowerCase()),
  );

  return (
    <>
      <TopBar title="Explore Projects" />

      <div className="p-4 space-y-4">
        <p className="text-[12px] text-(--text-secondary)">
          Browse active projects and decide where to collaborate.
        </p>

        <Input
          placeholder="Search projects"
          onChange={(text) => setQuery(text)}
        />
      </div>
      <div className="p-4 overflow-y-auto">
        {filteredProjects.length === 0 ? (
          <p className="text-[12px] text-(--text-muted) text-center">
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
