import { useEffect, useState } from "react";
import { useParams, useSearchParams } from "react-router";

import { TopBar } from "../components/topbar";
import {
  Table,
  TableHeader,
  TableBody,
  TableRow,
  TableHead,
  TableCell,
} from "../components/table";
import { Button } from "../components/button";
import { Input } from "../components/input";
import { Tab } from "../components/tab";
import {
  MapJoinRequest,
  MapMember,
  MapProjectDetails,
  ROLES,
  type JoinRequest,
  type JoinRequestsResponseApi,
  type Member,
  type MembersResponse,
  type ProjectDetails,
  type ProjectDetailsApi,
} from "../types/project";
import { MapTask, type Task, type TasksResponseApi } from "../types/task";
import { ApiFetch } from "../utils/api";
import { AddTaskModal } from "../components/add-task";
import { TaskDrawer } from "./task-drawer";
import { JOIN_STATUS } from "../types/explore";
import { renderLocalTime } from "../utils";

export default function ProjectDetailsPage() {
  const [activeTab, setActiveTab] = useState<"tasks" | "members" | "requests">(
    "tasks",
  );

  const { id: projectId } = useParams();
  const [searchParams, setSearchParams] = useSearchParams();

  const [details, setDetails] = useState<ProjectDetails>();
  const [tasks, setTasks] = useState<Task[]>([]);
  const [members, setMembers] = useState<Member[]>([]);
  const [joinRequests, setJoinRequests] = useState<JoinRequest[]>([]);

  const [addTask, setAddTask] = useState<boolean>(false);
  const [_editProject, setEditProject] = useState<boolean>(false);

  async function getProjectDetails(id: string) {
    try {
      const response = await ApiFetch(`/projects/${id}`);
      if (response.ok) {
        const respondeData = await response.json();
        const data: ProjectDetailsApi = respondeData.data;
        if (data) {
          setDetails(MapProjectDetails(data));
        }
      } else {
        throw new Error("Failed to get project details.");
      }
    } catch (err) {
      console.error(err);
    }
  }

  async function getProjectTasks(id: string) {
    try {
      const response = await ApiFetch(`/projects/${id}/tasks`);
      if (response.ok) {
        const respondeData = await response.json();
        const data: TasksResponseApi = respondeData.data;
        if (data) {
          setTasks(data.tasks.map(MapTask));
        }
      } else {
        throw new Error("Failed to get project tasks.");
      }
    } catch (err) {
      console.error(err);
    }
  }

  async function getProjectMembers(id: string) {
    try {
      const response = await ApiFetch(`/projects/${id}/members`);
      if (response.ok) {
        const respondeData = await response.json();
        const data: MembersResponse = respondeData.data;
        if (data) {
          setMembers(data.members.map(MapMember));
        }
      } else {
        throw new Error("Failed to get project members.");
      }
    } catch (err) {
      console.error(err);
    }
  }

  async function getJoinRequests(id: string) {
    try {
      const response = await ApiFetch(`/projects/${id}/join-requests`);
      if (response.ok) {
        const respondeData = await response.json();
        const data: JoinRequestsResponseApi = respondeData.data;
        if (data) {
          setJoinRequests(data.join_requests.map(MapJoinRequest));
        }
      } else {
        throw new Error("Failed to get project join requests.");
      }
    } catch (err) {
      console.error(err);
    }
  }

  useEffect(() => {
    if (projectId) {
      getProjectDetails(projectId);
      getProjectTasks(projectId);
      getProjectMembers(projectId);
      getJoinRequests(projectId);
    }
  }, [projectId]);

  const taskId = searchParams.get("task"); // string | null
  const closeTaskDrawer = () => {
    searchParams.delete("task");
    setSearchParams(searchParams);
  };
  const openTask = (taskId: string) => {
    searchParams.set("task", taskId);
    setSearchParams(searchParams);
  };

  return (
    <>
      <TopBar
        title="Projects / Project Details"
        actions={
          <div className="flex gap-2">
            {details?.role === ROLES.OWNER && (
              <Button variant="secondary" onClick={() => setEditProject(true)}>
                Edit Project
              </Button>
            )}
            {details?.role === ROLES.OWNER && (
              <Button onClick={() => setAddTask(true)}>Add Task</Button>
            )}
          </div>
        }
      />

      <div className="flex-1 overflow-y-auto p-6 space-y-6">
        {/* Project meta */}
        <div className="space-y-3">
          <h1 className="max-w-3xl truncate text-2xl font-semibold text-text-primary tracking-tight">
            {details?.name || "—"}
          </h1>
          <p className="max-w-2xl text-sm text-text-secondary leading-relaxed">
            {details?.description || "—"}
          </p>
          {details?.skills && (
            <p className="text-xs text-text-muted">Skills: {details.skills}</p>
          )}

          {/* Stats row */}
          <div className="flex items-center gap-1 flex-wrap">
            {[
              { label: "Unassigned", value: details?.unassignedTasks },
              { label: "Ongoing", value: details?.ongoingTasks },
              { label: "Completed", value: details?.completedTasks },
              { label: "Abandoned", value: details?.abandonedTasks },
            ].map(({ label, value }) => (
              <div
                key={label}
                className="flex items-center gap-1.5 rounded-md border border-border bg-bg-elevated px-3 py-1.5"
              >
                <span className="text-xs text-text-muted">{label}</span>
                <span className="text-sm font-semibold text-text-primary">
                  {value ?? 0}
                </span>
              </div>
            ))}
            <div className="flex items-center gap-1.5 rounded-md border border-border bg-bg-elevated px-3 py-1.5">
              <span className="text-xs text-text-muted">Members</span>
              <span className="text-sm font-semibold text-text-primary">
                {details?.membersCount ?? 0}
              </span>
            </div>
          </div>
        </div>

        {/* Tabs */}
        <div className="flex gap-5 border-b border-border-muted">
          <Tab
            label="Tasks"
            active={activeTab === "tasks"}
            onClick={() => setActiveTab("tasks")}
          />
          <Tab
            label="Members"
            active={activeTab === "members"}
            onClick={() => setActiveTab("members")}
          />
          <Tab
            label="Join Requests"
            active={activeTab === "requests"}
            onClick={() => setActiveTab("requests")}
          />
        </div>

        {/* Tab panels */}
        {activeTab === "tasks" && (
          <TasksSection onOpenTask={openTask} tasks={tasks} />
        )}
        {activeTab === "members" && <MembersSection members={members} />}
        {activeTab === "requests" && (
          <JoinRequestsSection requests={joinRequests} />
        )}
      </div>

      <AddTaskModal
        projectId={projectId}
        members={members}
        open={addTask}
        onClose={() => setAddTask(false)}
      />
      <TaskDrawer
        open={Boolean(taskId)}
        taskId={taskId}
        projectId={projectId}
        members={members}
        onClose={closeTaskDrawer}
        role={details?.role || ROLES.MEMBER}
      />
    </>
  );
}

