import { type ReactNode, useEffect, useState } from "react";

import { ROLES, type Member, type Role } from "../types/project";
import {
  MapTaskComment,
  MapTaskDetails,
  TASK_STATUS,
  type TaskComment,
  type TaskCommentsResponseApi,
  type TaskDetailApi,
} from "../types/task";
import { Drawer } from "../components/drawer";
import { Button } from "../components/button";
import { ApiFetch } from "../utils/api";
import AssigneeSelector from "../components/assignee-selector";
import { StatusSelector } from "../components/status-selector";
import { useAuth } from "../context/auth";

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
  role: Role;
}) {
  const canEditAll = role === ROLES.OWNER;

  const [editMode, setEditMode] = useState(false);
  const [dirty, setDirty] = useState(false);

  const { user } = useAuth();

  const [title, setTitle] = useState<string>("");
  const [description, setDescription] = useState<string>("");
  const [status, setStatus] = useState<string>("Unassigned");

  const [editedTitle, setEditedTitle] = useState<string>("");
  const [editedDescription, setEditedDescription] = useState<string>("");
  const [editedStatus, setEditedStatus] = useState<string>("Unassigned");

  const [initialAssignees, setInitialAssignees] = useState<Member[]>([]);
  const [currentAssignees, setCurrentAssignees] = useState<Member[]>([]);

  const [comments, setComments] = useState<TaskComment[]>([]);

  const [comment, setComment] = useState<string>("");

  async function getProjectTask(projectId: string, taskId: string) {
    try {
      const response = await ApiFetch(`/projects/${projectId}/tasks/${taskId}`);
      if (response.ok) {
        const { data } = await response.json();
        const task: TaskDetailApi = data;
        if (task) {
          const taskDetails = MapTaskDetails(task);
          setTitle(taskDetails.title || "");
          setEditedTitle(taskDetails.title || "");
          setDescription(taskDetails.description || "");
          setEditedDescription(taskDetails.description || "");
          setStatus(taskDetails.status || TASK_STATUS.UNASSIGNED);
          setEditedStatus(taskDetails.status || TASK_STATUS.UNASSIGNED);

          setInitialAssignees(taskDetails.assignees || []);
          setCurrentAssignees(taskDetails.assignees || []);
        }
      } else {
        throw new Error("Failed to get task details.");
      }
    } catch (err) {
      console.error(err);
    }
  }

  async function getTaskComment(projectId: string, taskId: string) {
    try {
      const response = await ApiFetch(
        `/projects/${projectId}/tasks/${taskId}/comments`,
      );
      if (response.ok) {
        const responseData = await response.json();
        const data: TaskCommentsResponseApi = responseData.data;
        if (data) {
          setComments(data.comments.map(MapTaskComment));
        }
      } else {
        throw new Error("Failed to get task comments.");
      }
    } catch (err) {
      console.error(err);
    }
  }

  useEffect(() => {
    if (projectId && taskId) {
      getProjectTask(projectId, taskId);
      getTaskComment(projectId, taskId);
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

  function AssigneeDiff(a1: Member[], a2: Member[]): string[] {
    return a1
      .filter(
        (a) =>
          a2.findIndex((ia) => ia.avatar.userId === a.avatar.userId) === -1,
      )
      .map((a) => a.avatar.userId);
  }

  async function handleEditSave() {
    const assigneesToAdd = AssigneeDiff(currentAssignees, initialAssignees);
    const assigneesToRemove = AssigneeDiff(initialAssignees, currentAssignees);

    function nullFilter(editedValue: string, originalValue: string) {
      if (editedValue !== originalValue) return editedValue;

      return null;
    }

    const data = {
      title: nullFilter(editedTitle, title),
      description: nullFilter(editedDescription, description),
      status: nullFilter(editedStatus, status),

      assignees_to_add: assigneesToAdd.length > 0 ? assigneesToAdd : null,
      assignees_to_remove:
        assigneesToRemove.length > 0 ? assigneesToRemove : null,
    };

    try {
      const response = await ApiFetch(
        `/projects/${projectId}/tasks/${taskId}`,
        {
          method: "PUT",
          body: JSON.stringify(data),
        },
      );
      if (!response.ok) {
        throw new Error("Failed to update the task.");
      }
    } catch (err) {
      console.error(err);
    } finally {
      setEditMode(false);
      setDirty(false);
    }
  }

  function handleAssigneeEdit(assignees: string[]) {
    setCurrentAssignees(() => {
      return members.filter((m) =>
        assignees.find((a) => a === m.avatar.userId),
      );
    });
  }

  async function handleAddComment() {
    const data = {
      comment: comment,
      user_id: user?.userId,
    };
    try {
      const response = await ApiFetch(
        `/projects/${projectId}/tasks/${taskId}/comments`,
        {
          method: "POST",
          body: JSON.stringify(data),
        },
      );
      if (!response.ok) {
        throw new Error("Failed to update the task.");
      }
    } catch (err) {
      console.error(err);
    } finally {
      setComment("");
      // Refresh comments
      if (projectId && taskId) {
        getProjectTask(projectId, taskId);
      }
    }
  }

  const isAssignee =
    initialAssignees.findIndex((a) => a.avatar.userId === user?.userId) !== -1;

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
                setEditedTitle(title);
                setEditedDescription(description);
                setEditedStatus(status);
              }}
            >
              Cancel
            </Button>
            <Button
              disabled={
                editedTitle === title &&
                editedDescription === description &&
                editedStatus === status &&
                AssigneeDiff(currentAssignees, initialAssignees).length === 0
              }
              onClick={handleEditSave}
            >
              Save
            </Button>
          </div>
        ) : null
      }
    >
      {/* Title + meta */}
      <div className="space-y-2 mb-6">
        {editMode && (canEditAll || isAssignee) ? (
          <input
            value={editedTitle}
            onChange={(e) => {
              setEditedTitle(e.target.value);
              setDirty(true);
            }}
            className="w-full h-9 rounded-md bg-bg-elevated px-3 text-sm text-text-primary border border-border outline-none transition duration-fast focus:border-primary focus:shadow-focus-primary"
          />
        ) : (
          <h3 className="text-base font-semibold text-text-primary tracking-snug">
            {title}
          </h3>
        )}

        {!editMode && (
          <p className="text-xs text-text-muted">
            Status: <span className="text-text-secondary">{status}</span>
            {" · "}
            Assignee:{" "}
            <span className="text-text-secondary">
              {currentAssignees.map((a) => a.avatar.username).join(", ") ||
                "Unassigned"}
            </span>
          </p>
        )}

        {!editMode && (canEditAll || isAssignee) && (
          <Button variant="secondary" onClick={() => setEditMode(true)}>
            Edit
          </Button>
        )}
      </div>

      {/* Description */}
      <section className="mb-6 space-y-1.5">
        <h4 className="text-xs font-medium text-text-muted uppercase tracking-wider">
          Description
        </h4>
        {editMode && (canEditAll || isAssignee) ? (
          <textarea
            rows={4}
            value={editedDescription}
            onChange={(e) => {
              setEditedDescription(e.target.value);
              setDirty(true);
            }}
            className="w-full rounded-md bg-bg-elevated px-3 py-2 text-sm text-text-primary border border-border outline-none resize-none transition duration-fast focus:border-primary focus:shadow-focus-primary"
          />
        ) : (
          <p className="text-sm text-text-secondary leading-relaxed">
            {description || "No description provided."}
          </p>
        )}
      </section>

      {/* Edit fields */}
      {editMode && (canEditAll || isAssignee) && (
        <section className="mb-6 space-y-4">
          <div className="flex flex-col gap-1.5">
            <StatusSelector status={editedStatus} onChange={setEditedStatus} />
          </div>
          <div className="flex flex-col gap-1.5">
            <AssigneeSelector
              initialAssignees={initialAssignees.map((a) => a.avatar.userId)}
              members={members}
              onChange={handleAssigneeEdit}
              isDisabled={!canEditAll}
            />
          </div>
        </section>
      )}

      {/* Comments */}
      <section className="space-y-4">
        <h4 className="text-xs font-medium text-text-muted uppercase tracking-wider">
          Comments
        </h4>

        {!editMode && (
          <div className="flex flex-col gap-2">
            <textarea
              placeholder="Add a comment…"
              value={comment}
              onChange={(e) => setComment(e.target.value)}
              rows={2}
              className="w-full rounded-md bg-bg-elevated px-3 py-2 text-sm text-text-primary border border-border outline-none resize-none transition duration-fast placeholder:text-text-muted focus:border-primary focus:shadow-focus-primary"
            />
            <Button onClick={handleAddComment} className="self-end">
              Send
            </Button>
          </div>
        )}

        <div className="space-y-4">
          {comments.length === 0 ? (
            <p className="text-sm text-text-muted">No comments yet.</p>
          ) : (
            comments.map((comment, index) => (
              <Comment
                key={index}
                author={comment.avatar.displayName || comment.avatar.username}
                time={comment.createdAt}
              >
                {comment.content}
              </Comment>
            ))
          )}
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
  const createdAtMs = new Date(time).getTime();
  const diffMs = Date.now() - createdAtMs;

  const sec = Math.floor(diffMs / 1000);
  const min = Math.floor(sec / 60);
  const hr = Math.floor(min / 60);
  const day = Math.floor(hr / 24);
  const month = Math.floor(day / 30);
  const year = Math.floor(day / 365);

  let timeAgo = "";
  if (year >= 1) timeAgo = `${year} year${year > 1 ? "s" : ""}`;
  else if (month >= 1) timeAgo = `${month} month${month > 1 ? "s" : ""}`;
  else if (day >= 1) timeAgo = `${day} day${day > 1 ? "s" : ""}`;
  else if (hr >= 1) timeAgo = `${hr} hr${hr > 1 ? "s" : ""}`;
  else if (min >= 1) timeAgo = `${min} min${min > 1 ? "s" : ""}`;
  else timeAgo = "just now";

  return (
    <div className="space-y-1">
      <div className="flex items-center gap-1.5">
        <div className="h-5 w-5 rounded-full bg-bg-elevated border border-border flex items-center justify-center text-[10px] font-semibold text-primary shrink-0">
          {author?.charAt(0).toUpperCase()}
        </div>
        <span className="text-xs font-medium text-text-primary">{author}</span>
        <span className="text-xs text-text-muted">· {timeAgo}</span>
      </div>
      <p className="text-sm text-text-secondary leading-relaxed pl-6">
        {children}
      </p>
    </div>
  );
}
