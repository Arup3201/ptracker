import { useState } from "react";
import { TopBar } from "../components/topbar.tsx";
import { Button } from "../components/button.tsx";
import { CreateProjectModal } from "../components/create-project.tsx";

export function Dashboard() {
  const [showModal, setShowModal] = useState(false);

  return (
    <>
      <TopBar
        title="Dashboard"
        actions={
          <Button onClick={() => setShowModal(true)}>New Project</Button>
        }
      />

      <div className="flex-1 overflow-y-auto p-4 space-y-6">
        <p className="text-center text-xl text-gray-500">
          Currently dashboard not avaiable.
        </p>
      </div>

      <CreateProjectModal
        open={showModal}
        onClose={() => setShowModal(false)}
      />
    </>
  );
}
