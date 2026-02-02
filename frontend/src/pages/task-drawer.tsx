import { type ReactNode, useEffect, useState } from "react";

import { ROLES, type Member, type TaskRole } from "../types/project";
import { MapTaskDetails, TASK_STATUS, type TaskDetailApi } from "../types/task";
import { Drawer } from "../components/drawer";
import { Button } from "../components/button";
import { ApiRequest } from "../api/request";
import AssigneeSelector from "../components/assignee-selector";
import { StatusSelector } from "../components/status-selector";

export function TaskDrawer({
  open,
  taskId,
  projectId,
  members,
  onClose,
  role,
}: {
  open: boolean;
  taskId: string | null;
  projectId?: string;
  members: Member[];
  onClose: () => void;
  role: TaskRole;
}) {
  const canEditAll = role === ROLES.OWNER;
  const canEditPartial = role === ROLES.ASSIGNEE;

  const [editMode, setEditMode] = useState(false);
  const [dirty, setDirty] = useState(false);

  const [title, setTitle] = useState<string>("");
  const [description, setDescription] = useState<string>("");
  const [status, setStatus] = useState<string>("Unassigned");
  const [assignees, setAssignees] = useState<Member[]>([]);

  async function getProjectTask(projectId: string, taskId: string) {
    try {
      const data = await ApiRequest<TaskDetailApi>(
        `/projects/${projectId}/tasks/${taskId}`,
        "GET",
        null,
      );
      if (data) {
        const taskDetails = MapTaskDetails(data);
        setTitle(taskDetails.title || "");
        setDescription(taskDetails.description || "");
        setStatus(taskDetails.status || TASK_STATUS.UNASSIGNED);
        setAssignees(taskDetails.assignees || []);
      }
    } catch (err) {
      console.error(err);
    }
  }

  useEffect(() => {
    if (projectId && taskId) {
      getProjectTask(projectId, taskId);
    }
  }, [projectId, taskId]);

  const requestClose = () => {
    if (dirty && editMode) {
      if (!confirm("Discard changes?")) return;
    }
    setEditMode(false);
    setDirty(false);
    onClose();
  };

  return (
    <Drawer
      open={open}
      onClose={requestClose}
      title={title}
      footer={
        editMode ? (
          <div className="flex justify-end gap-2">
            <Button
              variant="secondary"
              onClick={() => {
                setEditMode(false);
                setDirty(false);
              }}
            >
              Cancel
            </Button>
            <Button
              onClick={() => {
                setEditMode(false);
                setDirty(false);
              }}
            >
              Save
            </Button>
          </div>
        ) : null
      }
    >
      {/* Header Metadata */}
      <div className="space-y-1 mb-4">
        {editMode ? (
          <input
            value={title}
            onChange={(e) => {
              setTitle(e.target.value);
              setDirty(true);
            }}
            className="w-full h-8 rounded-xs bg-(--bg-surface) px-2 text-sm border border-(--border-default) outline-none focus:border-(--primary)"
          />
        ) : (
          <h3 className="text-sm font-medium text-(--text-primary)">{title}</h3>
        )}

        {!editMode && (
          <div className="text-xs text-(--text-muted)">
            Status: {status} · Assignee:{" "}
            {assignees.map((a) => a.username).join(", ") || "Unassigned"}
          </div>
        )}

        {!editMode && (canEditAll || canEditPartial) && (
          <Button variant="secondary" onClick={() => setEditMode(true)}>
            Edit
          </Button>
        )}
      </div>

      {/* Description */}
      <section className="mb-6">
        <h4 className="text-xs font-medium text-(--text-primary) mb-1">
          Description
        </h4>

        {editMode ? (
          <textarea
            rows={4}
            value={description}
            onChange={(e) => {
              setDescription(e.target.value);
              setDirty(true);
            }}
            className="w-full rounded-xs bg-(--bg-surface) px-2 py-1 text-sm border border-(--border-default) outline-none resize-none focus:border-(--primary)"
          />
        ) : (
          <p className="text-sm text-(--text-secondary)">
            {description || "No description provided"}
          </p>
        )}
      </section>

      {/* Metadata Editing (Owner only) */}
      {editMode && canEditAll && (
        <section className="mb-6 space-y-3">
          <div className="flex flex-col gap-1">
            <StatusSelector status={status} onChange={setStatus} />
          </div>

          <div className="flex flex-col gap-1">
            <AssigneeSelector
              initialAssignees={assignees.map((a) => a.userId)}
              members={members}
            />
          </div>
        </section>
      )}

      {/* Comments */}
      <section>
        <h4 className="text-xs font-medium text-(--text-primary) mb-2">
          Comments
        </h4>

        <div className="space-y-3 mb-3">
          {!editMode && (
            <textarea
              placeholder="Add a comment…"
              rows={2}
              className="w-full rounded-xs bg-(--bg-surface) px-2 py-1 text-sm border border-(--border-default) outline-none resize-none focus:border-(--primary)"
            />
          )}

          <Comment author="Arup" time="2h ago">
            Please make sure this works with prod DB as well.
          </Comment>
          <Comment author="Rahul" time="1h ago">
            Working on it, will update soon.
          </Comment>
        </div>
      </section>
    </Drawer>
  );
}

function Comment({
  author,
  time,
  children,
}: {
  author: string;
  time: string;
  children: ReactNode;
}) {
  return (
    <div className="space-y-0.5">
      <div className="text-xs text-(--text-muted)">
        <span className="font-medium text-(--text-primary)">{author}</span> ·{" "}
        {time}
      </div>
      <p className="text-sm text-(--text-secondary)">{children}</p>
    </div>
  );
}
