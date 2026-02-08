import { useEffect, useState } from "react";
import { TopBar } from "../components/topbar.tsx";
import { Button } from "../components/button.tsx";
import { CreateProjectModal } from "../components/create-project.tsx";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../components/table.tsx";
import { renderTaskSignal } from "../utils.ts";
import {
  MapDashboardProject,
  MapDashboardTask,
  type DashboardProject,
  type DashboardProjectsResponse,
  type DashboardTask,
  type DashboardTasksResponse,
} from "../types/dashboard.ts";
import { ApiRequest } from "../api/request.ts";
import { useWebsocket } from "../hooks/ws.tsx";
import { WEBSOCKET_URL } from "../api/ws.ts";

export function Dashboard() {
  const [showModal, setShowModal] = useState(false);

  const [recentAssignedTasks, setRecentAssignedTasks] = useState<
    DashboardTask[]
  >([]);
  const [recentUnassignedTasks, setRecentUnassignedTasks] = useState<
    DashboardTask[]
  >([]);
  const [recentlyCreatedProjects, setRecentlyCreatedProjects] = useState<
    DashboardProject[]
  >([]);
  const [recentlyJoinedProjects, setRecentlyJoinedProjects] = useState<
    DashboardProject[]
  >([]);

  const { isConnected, reconnect } = useWebsocket({
    url: WEBSOCKET_URL,
    reconnect: true,
  });

  useEffect(() => {
    console.log(isConnected);

    if (!isConnected) {
      console.log("Trying to reconnect...");
      reconnect();
    }
  }, [isConnected]);

  async function getRecentAssignedTasks() {
    try {
      const data = await ApiRequest<DashboardTasksResponse>(
        "/dashboard/tasks/assigned",
        "GET",
        null,
      );
      if (data) {
        setRecentAssignedTasks(data.tasks.map(MapDashboardTask));
      }
    } catch (err) {
      console.error(err);
    }
  }

  async function getRecentUnassignedTasks() {
    try {
      const data = await ApiRequest<DashboardTasksResponse>(
        "/dashboard/tasks/unassigned",
        "GET",
        null,
      );
      if (data) {
        setRecentUnassignedTasks(data.tasks.map(MapDashboardTask));
      }
    } catch (err) {
      console.error(err);
    }
  }

  async function getRecentlyCreatedTasks() {
    try {
      const data = await ApiRequest<DashboardProjectsResponse>(
        "/dashboard/projects/created",
        "GET",
        null,
      );
      if (data) {
        setRecentlyCreatedProjects(data.projects.map(MapDashboardProject));
      }
    } catch (err) {
      console.error(err);
    }
  }

  async function getRecentlyJoinedTasks() {
    try {
      const data = await ApiRequest<DashboardProjectsResponse>(
        "/dashboard/projects/joined",
        "GET",
        null,
      );
      if (data) {
        setRecentlyJoinedProjects(data.projects.map(MapDashboardProject));
      }
    } catch (err) {
      console.error(err);
    }
  }

  useEffect(() => {
    getRecentAssignedTasks();
    getRecentUnassignedTasks();
    getRecentlyCreatedTasks();
    getRecentlyJoinedTasks();
  }, []);

  return (
    <>
      <TopBar
        title="Dashboard"
        actions={
          <Button onClick={() => setShowModal(true)}>New Project</Button>
        }
      />

      <div className="flex-1 overflow-y-auto p-4 space-y-6">
        <h2 className="font-semibold">Recent Assigned Tasks</h2>

        <Table>
          <TableHeader>
            <TableHead>Task</TableHead>
            <TableHead>Project</TableHead>
            <TableHead>Status</TableHead>
            <TableHead align="right">Updated</TableHead>
          </TableHeader>

          <TableBody>
            {recentAssignedTasks.length === 0 && (
              <tr>
                <td
                  colSpan={4}
                  className="px-3 py-6 text-center text-sm text-(--text-muted)"
                >
                  No tasks found
                </td>
              </tr>
            )}
            {recentAssignedTasks.map((task) => (
              <TableRow key={task.id}>
                <TableCell>
                  <a
                    href="#"
                    onClick={(e) => {
                      e.preventDefault();
                    }}
                  >
                    <span className="font-medium">{task.title}</span>
                  </a>
                </TableCell>
                <TableCell muted>{task.projectName}</TableCell>
                <TableCell muted>{task.status}</TableCell>
                <TableCell align="right" muted>
                  {task.updatedAt}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>

        <h2 className="font-semibold">Recent Unassigned Tasks</h2>

        <Table>
          <TableHeader>
            <TableHead>Task</TableHead>
            <TableHead>Assignee</TableHead>
            <TableHead>Status</TableHead>
            <TableHead align="right">Updated</TableHead>
          </TableHeader>

          <TableBody>
            {recentUnassignedTasks.length === 0 && (
              <tr>
                <td
                  colSpan={4}
                  className="px-3 py-6 text-center text-sm text-(--text-muted)"
                >
                  No tasks found
                </td>
              </tr>
            )}
            {recentUnassignedTasks.map((task) => (
              <TableRow key={task.id}>
                <TableCell>
                  <a
                    href="#"
                    onClick={(e) => {
                      e.preventDefault();
                    }}
                  >
                    <span className="font-medium">{task.title}</span>
                  </a>
                </TableCell>
                <TableCell muted>{task.projectName}</TableCell>
                <TableCell muted>{task.status}</TableCell>
                <TableCell align="right" muted>
                  {task.updatedAt}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>

        <h2 className="font-semibold">Recently Created Projects</h2>

        <Table>
          <TableHeader>
            <TableHead>Project</TableHead>
            <TableHead>Tasks</TableHead>
            <TableHead align="right">Updated</TableHead>
          </TableHeader>

          <TableBody>
            {recentlyCreatedProjects.length === 0 && (
              <tr>
                <td
                  colSpan={4}
                  className="px-3 py-6 text-center text-sm text-(--text-muted)"
                >
                  No recently created projects found
                </td>
              </tr>
            )}
            {recentlyCreatedProjects.map((project) => (
              <TableRow key={project.id}>
                <TableCell>
                  <a
                    href="#"
                    onClick={(e) => {
                      e.preventDefault();
                    }}
                  >
                    <span className="font-medium">{project.name}</span>
                  </a>
                </TableCell>
                <TableCell>{renderTaskSignal(project)}</TableCell>
                <TableCell align="right" muted>
                  {project.updatedAt}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>

        <h2 className="font-semibold">Recently Joined Projects</h2>

        <Table>
          <TableHeader>
            <TableHead>Project</TableHead>
            <TableHead>Tasks</TableHead>
            <TableHead align="right">Updated</TableHead>
          </TableHeader>

          <TableBody>
            {recentlyJoinedProjects.length === 0 && (
              <tr>
                <td
                  colSpan={4}
                  className="px-3 py-6 text-center text-sm text-(--text-muted)"
                >
                  No recently joined projects found
                </td>
              </tr>
            )}
            {recentlyJoinedProjects.map((project) => (
              <TableRow key={project.id}>
                <TableCell>
                  <a
                    href="#"
                    onClick={(e) => {
                      e.preventDefault();
                    }}
                  >
                    <span className="font-medium">{project.name}</span>
                  </a>
                </TableCell>
                <TableCell>{renderTaskSignal(project)}</TableCell>
                <TableCell align="right" muted>
                  {project.updatedAt}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      <CreateProjectModal
        open={showModal}
        onClose={() => setShowModal(false)}
      />
    </>
  );
}
