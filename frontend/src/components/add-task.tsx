import { useState } from "react";
import clsx from "clsx";
import { Modal } from "./modal";
import { Button } from "./button";
import { ApiRequest } from "../api/request";
import AssigneeSelector from "./assignee-selector";
import type { Member } from "../types/project";

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
      assignee: [],
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
        <form onSubmit={handleSubmit} className="px-4 py-4 flex flex-col gap-3">
          {/* Title */}
          <div className="flex flex-col gap-1">
            <label className="text-[12px] font-medium text-(--text-primary)">
              Title <span className="text-(--danger)">*</span>
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

          {/* Description */}
          <div className="flex flex-col gap-1">
            <label className="text-[12px] font-medium text-(--text-primary)">
              Description
            </label>

            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Add more context if needed"
              rows={3}
              className={clsx(
                "rounded-xs bg-(--bg-surface) px-2 py-1 text-sm text-(--text-primary)",
                "border border-(--border-default) outline-none resize-none",
                "focus:border-(--primary)",
              )}
            />
          </div>

          {/* Assignees */}
          <div className="flex flex-col gap-1">
            <AssigneeSelector members={members} />
          </div>

          {/* Status (optional / advanced) */}
          <div className="flex flex-col gap-1">
            <label className="text-[12px] font-medium text-(--text-primary)">
              Status
            </label>

            <select
              value={status}
              onChange={(e) => setStatus(e.target.value)}
              className={clsx(
                "h-8 rounded-xs bg-(--bg-surface) px-2 text-sm text-(--text-primary)",
                "border border-(--border-default) outline-none",
                "focus:border-(--primary)",
              )}
            >
              <option value="Unassigned">Unassigned</option>
              <option value="Ongoing">Ongoing</option>
              <option value="Completed">Completed</option>
              <option value="Abandoned">Abandoned</option>
            </select>
          </div>

          {/* Actions */}
          <div className="flex justify-end gap-2 px-4 py-3 border-t border-(--border-muted)">
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
