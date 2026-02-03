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
import { ApiRequest } from "../api/request";
import AssigneeSelector from "../components/assignee-selector";
import { StatusSelector } from "../components/status-selector";
import { useCurrentUser } from "../hooks/current_user";

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

  const currentUser = useCurrentUser();

  const [title, setTitle] = useState<string>("");
  const [description, setDescription] = useState<string>("");
  const [status, setStatus] = useState<string>("Unassigned");
  const [initialAssignees, setInitialAssignees] = useState<Member[]>([]);
  const [currentAssignees, setCurrentAssignees] = useState<Member[]>([]);
  const [comments, setComments] = useState<TaskComment[]>([]);

  const [comment, setComment] = useState<string>("");

  async function getProjectTask(projectId: string, taskId: string) {
    try {
      const taskDetailsResponse = await ApiRequest<TaskDetailApi>(
        `/projects/${projectId}/tasks/${taskId}`,
        "GET",
        null,
      );
      if (taskDetailsResponse) {
        const taskDetails = MapTaskDetails(taskDetailsResponse);
        setTitle(taskDetails.title || "");
        setDescription(taskDetails.description || "");
        setStatus(taskDetails.status || TASK_STATUS.UNASSIGNED);

        setInitialAssignees(taskDetails.assignees || []);
        setCurrentAssignees(taskDetails.assignees || []);
      }

      const commentsResponse = await ApiRequest<TaskCommentsResponseApi>(
        `/projects/${projectId}/tasks/${taskId}/comments`,
        "GET",
        null,
      );
      if (commentsResponse) {
        setComments(commentsResponse.comments.map(MapTaskComment));
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

  async function handleEditSave() {
    const assigneesToAdd = currentAssignees
      .filter(
        (a) =>
          initialAssignees.findIndex((ia) => ia.userId === a.userId) === -1,
      )
      .map((a) => a.userId);
    const assigneesToRemove = initialAssignees
      .filter(
        (a) =>
          currentAssignees.findIndex((ca) => ca.userId === a.userId) === -1,
      )
      .map((a) => a.userId);

    const data = {
      title: title,
      description: description,
      status: status,
      assignees_to_add: assigneesToAdd,
      assignees_to_remove: assigneesToRemove,
    };

    try {
      await ApiRequest(`/projects/${projectId}/tasks/${taskId}`, "PUT", data);
    } catch (err) {
      console.error(err);
    } finally {
      setEditMode(false);
      setDirty(false);
    }
  }

  function handleAssigneeEdit(assignees: string[]) {
    setCurrentAssignees(() => {
      return members.filter((m) => assignees.find((a) => a === m.userId));
    });
  }

  async function handleAddComment() {
    const data = {
      comment: comment,
      user_id: currentUser?.id,
    };
    try {
      await ApiRequest(
        `/projects/${projectId}/tasks/${taskId}/comments`,
        "POST",
        data,
      );
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
    initialAssignees.findIndex((a) => a.userId === currentUser?.id) !== -1;

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
            <Button onClick={handleEditSave}>Save</Button>
          </div>
        ) : null
      }
    >
      <div className="space-y-1 mb-4">
        {editMode && (canEditAll || isAssignee) ? (
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
            {currentAssignees.map((a) => a.username).join(", ") || "Unassigned"}
          </div>
        )}

        {!editMode && (canEditAll || isAssignee) && (
          <Button variant="secondary" onClick={() => setEditMode(true)}>
            Edit
          </Button>
        )}
      </div>

      <section className="mb-6">
        <h4 className="text-xs font-medium text-(--text-primary) mb-1">
          Description
        </h4>

        {editMode && (canEditAll || isAssignee) ? (
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

      {editMode && (canEditAll || isAssignee) && (
        <section className="mb-6 space-y-3">
          <div className="flex flex-col gap-1">
            <StatusSelector status={status} onChange={setStatus} />
          </div>

          <div className="flex flex-col gap-1">
            <AssigneeSelector
              initialAssignees={initialAssignees.map((a) => a.userId)}
              members={members}
              onChange={handleAssigneeEdit}
              isDisabled={!canEditAll}
            />
          </div>
        </section>
      )}

      <section>
        <h4 className="text-xs font-medium text-(--text-primary) mb-2">
          Comments
        </h4>

        <div className="space-y-3 mb-3">
          {!editMode && (
            <div className="flex flex-col gap-2 mb-4">
              <textarea
                placeholder="Add a comment…"
                value={comment}
                onChange={(e) => setComment(e.target.value)}
                rows={2}
                className="w-full rounded-xs bg-(--bg-surface) px-2 py-1 text-sm border border-(--border-default) outline-none resize-none focus:border-(--primary)"
              />
              <Button
                onClick={handleAddComment}
                className="bg-(--primary) text-white px-4 py-1 rounded-xs self-end"
              >
                Send
              </Button>
            </div>
          )}

          {comments.length === 0 ? (
            <p className="text-sm text-(--text-secondary)">No comments yet.</p>
          ) : (
            comments.map((comment, index) => (
              <Comment
                key={index}
                author={comment.user.username}
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
    <div className="space-y-0.5">
      <div className="text-xs text-(--text-muted)">
        <span className="font-medium text-(--text-primary)">{author}</span> ·{" "}
        {timeAgo}
      </div>
      <p className="text-sm text-(--text-secondary)">{children}</p>
    </div>
  );
}
