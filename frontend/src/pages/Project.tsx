import React, { useState, useEffect } from "react";
import { useParams } from "react-router";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import { Badge } from "@/components/ui/badge";
import { Alert, AlertDescription } from "@/components/ui/alert";
import {
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from "@tanstack/react-table";
import {
  Plus,
  Edit,
  Trash2,
  UserPlus,
  User,
  Clock,
  CheckCircle2,
  Circle,
  Loader2,
} from "lucide-react";

import { HttpGet, HttpPost, HttpPut } from "@/utils/http";
import useAuth from "@/hooks/auth";

import type { ProjectRole, Project as ProjectType } from "@/types/project";
import type { Task, TeamMember, NewTaskData, TaskStatus } from "@/types/task";
import type { MemberResponse, TaskResponse } from "@/types/response";

const Project: React.FC = () => {
  const { id: project_id } = useParams();
  const { user } = useAuth();

  const [isLoading, setIsLoading] = useState(false);

  const [project, setProject] = useState<ProjectType>({} as ProjectType);
  const [tasks, setTasks] = useState<Task[]>([]);
  const [teamMembers, setTeamMembers] = useState<TeamMember[]>([]);

  const [projectRole, setProjectRole] = useState<ProjectRole>();

  const getProject = async (projectId: string) => {
    setIsLoading(true);
    try {
      const data = await HttpGet(`/projects/${projectId}`);
      setProject({
        id: data.project.id,
        name: data.project.name,
        description: data.project.description,
        deadline: data.project.deadline,
        code: data.project.code,
      });
      setProjectRole(data.project.role);
      setTasks(() =>
        data.tasks.map((task: TaskResponse) => ({
          id: task.id,
          name: task.name,
          description: task.description,
          status: task.status,
          assignee: {
            id: task.assignee,
            name: task.assignee_name,
            email: task.assignee_email,
          },
        }))
      );
    } catch (err) {
      console.error(`getProject failed: ${err}`);
    } finally {
      setIsLoading(false);
    }
  };
  const getMembers = async (projectId: string) => {
    try {
      const data = await HttpGet(`/projects/${projectId}/members`);
      setTeamMembers(() =>
        data.members.map((member: MemberResponse) => ({
          id: member.user_id,
          email: member.email,
          name: member.name,
          joinedAt: member.joined_at,
          role: member.role,
        }))
      );
    } catch (err) {
      console.error(`getMembers failed: ${err}`);
    }
  };
  useEffect(() => {
    if (project_id) {
      getProject(project_id);
      getMembers(project_id);
    }
  }, [project_id]);

  const [newTask, setNewTask] = useState<NewTaskData>({
    name: "",
    description: "",
    assignee: "",
    status: "To Do",
  });
  const [editTask, setEditTask] = useState<NewTaskData>({
    name: "",
    description: "",
    assignee: "",
    status: "To Do",
  });

  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [createError, setCreateError] = useState("");
  const [isCreating, setIsCreating] = useState(false);

  const [editingDialogOpen, setEditingDialogOpen] = useState(false);
  const [editError, setEditError] = useState("");
  const [isEditing, setIsEditing] = useState(false);

  const handleCreateTask = async () => {
    setCreateError("");

    if (
      !newTask.name.trim() ||
      !newTask.description.trim() ||
      !newTask.assignee
    ) {
      setCreateError("All fields are required");
      return;
    }

    setIsCreating(true);

    try {
      const data = await HttpPost(`/projects/${project_id}/tasks/`, {
        name: newTask.name,
        description: newTask.description,
        assignee: newTask.assignee,
        status: newTask.status,
      });

      console.log("Creating task:", newTask);

      setTasks((prev) => [
        ...prev,
        {
          id: data.task.id,
          name: data.task.name,
          description: data.task.description,
          assignee: {
            id: data.task.assignee,
            name: data.task.assignee_name,
            email: data.task.assignee_email,
          },
          status: data.task.status,
        },
      ]);

      // Reset form and close dialog
      setNewTask({
        name: "",
        description: "",
        assignee: "",
        status: "To Do",
      });
      setIsCreateDialogOpen(false);
    } catch (err) {
      setCreateError("Failed to create task. Please try again.");
    } finally {
      setIsCreating(false);
    }
  };

  const handleDeleteTask = async (taskId: string) => {
    try {
      // TODO: Replace with your actual delete task API call
      console.log("Deleting task:", taskId);
      setTasks((prev) => prev.filter((task) => task.id !== taskId));
    } catch (err) {
      console.error("Failed to delete task");
    }
  };

  const openEditTaskModal = (taskId: string) => {
    setEditingDialogOpen(true);
    const task = tasks.find((t) => t.id === taskId);
    if (!task) return;
    setEditTask({
      name: task.name,
      description: task.description,
      status: task.status,
      assignee: task.assignee.id,
    });
  };

  const handleEditTask = async () => {
    setEditError("");

    if (
      !editTask.name.trim() ||
      !editTask.description.trim() ||
      !editTask.assignee
    ) {
      setEditError("All fields are required");
      return;
    }

    setIsEditing(true);

    try {
      const data = await HttpPut(`/projects/${project_id}/tasks/`, {
        name: editTask.name,
        description: editTask.description,
      });

      console.log("Creating task:", editTask);

      // Reset form and close dialog
      setEditTask({
        name: "",
        description: "",
        assignee: "",
        status: "To Do",
      });
      setEditingDialogOpen(false);
    } catch (err) {
      setEditError("Failed to create task. Please try again.");
    } finally {
      setIsEditing(false);
    }
  }

  const handleViewTask = (taskId: string) => {
    // TODO: Implement view task details functionality
    console.log("View task:", taskId);
  };

  const handleAssignTask = (taskId: string) => {
    // TODO: Implement assign task functionality
    console.log("Assign task:", taskId);
  };

  const truncateDescription = (description: string, maxLength: number = 50) => {
    return description.length > maxLength
      ? description.substring(0, maxLength) + "..."
      : description;
  };

  const getStatusBadge = (status: TaskStatus) => {
    const statusConfig: any = {
      "To Do": {
        color: "bg-gray-100 text-gray-800",
        icon: <Circle className="w-3 h-3" />,
      },
      "In Progress": {
        color: "bg-blue-100 text-blue-800",
        icon: <Clock className="w-3 h-3" />,
      },
      Completed: {
        color: "bg-green-100 text-green-800",
        icon: <CheckCircle2 className="w-3 h-3" />,
      },
    };

    const config = statusConfig[status];

    return (
      <Badge className={`${config.color} flex items-center gap-1`}>
        {config.icon}
        {status}
      </Badge>
    );
  };

  // TanStack Table setup
  const columnHelper = createColumnHelper<Task>();

  const columns = [
    columnHelper.accessor("name", {
      header: "Task Name",
      cell: (info) => <div className="font-medium">{info.getValue()}</div>,
    }),
    columnHelper.accessor("description", {
      header: "Description",
      cell: (info) => (
        <span className="text-stone-600" title={info.getValue()}>
          {truncateDescription(info.getValue())}
        </span>
      ),
    }),
    columnHelper.accessor("assignee", {
      header: "Assignee",
      cell: (info) => (
        <div className="flex items-center gap-2">
          <User className="w-4 h-4 text-stone-400" />
          {info.getValue().name}
        </div>
      ),
    }),
    columnHelper.accessor("status", {
      header: "Status",
      cell: (info) => getStatusBadge(info.getValue()),
    }),
    columnHelper.display({
      id: "actions",
      header: "Actions",
      cell: (info) => (
        <div className="flex justify-center gap-1">
          <Button
                variant="ghost"
                onClick={() => {
                  if (
                    projectRole !== "Owner" &&
                    user?.id !== info.row.original.assignee.id
                  ) {
                    // user is not the owner or the assignee
                    return;
                  }
                  openEditTaskModal(info.row.original.id);
                }}
                disabled={
                  projectRole !== "Owner" &&
                  user?.id !== info.row.original.assignee.id
                }
                className="cursor-pointer"
              >
                <Edit className="w-4 h-4" />
              </Button>
          <Dialog open={editingDialogOpen} onOpenChange={setEditingDialogOpen}>
            <DialogContent className="sm:max-w-md">
              <DialogHeader>
                <DialogTitle>Edit Task</DialogTitle>
                <DialogDescription>
                  Fill in the details below to edit this task
                </DialogDescription>
              </DialogHeader>

              <div className="space-y-4">
                {editError && (
                  <Alert variant="destructive">
                    <AlertDescription>{editError}</AlertDescription>
                  </Alert>
                )}

                <div className="space-y-2">
                  <Label htmlFor="taskName">Task Name</Label>
                  <Input
                    id="taskName"
                    placeholder="Enter task name"
                    value={editTask.name}
                    onChange={(e) =>
                      setEditTask((prev) => ({
                        ...prev,
                        name: e.target.value,
                      }))
                    }
                    disabled={isEditing}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="taskDescription">Description</Label>
                  <Textarea
                    id="taskDescription"
                    placeholder="Enter task description"
                    value={editTask.description}
                    onChange={(e) =>
                      setEditTask((prev) => ({
                        ...prev,
                        description: e.target.value,
                      }))
                    }
                    disabled={isEditing}
                    rows={3}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="taskAssignee">Assignee</Label>
                  <Select
                    value={editTask.assignee}
                    onValueChange={(value) =>
                      setEditTask((prev) => ({ ...prev, assignee: value }))
                    }
                    disabled={isEditing}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="Select assignee" />
                    </SelectTrigger>
                    <SelectContent>
                      {teamMembers.map((member) => (
                        <SelectItem key={member.id} value={member.id}>
                          <div className="flex items-center gap-2">
                            <User className="w-4 h-4" />
                            {member.name}
                          </div>
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="taskStatus">Status</Label>
                  <Select
                    value={editTask.status}
                    onValueChange={(value: TaskStatus) =>
                      setEditTask((prev) => ({ ...prev, status: value }))
                    }
                    disabled={isEditing}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="To Do">To Do</SelectItem>
                      <SelectItem value="In Progress">In Progress</SelectItem>
                      <SelectItem value="Done">Done</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>

              <DialogFooter className="flex gap-3">
                <Button
                  variant="outline"
                  onClick={() => setIsCreateDialogOpen(false)}
                  disabled={isEditing}
                >
                  Cancel
                </Button>
                <Button
                  onClick={handleEditTask}
                  disabled={
                    isEditing ||
                    projectRole !== "Owner" ||
                    user?.id !== editTask.assignee
                  }
                >
                  {isEditing ? "Editing..." : "Edit Task"}
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
          <Button
            variant="ghost"
            onClick={() => {
              if (projectRole !== "Owner") {
                // user is not the owner
                return;
              }
              handleAssignTask(info.row.original.id);
            }}
            className="cursor-pointer"
            disabled={projectRole !== "Owner"}
          >
            <UserPlus className="w-4 h-4" />
          </Button>
          <Button
            variant="ghost"
            onClick={() => {
              if (projectRole !== "Owner") {
                // user is not the owner
                return;
              }
              handleDeleteTask(info.row.original.id);
            }}
            className="text-destructive hover:text-red-600 cursor-pointer"
            disabled={projectRole !== "Owner"}
          >
            <Trash2 className="w-4 h-4" />
          </Button>
        </div>
      ),
    }),
  ];

  const table = useReactTable({
    data: tasks,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <div className="bg-stone-50 p-6 min-h-screen">
      <div className="mx-auto max-w-7xl">
        {/* Header */}
        <div className="mb-8">
          <h1 className="mb-2 font-bold text-stone-900 text-3xl">
            {project.name}
          </h1>
          <p className="text-stone-600">Manage tasks and track progress</p>
        </div>

        {/* Tasks Table */}
        <Card>
          <CardHeader className="flex flex-row justify-between items-center">
            <CardTitle className="text-xl">Tasks</CardTitle>

            {/* Create New Task Button - moved to the right */}
            <Dialog
              open={isCreateDialogOpen}
              onOpenChange={setIsCreateDialogOpen}
            >
              <DialogTrigger asChild>
                <Button className="flex items-center gap-2">
                  <Plus className="w-4 h-4" />
                  Create New Task
                </Button>
              </DialogTrigger>

              <DialogContent className="sm:max-w-md">
                <DialogHeader>
                  <DialogTitle>Create New Task</DialogTitle>
                  <DialogDescription>
                    Fill in the details below to create a new task
                  </DialogDescription>
                </DialogHeader>

                <div className="space-y-4">
                  {createError && (
                    <Alert variant="destructive">
                      <AlertDescription>{createError}</AlertDescription>
                    </Alert>
                  )}

                  <div className="space-y-2">
                    <Label htmlFor="taskName">Task Name</Label>
                    <Input
                      id="taskName"
                      placeholder="Enter task name"
                      value={newTask.name}
                      onChange={(e) =>
                        setNewTask((prev) => ({
                          ...prev,
                          name: e.target.value,
                        }))
                      }
                      disabled={isCreating}
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="taskDescription">Description</Label>
                    <Textarea
                      id="taskDescription"
                      placeholder="Enter task description"
                      value={newTask.description}
                      onChange={(e) =>
                        setNewTask((prev) => ({
                          ...prev,
                          description: e.target.value,
                        }))
                      }
                      disabled={isCreating}
                      rows={3}
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="taskAssignee">Assignee</Label>
                    <Select
                      value={newTask.assignee}
                      onValueChange={(value) =>
                        setNewTask((prev) => ({ ...prev, assignee: value }))
                      }
                      disabled={isCreating}
                    >
                      <SelectTrigger>
                        <SelectValue placeholder="Select assignee" />
                      </SelectTrigger>
                      <SelectContent>
                        {teamMembers.map((member) => (
                          <SelectItem key={member.id} value={member.id}>
                            <div className="flex items-center gap-2">
                              <User className="w-4 h-4" />
                              {member.name}
                            </div>
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="taskStatus">Status</Label>
                    <Select
                      value={newTask.status}
                      onValueChange={(value: TaskStatus) =>
                        setNewTask((prev) => ({ ...prev, status: value }))
                      }
                      disabled={isCreating}
                    >
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="To Do">To Do</SelectItem>
                        <SelectItem value="In Progress">In Progress</SelectItem>
                        <SelectItem value="Done">Done</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>

                <DialogFooter className="flex gap-3">
                  <Button
                    variant="outline"
                    onClick={() => setIsCreateDialogOpen(false)}
                    disabled={isCreating}
                  >
                    Cancel
                  </Button>
                  <Button
                    onClick={handleCreateTask}
                    disabled={isCreating || projectRole !== "Owner"}
                  >
                    {isCreating ? "Creating..." : "Create Task"}
                  </Button>
                </DialogFooter>
              </DialogContent>
            </Dialog>
          </CardHeader>

          <CardContent>
            {tasks.length === 0 ? (
              isLoading ? (
                <Loader2 className="mx-auto mt-2 animate-spin" size={24} />
              ) : (
                <div className="py-12 text-center">
                  <p className="mb-4 text-stone-500">No tasks found</p>
                  <p className="text-stone-400 text-sm">
                    Create a new task to get started
                  </p>
                </div>
              )
            ) : (
              <div className="border rounded-md">
                <table className="w-full text-sm caption-bottom">
                  <thead className="[&_tr]:border-b">
                    {table.getHeaderGroups().map((headerGroup) => (
                      <tr
                        key={headerGroup.id}
                        className="[&_th:last-child]:border-r-0"
                      >
                        {headerGroup.headers.map((header) => (
                          <th
                            key={header.id}
                            className="px-4 border-r-1 h-12 font-medium text-stone-500 text-center align-middle b"
                          >
                            {header.isPlaceholder
                              ? null
                              : flexRender(
                                  header.column.columnDef.header,
                                  header.getContext()
                                )}
                          </th>
                        ))}
                      </tr>
                    ))}
                  </thead>
                  <tbody className="[&_tr:last-child]:border-0">
                    {table.getRowModel().rows.map((row) => (
                      <tr
                        key={row.id}
                        className="hover:bg-stone-50 [&_td:last-child]:border-r-0 border-b transition-colors"
                      >
                        {row.getVisibleCells().map((cell) => (
                          <td
                            key={cell.id}
                            className="p-4 border-r-1 align-middle"
                          >
                            {flexRender(
                              cell.column.columnDef.cell,
                              cell.getContext()
                            )}
                          </td>
                        ))}
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default Project;
