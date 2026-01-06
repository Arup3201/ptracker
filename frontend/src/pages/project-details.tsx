import { useState } from "react";

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

type TaskStatus = "Unassigned" | "Ongoing" | "Completed" | "Abandoned";

type Task = {
  id: string;
  title: string;
  assignee?: string;
  status: TaskStatus;
  updatedAt: string;
};

type Member = {
  id: string;
  name: string;
  role: "Owner" | "Member";
  joinedAt: string;
};

type JoinRequest = {
  id: string;
  name: string;
  note: string;
};

export default function ProjectDetailsPage() {
  const [activeTab, setActiveTab] = useState<"tasks" | "members" | "requests">(
    "tasks"
  );

  const [tasks, _] = useState<Task[]>([]);
  const [members, __] = useState<Member[]>([]);
  const [requests, ___] = useState<JoinRequest[]>([]);

  return (
    <>
      <TopBar
        title="Projects / Project Details"
        actions={
          <div className="flex gap-1">
            <Button variant="secondary">Edit Project</Button>
            <Button>Add Task</Button>
          </div>
        }
      />

      <div className="flex-1 p-4 space-y-4">
        <div className="space-y-3">
          <h1 className="max-w-3xl truncate text-lg font-semibold text-(--text-primary)">
            Collaborative Project Tracker
          </h1>

          <p className="max-w-2xl text-sm text-(--text-secondary)">
            A platform to manage collaborative software projects with clear
            ownership and task visibility.
          </p>

          <p className="text-xs text-(--text-muted)">
            Skills: PostgreSQL, React, Distributed Systems
          </p>

          <div className="flex gap-4 text-sm text-(--text-secondary)">
            <span>
              Unassigned: <strong className="text-(--text-primary)">3</strong>
            </span>
            <span>
              Ongoing: <strong className="text-(--text-primary)">5</strong>
            </span>
            <span>
              Completed: <strong className="text-(--text-primary)">12</strong>
            </span>
            <span>
              Abandoned: <strong className="text-(--text-primary)">1</strong>
            </span>
          </div>

          <div className="text-sm text-(--text-muted)">
            Members:{" "}
            <span className="font-medium text-(--text-primary)">6</span>
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
