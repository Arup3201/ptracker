import { useState } from "react";
import clsx from "clsx";
import { Modal } from "./modal";
import { Button } from "./button";
import { ApiRequest } from "../api/request";
import AssigneeSelector from "./assignee-selector";
import type { Member } from "../types/project";
import { StatusSelector } from "./status-selector";

type AddTaskModalProps = {
  projectId: string | undefined;
  open: boolean;
  members: Member[];
  onClose: () => void;
};

export const AddTaskModal = ({
  projectId,
  open,
  members,
  onClose,
}: AddTaskModalProps) => {
  const [title, setTitle] = useState<string | undefined>(undefined);
  const [description, setDescription] = useState<string | undefined>(undefined);
  const [status, setStatus] = useState<string>("Unassigned");
  const [assignees, setAssignees] = useState<string[]>([]);

  const [error, setError] = useState<string | null>(null);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();

    if (!title?.trim()) {
      setError("Task title is required");
      return;
    }

    if (!projectId) {
      console.error("No project id");
      return;
    }

    const data = {
      title: title.trim(),
      description: description?.trim() || undefined,
      assignees: assignees,
      status: status,
    };
    await ApiRequest(`/projects/${projectId}/tasks`, "POST", data);

    onClose();
  }

  return (
    <Modal
      open={open}
      title="New Task"
      body={
        <form onSubmit={handleSubmit} className="flex flex-col gap-4 p-4">
          {/* Title */}
          <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium text-text-primary">
              Title <span className="text-danger">*</span>
            </label>
            <input
              autoFocus
              value={title}
              onChange={(e) => {
                setTitle(e.target.value);
                setError(null);
              }}
              placeholder="e.g. Set up database migrations"
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
              placeholder="Add more context if needed"
              rows={3}
              className={clsx(
                "w-full rounded-md bg-bg-elevated px-3 py-2 text-sm text-text-primary",
                "border border-border outline-none resize-none transition duration-fast",
                "placeholder:text-text-muted",
                "focus:border-primary focus:shadow-focus-primary",
              )}
            />
          </div>

          {/* Assignees */}
          <div className="flex flex-col gap-1.5">
            <AssigneeSelector members={members} onChange={setAssignees} />
          </div>

          {/* Status */}
          <div className="flex flex-col gap-1.5">
            <StatusSelector status={status} onChange={setStatus} />
          </div>

          {/* Footer */}
          <div className="flex justify-end gap-2 border-t border-border-muted pt-4 mt-1">
            <Button type="button" variant="secondary" onClick={onClose}>
              Cancel
            </Button>
            <Button type="submit" disabled={!title?.trim()}>
              Create task
            </Button>
          </div>
        </form>
      }
    />
  );
};
