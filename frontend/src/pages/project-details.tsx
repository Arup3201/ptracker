import { useEffect, useState } from "react";
import { useParams } from "react-router";

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
  MapProjectDetails,
  type ProjectDetails,
  type ProjectDetailsApi,
} from "../types/project";
import { MapTask, type Task, type TasksResponseApi } from "../types/task";
import type { Member } from "../types/member";
import type { JoinRequest } from "../types/join-request";
import { ApiRequest } from "../api/request";
import { AddTaskModal } from "../components/add-task";

export default function ProjectDetailsPage() {
  const [activeTab, setActiveTab] = useState<"tasks" | "members" | "requests">(
    "tasks"
  );

  const { id: projectId } = useParams();

  const [details, setDetails] = useState<ProjectDetails>();
  const [tasks, setTasks] = useState<Task[]>([]);
  const [members, __] = useState<Member[]>([]);
  const [requests, ___] = useState<JoinRequest[]>([]);

  const [addTask, setAddTask] = useState<boolean>(false);
  const [editProject, setEditProject] = useState<boolean>(false);

  async function getProjectDetails(id: string) {
    try {
      const data = await ApiRequest<ProjectDetailsApi>(
        `/projects/${id}`,
        "GET",
        null
      );
      if (data) {
        setDetails(MapProjectDetails(data));
      }
    } catch (err) {
      console.error(err);
    }
  }

  async function getProjectTasks(id: string) {
    try {
      const data = await ApiRequest<TasksResponseApi>(
        `/projects/${id}/tasks`,
        "GET",
        null
      );
      if (data) {
        setTasks(data.tasks.map(MapTask));
      }
    } catch (err) {
      console.error(err);
    }
  }

  useEffect(() => {
    if (projectId) {
      getProjectDetails(projectId);
      getProjectTasks(projectId);
    }
  }, [projectId]);

  return (
    <>
      <TopBar
        title="Projects / Project Details"
        actions={
          <div className="flex gap-1">
            <Button variant="secondary" onClick={() => setEditProject(true)}>
              Edit Project
            </Button>
            <Button onClick={() => setAddTask(true)}>Add Task</Button>
          </div>
        }
      />

      <div className="flex-1 p-4 space-y-4">
        <div className="space-y-3">
          <h1 className="max-w-3xl truncate text-lg font-semibold text-(--text-primary)">
            {details?.name || "-"}
          </h1>

          <p className="max-w-2xl text-sm text-(--text-secondary)">
            {details?.description || "-"}
          </p>

          <p className="text-xs text-(--text-muted)">
            Skills: {details?.skills || "-"}
          </p>

          <div className="flex gap-4 text-sm text-(--text-secondary)">
            <span>
              Unassigned:{" "}
              <strong className="text-(--text-primary)">
                {details?.unassignedTasks}
              </strong>
            </span>
            <span>
              Ongoing:{" "}
              <strong className="text-(--text-primary)">
                {details?.ongoingTasks}
              </strong>
            </span>
            <span>
              Completed:{" "}
              <strong className="text-(--text-primary)">
                {details?.completedTasks}
              </strong>
            </span>
            <span>
              Abandoned:{" "}
              <strong className="text-(--text-primary)">
                {details?.abandonedTasks}
              </strong>
            </span>
          </div>

          <div className="text-sm text-(--text-muted)">
            Members:{" "}
            <span className="font-medium text-(--text-primary)">
              {details?.memberCount}
            </span>
          </div>
        </div>

        <div className="flex gap-4 border-b border-(--border-muted)">
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

        {activeTab === "tasks" && <TasksSection tasks={tasks} />}
        {activeTab === "members" && <MembersSection members={members} />}
        {activeTab === "requests" && (
          <JoinRequestsSection requests={requests} />
        )}
      </div>
      <AddTaskModal
        projectId={projectId}
        open={addTask}
        onClose={() => setAddTask(false)}
      />
    </>
  );
}

function TasksSection({ tasks }: { tasks: Task[] }) {
  return (
    <div className="space-y-3">
      <div className="w-80">
        <Input placeholder="Search tasks" onChange={() => {}} />
      </div>

      <Table>
        <TableHeader>
          <TableHead>Task</TableHead>
          <TableHead>Assignee</TableHead>
          <TableHead>Status</TableHead>
          <TableHead align="right">Updated</TableHead>
        </TableHeader>

        <TableBody>
          {tasks.length === 0 && (
            <tr>
              <td
                colSpan={4}
                className="px-3 py-6 text-center text-sm text-(--text-muted)"
              >
                No tasks found
              </td>
            </tr>
          )}
          {tasks.map((task) => (
            <TableRow key={task.id}>
              <TableCell>
                <span className="font-medium">{task.title}</span>
              </TableCell>
              <TableCell muted>{task.assignee ?? "—"}</TableCell>
              <TableCell muted>{task.status}</TableCell>
              <TableCell align="right" muted>
                {task.updatedAt}
              </TableCell>
            </TableRow>
          ))}
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
        {members.length === 0 && (
          <tr>
            <td
              colSpan={4}
              className="px-3 py-6 text-center text-sm text-(--text-muted)"
            >
              No members found
            </td>
          </tr>
        )}
        {members.map((member) => (
          <TableRow key={member.id}>
            <TableCell>
              <span className="font-medium">{member.name}</span>
            </TableCell>
            <TableCell muted>{member.role}</TableCell>
            <TableCell align="right" muted>
              {member.joinedAt}
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}

function JoinRequestsSection({ requests }: { requests: JoinRequest[] }) {
  return (
    <div className="divide-y divide-(--border-muted) rounded-md border border-(--border-default)">
      {requests.length === 0 && (
        <div className="px-3 py-6 text-center text-sm text-(--text-muted)">
          No join requests found
        </div>
      )}
      {requests.map((req) => (
        <div key={req.id} className="flex items-start justify-between p-4">
          <div className="space-y-1">
            <div className="text-sm font-medium">{req.name}</div>
            <p className="text-sm text-(--text-secondary)">{req.note}</p>
          </div>

          <div className="flex gap-2">
            <Button variant="secondary">Reject</Button>
            <Button>Accept</Button>
          </div>
        </div>
      ))}
    </div>
  );
}
