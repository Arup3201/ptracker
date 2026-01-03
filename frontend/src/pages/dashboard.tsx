import { TopBar } from "../components/topbar.tsx";
import { Button } from "../components/button.tsx";

export function Dashboard() {
  return (
    <>
      <TopBar title="Dashboard" actions={<Button>New Project</Button>} />

      <div className="flex-1 overflow-y-auto p-4 space-y-6">
        <p className="text-center text-xl text-gray-500">
          Currently dashboard not avaiable.
        </p>
      </div>
    </>
  );
}
