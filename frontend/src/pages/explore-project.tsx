import { useEffect, useState } from "react";
import { useParams } from "react-router";
import { TopBar } from "../components/topbar";
import { Button } from "../components/button";
import { ApiRequest } from "../api/request";
import {
  MapExploredProjectDetails,
  type ExploredProjectDetails,
  type ExploredProjectDetailsApi,
} from "../types/explore";

export default function ProjectExplorePage() {
  const { id: projectId } = useParams();
  const [details, setDetails] = useState<ExploredProjectDetails>();

  async function getProjectDetails(id: string) {
    try {
      const data = await ApiRequest<ExploredProjectDetailsApi>(
        `/explore/projects/${id}`,
        "GET",
        null,
      );
      if (data) {
        setDetails(MapExploredProjectDetails(data));
      }
    } catch (err) {
      console.error(err);
    }
  }

  useEffect(() => {
    if (projectId) {
      getProjectDetails(projectId);
    }
  }, [projectId]);

  const handleJoin = () => {};

  return (
    <>
      <TopBar
        title="Projects / Project Details"
        actions={
          <div className="flex gap-1">
            <Button onClick={handleJoin}>Join Project</Button>
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
        </div>
      </div>
    </>
  );
}
