import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Textarea } from "@/components/ui/textarea";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Calendar, Plus, CheckCircle, Loader2 } from "lucide-react";

import { HttpGet, HttpPost } from "@/utils/http";

import type { Project, NewProjectData } from "@/types/project";

const Dashboard: React.FC = () => {
  const navigate = useNavigate();

  const [isLoading, setIsLoading] = useState(false);
  const [projects, setProjects] = useState<Project[]>([]);
  const getProjects = async () => {
    setIsLoading(true);
    try {
      const data = await HttpGet("/projects/");
      setProjects(() =>
        data.projects.map((project: Project) => ({
          id: project.id,
          name: project.name,
          description: project.description,
          deadline: project.deadline,
          code: project.code,
        }))
      );
    } catch (err) {
      console.error(`getProjects failed: ${err}`);
    } finally {
      setIsLoading(false);
    }
  };
  useEffect(() => {
    getProjects();
  }, []);

  const [joinCode, setJoinCode] = useState("");
  const [joinError, setJoinError] = useState("");
  const [joinSuccess, setJoinSuccess] = useState("");
  const [isJoining, setIsJoining] = useState(false);

  const [newProject, setNewProject] = useState<NewProjectData>({
    name: "",
    description: "",
    deadline: "",
  });
  const [createError, setCreateError] = useState("");
  const [isCreating, setIsCreating] = useState(false);
  const [isDialogOpen, setIsDialogOpen] = useState(false);

  const handleJoinProject = async () => {
    setJoinError("");
    setJoinSuccess("");

    if (!joinCode.trim()) {
      setJoinError("Please enter a project code");
      return;
    }

    setIsJoining(true);

    try {
      const data = await HttpPost(`/projects/join/code/${joinCode}`, {});

      console.log("Joining project with code:", joinCode);
      setJoinSuccess(`Successfully joined the project ${data.project.name}!`);
      setJoinCode("");

      await getProjects();
    } catch (err) {
      setJoinError("Invalid project code. Please check and try again.");
    } finally {
      setIsJoining(false);
    }
  };

  const handleCreateProject = async () => {
    setCreateError("");

    if (
      !newProject.name.trim() ||
      !newProject.description.trim() ||
      !newProject.deadline
    ) {
      setCreateError("All fields are required");
      return;
    }

    // Check if deadline is in the future
    const selectedDate = new Date(newProject.deadline);
    const today = new Date();
    today.setHours(0, 0, 0, 0);

    if (selectedDate < today) {
      setCreateError("Deadline must be in the future");
      return;
    }

    setIsCreating(true);

    try {
      const data = await HttpPost("/projects/", {
        name: newProject.name,
        description: newProject.description,
        deadline: newProject.deadline,
      });

      console.log("Creating project:", newProject);

      setProjects((prev) => [
        ...prev,
        {
          id: data.project.id,
          name: data.project.name,
          description: data.project.description,
          deadline: data.project.deadline,
          code: data.project.code,
        },
      ]);

      // Reset form and close dialog
      setNewProject({ name: "", description: "", deadline: "" });
      setIsDialogOpen(false);
    } catch (err) {
      setCreateError("Failed to create project. Please try again.");
    } finally {
      setIsCreating(false);
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  };

  return (
    <div className="bg-stone-50 p-6 min-h-screen">
      <div className="mx-auto max-w-4xl">
        <div className="mb-8">
          <h1 className="mb-2 font-bold text-stone-900 text-3xl">
            Projects Dashboard
          </h1>
          <p className="text-stone-600">
            Manage your projects and collaborate with your team
          </p>
        </div>

        {/* Join Project Section */}
        <Card className="mb-8">
          <CardHeader>
            <CardTitle className="text-lg">Join a Project</CardTitle>
            <CardDescription>
              Enter a project code to join an existing project
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex gap-3">
              <div className="flex-1">
                <Input
                  placeholder="Enter project code"
                  value={joinCode}
                  onChange={(e) => setJoinCode(e.target.value)}
                  disabled={isJoining}
                />
              </div>
              <Button
                onClick={handleJoinProject}
                disabled={isJoining || !joinCode.trim()}
              >
                {isJoining ? "Joining..." : "Join"}
              </Button>
            </div>
            {joinError && (
              <Alert variant="destructive" className="mt-3">
                <AlertDescription>{joinError}</AlertDescription>
              </Alert>
            )}
            {joinSuccess && (
              <Alert className="bg-green-50 mt-3 border-green-200">
                <CheckCircle className="w-4 h-4 text-green-600" />
                <AlertDescription className="text-green-800">
                  {joinSuccess}
                </AlertDescription>
              </Alert>
            )}
          </CardContent>
        </Card>

        {/* Projects List */}
        <div className="space-y-4">
          <h2 className="font-semibold text-stone-900 text-xl">
            Your Projects
          </h2>

          <div className="gap-2 grid grid-cols-2">
            {projects.length === 0 ? (
              isLoading ? (
                <Loader2 className="mx-auto mt-2 animate-spin" size={24} />
              ) : (
                <Card>
                  <CardContent className="py-12 text-center">
                    <p className="mb-4 text-stone-500">No projects found</p>
                    <p className="text-stone-400 text-sm">
                      Create a new project or join an existing one to get
                      started
                    </p>
                  </CardContent>
                </Card>
              )
            ) : (
              projects.map((project) => (
                <Card key={project.id}>
                  <CardHeader>
                    <div className="flex justify-between items-start">
                      <div className="flex-1">
                        <CardTitle className="flex justify-between text-stone-900 text-lg">
                          <h2
                            onClick={() => navigate(`/projects/${project.id}`)}
                            className="inline-block hover:text-blue-500 hover:underline cursor-pointer"
                          >
                            {project.name}
                          </h2>
                          <span className="inline-block px-2 py-1 border border-blue-200 rounded-sm text-blue-500 text-sm">{project.code}</span>
                        </CardTitle>
                        <CardDescription className="mt-1">
                          {project.description}
                        </CardDescription>
                      </div>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <div className="flex justify-between items-center text-stone-600 text-sm">
                      <div className="flex items-center gap-4">
                        <div className="flex items-center gap-1">
                          <Calendar className="w-4 h-4" />
                          <span>Due {formatDate(project.deadline)}</span>
                        </div>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              ))
            )}
          </div>

          {/* Create New Project Button */}
          <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
            <DialogTrigger asChild>
              <Card className="hover:shadow-md border-2 border-stone-300 border-dashed transition-shadow cursor-pointer">
                <CardContent className="flex justify-center items-center py-4">
                  <div className="text-center">
                    <Plus className="mx-auto mb-2 w-8 h-8 text-stone-400" />
                    <p className="font-medium text-stone-600">
                      Create New Project
                    </p>
                    <p className="text-stone-400 text-sm">
                      Start a new collaborative project
                    </p>
                  </div>
                </CardContent>
              </Card>
            </DialogTrigger>

            <DialogContent className="sm:max-w-md">
              <DialogHeader>
                <DialogTitle>Create New Project</DialogTitle>
                <DialogDescription>
                  Fill in the details below to create a new project
                </DialogDescription>
              </DialogHeader>

              <div className="space-y-4">
                {createError && (
                  <Alert variant="destructive">
                    <AlertDescription>{createError}</AlertDescription>
                  </Alert>
                )}

                <div className="space-y-2">
                  <Label htmlFor="projectName">Project Name</Label>
                  <Input
                    id="projectName"
                    placeholder="Enter project name"
                    value={newProject.name}
                    onChange={(e) =>
                      setNewProject((prev) => ({
                        ...prev,
                        name: e.target.value,
                      }))
                    }
                    disabled={isCreating}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="projectDescription">Description</Label>
                  <Textarea
                    id="projectDescription"
                    placeholder="Enter project description"
                    value={newProject.description}
                    onChange={(e) =>
                      setNewProject((prev) => ({
                        ...prev,
                        description: e.target.value,
                      }))
                    }
                    disabled={isCreating}
                    rows={3}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="projectDeadline">Deadline</Label>
                  <Input
                    id="projectDeadline"
                    type="date"
                    value={newProject.deadline}
                    onChange={(e) =>
                      setNewProject((prev) => ({
                        ...prev,
                        deadline: e.target.value,
                      }))
                    }
                    disabled={isCreating}
                  />
                </div>
              </div>

              <DialogFooter className="flex gap-3">
                <Button
                  variant="outline"
                  onClick={() => setIsDialogOpen(false)}
                  disabled={isCreating}
                >
                  Cancel
                </Button>
                <Button onClick={handleCreateProject} disabled={isCreating}>
                  {isCreating ? "Creating..." : "Create Project"}
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