function TasksSection({
  tasks,
  onOpenTask,
}: {
  tasks: Task[];
  onOpenTask: (taskId: string) => void;
}) {
  return (
    <div className="space-y-3">
      <div className="w-72">
        <Input placeholder="Search tasks..." onChange={() => {}} />
      </div>

      <Table>
        <TableHeader>
          <TableHead>Task</TableHead>
          <TableHead>Assignee</TableHead>
          <TableHead>Status</TableHead>
          <TableHead align="right">Updated</TableHead>
        </TableHeader>
        <TableBody>
          {tasks.length === 0 ? (
            <tr>
              <td
                colSpan={4}
                className="px-4 py-8 text-center text-sm text-text-muted"
              >
                No tasks found
              </td>
            </tr>
          ) : (
            tasks.map((task) => (
              <TableRow key={task.id}>
                <TableCell>
                  <a
                    href="#"
                    onClick={(e) => {
                      e.preventDefault();
                      onOpenTask(task.id);
                    }}
                    className="font-medium text-text-primary hover:text-primary transition duration-fast"
                  >
                    {task.title}
                  </a>
                </TableCell>
                <TableCell muted>
                  {task.assignees.map((a) => a.avatar.username).join(", ") ||
                    "—"}
                </TableCell>
                <TableCell muted>{task.status}</TableCell>
                <TableCell align="right" muted>
                  {task.updatedAt && renderLocalTime(task.updatedAt)}
                </TableCell>
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </div>
  );
}

function MembersSection({ members }: { members: Member[] }) {
  return (
    <Table>
      <TableHeader>
        <TableHead>Name</TableHead>
        <TableHead>Role</TableHead>
        <TableHead align="right">Joined</TableHead>
      </TableHeader>
      <TableBody>
        {members.length === 0 ? (
          <tr>
            <td
              colSpan={2}
              className="px-4 py-8 text-center text-sm text-text-muted"
            >
              No members found
            </td>
          </tr>
        ) : (
          members.map((member) => (
            <TableRow key={member.avatar.userId}>
              <TableCell>
                <div className="flex items-center gap-2.5">
                  <div className="h-7 w-7 rounded-full bg-bg-elevated border border-border flex items-center justify-center text-xs font-semibold text-primary shrink-0">
                    {member.avatar.displayName
                      ? member.avatar.displayName?.charAt(0).toUpperCase()
                      : member.avatar.username}
                  </div>
                  <span className="font-medium text-text-primary">
                    {member.avatar.displayName || member.avatar.username}
                  </span>
                </div>
              </TableCell>
              <TableCell>{member.role}</TableCell>
              <TableCell align="right" muted>
                {renderLocalTime(member.createdAt)}
              </TableCell>
            </TableRow>
          ))
        )}
      </TableBody>
    </Table>
  );
}

function JoinRequestsSection({ requests }: { requests: JoinRequest[] }) {
  const handleUpdate = async (
    projectId: string,
    userId: string,
    joinStatus: string,
  ) => {
    try {
      const response = await ApiFetch(`/projects/${projectId}/join-requests`, {
        method: "PUT",
        body: JSON.stringify({ user_id: userId, join_status: joinStatus }),
      });
      if (!response.ok) {
        throw new Error("Failed to respond to the join request.");
      }
    } catch (err) {
      console.error(err);
    }
  };

  const pendingRequests = requests.filter(
    (r) => r.status === JOIN_STATUS.PENDING,
  );

  return (
    <div className="rounded-lg border border-border bg-bg-surface overflow-hidden">
      {pendingRequests.length === 0 ? (
        <div className="px-4 py-8 text-center text-sm text-text-muted">
          No pending join requests
        </div>
      ) : (
        <div className="divide-y divide-border-muted">
          {pendingRequests.map((req) => (
            <div
              key={req.projectId + req.avatar.userId}
              className="flex items-center justify-between gap-4 px-4 py-3 hover:bg-bg-elevated transition duration-fast"
            >
              <div className="flex items-center gap-2.5 min-w-0">
                <div className="h-7 w-7 rounded-full bg-bg-elevated border border-border flex items-center justify-center text-xs font-semibold text-primary shrink-0">
                  {req.avatar.displayName
                    ? req.avatar.displayName?.charAt(0).toUpperCase()
                    : req.avatar.username.charAt(0).toUpperCase()}
                </div>
                <div className="min-w-0">
                  <p className="text-sm font-medium text-text-primary truncate">
                    {req.avatar.displayName || req.avatar.username}
                  </p>
                  <p className="text-xs text-text-muted truncate">
                    {req.avatar.email}
                  </p>
                </div>
              </div>

              <div className="flex gap-2 shrink-0">
                <Button
                  variant="secondary"
                  onClick={() =>
                    handleUpdate(req.projectId, req.avatar.userId, "Rejected")
                  }
                >
                  Reject
                </Button>
                <Button
                  onClick={() =>
                    handleUpdate(req.projectId, req.avatar.userId, "Accepted")
                  }
                >
                  Accept
                </Button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
