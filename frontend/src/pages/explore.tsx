import { useState } from "react";
import type { Project } from "../types/project";
import { Button } from "../components/button";
import { TopBar } from "../components/topbar";
import { Input } from "../components/input";

function ProjectCard({ project }: { project: Project }) {
  return (
    <div className="rounded-sm border border-(--border-default) bg-(--bg-surface) p-4 shadow-[0_1px_2px_rgba(0,0,0,0.4)">
      <h3 className="text-[14px] font-semibold text-(--text-primary) leading-snug">
        {project.name}
      </h3>
      <p className="mt-2 text-[13px] text-(--text-secondary) leading-relaxed line-clamp-3">
        {project.description}
      </p>
      <div className="mt-4">
        <Button variant="secondary">View Project</Button>
      </div>
    </div>
  );
}

export default function ExploreProjectsPage() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [query, setQuery] = useState("");

  const filteredProjects = projects.filter(
    (project) =>
      project.name.toLowerCase().includes(query.toLowerCase()) ||
      project.description.toLowerCase().includes(query.toLowerCase())
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
