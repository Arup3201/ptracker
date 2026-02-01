import { useState } from "react";
import clsx from "clsx";
import { Modal } from "./modal";
import { Button } from "./button";
import { ApiRequest } from "../api/request";
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

    const projectId = await ApiRequest<string>("/projects", "POST", {
      name: name.trim(),
      description: description.trim() || undefined,
      skills: parsedSkills.join(", "),
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
        <>
          <form
            onSubmit={handleSubmit}
            className="px-4 py-4 flex flex-col gap-3"
          >
            <div className="flex flex-col gap-1">
              <label className="text-[12px] font-medium text-(--text-primary)">
                Project name <span className="text-(--danger)">*</span>
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
                  "h-8 rounded-xs bg-(--bg-surface) px-2 text-sm text-(--text-primary)",
                  "border outline-none",
                  error
                    ? "border-(--danger)"
                    : "border-(--border-default) focus:border-(--primary)",
                )}
              />

              {error && (
                <span className="text-[11px] text-(--danger)">{error}</span>
              )}
            </div>

            <div className="flex flex-col gap-1">
              <label className="text-[12px] font-medium text-(--text-primary)">
                Description
              </label>

              <textarea
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="Short description of what this project is about"
                rows={3}
                className={clsx(
                  "rounded-xs bg-(--bg-surface) px-2 py-1 text-sm text-(--text-primary)",
                  "border border-(--border-default) outline-none resize-none",
                  "focus:border-(--primary)",
                )}
              />
            </div>

            <div className="flex flex-col gap-1">
              <label className="text-[12px] font-medium text-(--text-primary)">
                Skills
              </label>

              <input
                value={skills}
                onChange={(e) => setSkills(e.target.value)}
                placeholder="C, Java, Python"
                className={clsx(
                  "h-8 rounded-xs bg-(--bg-surface) px-2 text-sm text-(--text-primary)",
                  "border border-(--border-default) outline-none",
                  "focus:border-(--primary)",
                )}
              />

              <span className="text-[11px] text-(--text-muted)">
                Comma-separated values
              </span>
            </div>

            <div className="flex justify-end gap-2 px-4 py-3 border-t border-(--border-muted)">
              <Button type="button" variant="secondary" onClick={onClose}>
                Cancel
              </Button>

              <Button type="submit" disabled={!name.trim()}>
                Create project
              </Button>
            </div>
          </form>
        </>
      }
    />
  );
};
