import { useState } from "react";
import clsx from "clsx";
import { Modal } from "./modal";
import { Button } from "./button";
import { ApiFetch } from "../utils/api";
import { useNavigate } from "react-router";

type CreateProjectModalProps = {
  open: boolean;
  onClose: () => void;
};

export const CreateProjectModal = ({
  open,
  onClose,
}: CreateProjectModalProps) => {
  const navigate = useNavigate();

  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [skills, setSkills] = useState("");
  const [error, setError] = useState<string | null>(null);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();

    if (!name.trim()) {
      setError("Project name is required");
      return;
    }

    const parsedSkills = skills
      .split(",")
      .map((s) => s.trim())
      .filter(Boolean);

    const projectId = await ApiFetch("/projects", {
      method: "POST",
      body: JSON.stringify({
        name: name.trim(),
        description: description.trim() || undefined,
        skills: parsedSkills.join(", "),
      }),
    });

    onClose();
    if (projectId) {
      navigate("/projects/" + projectId);
    }
  }

  return (
    <Modal
      open={open}
      title="Create Project"
      body={
        <form onSubmit={handleSubmit} className="flex flex-col gap-4 p-4">
          {/* Project name */}
          <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium text-text-primary">
              Project name <span className="text-danger">*</span>
            </label>
            <input
              autoFocus
              value={name}
              onChange={(e) => {
                setName(e.target.value);
                setError(null);
              }}
              placeholder="e.g. Internal PM Tool"
              className={clsx(
                "h-9 w-full rounded-md bg-bg-elevated px-3 text-sm text-text-primary",
                "border outline-none transition duration-fast",
                "placeholder:text-text-muted",
                error
                  ? "border-danger focus:shadow-focus-danger"
                  : "border-border focus:border-primary focus:shadow-focus-primary",
              )}
            />
            {error && <span className="text-xs text-danger">{error}</span>}
          </div>

          {/* Description */}
          <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium text-text-primary">
              Description
            </label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Short description of what this project is about"
              rows={3}
              className={clsx(
                "w-full rounded-md bg-bg-elevated px-3 py-2 text-sm text-text-primary",
                "border border-border outline-none resize-none transition duration-fast",
                "placeholder:text-text-muted",
                "focus:border-primary focus:shadow-focus-primary",
              )}
            />
          </div>

          {/* Skills */}
          <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium text-text-primary">
              Skills
            </label>
            <input
              value={skills}
              onChange={(e) => setSkills(e.target.value)}
              placeholder="C, Java, Python"
              className={clsx(
                "h-9 w-full rounded-md bg-bg-elevated px-3 text-sm text-text-primary",
                "border border-border outline-none transition duration-fast",
                "placeholder:text-text-muted",
                "focus:border-primary focus:shadow-focus-primary",
              )}
            />
            <span className="text-xs text-text-muted">
              Comma-separated values
            </span>
          </div>

          {/* Footer */}
          <div className="flex justify-end gap-2 border-t border-border-muted pt-4 mt-1">
            <Button type="button" variant="secondary" onClick={onClose}>
              Cancel
            </Button>
            <Button type="submit" disabled={!name.trim()}>
              Create project
            </Button>
          </div>
        </form>
      }
    />
  );
};
