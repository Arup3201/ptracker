import { useEffect, useState } from "react";
import { useParams } from "react-router";
import { TopBar } from "../components/topbar";
import { Button } from "../components/button";
import { ApiFetch } from "../utils/api";
import {
  JOIN_STATUS,
  MapExploredProjectDetails,
  type ExploredProjectDetails,
  type ExploredProjectDetailsApi,
  type JoinStatus,
} from "../types/explore";

export default function ProjectExplorePage() {
  const { id: projectId } = useParams();
  const [details, setDetails] = useState<ExploredProjectDetails>();
  const [joinStatus, setJoinStatus] = useState<JoinStatus>(
    JOIN_STATUS.NOT_REQUESTED,
  );

  async function getProjectDetails(id: string) {
    try {
      const response = await ApiFetch(`/public/projects/${id}`);
      if (response.ok) {
        const respondeData = await response.json();
        const data: ExploredProjectDetailsApi = respondeData.data;
        if (data) {
          setDetails(MapExploredProjectDetails(data));
          setJoinStatus(data.join_status);
        }
      } else {
        throw new Error("Failed to fetch public project.");
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

  const handleJoin = async () => {
    try {
      const response = await ApiFetch(`/projects/${projectId}/join-requests`, {
        method: "POST",
      });
      if (response.ok) {
        setJoinStatus(() => JOIN_STATUS.PENDING);
      } else {
        throw new Error("Failed to create a request.");
      }
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <>
      <TopBar
        title="Projects / Project Details"
        actions={
          <Button
            disabled={joinStatus !== JOIN_STATUS.NOT_REQUESTED}
            onClick={handleJoin}
          >
            {joinStatus === JOIN_STATUS.NOT_REQUESTED
              ? "Join Project"
              : joinStatus}
          </Button>
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
          </div>
        </div>
      </div>
    </>
  );
}
