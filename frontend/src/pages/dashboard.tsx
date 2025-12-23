import { Sidebar } from "../components/sidebar.tsx";
import { TopBar } from "../components/topbar.tsx";
import { Card } from "../components/card.tsx";
import {
  Table,
  TableHeader,
  TableBody,
  TableRow,
  TableHead,
  TableCell,
} from "../components/table.tsx";
import { Button } from "../components/button.tsx";

export function Dashboard() {
  return (
    <>
      <Sidebar />

      <main className="flex flex-1 flex-col">
        <TopBar title="Dashboard" actions={<Button>New Project</Button>} />

        <div className="flex-1 overflow-y-auto p-4 space-y-6">
          <section>
            <SectionHeader title="Active Projects" />

            <div className="grid grid-cols-1 gap-3">
              <ProjectRow
                name="Payments Revamp"
                role="Owner"
                openTasks={5}
                totalTasks={12}
              />
              <ProjectRow
                name="Auth Refactor"
                role="Member"
                openTasks={2}
                totalTasks={6}
              />
            </div>
          </section>

          <section>
            <SectionHeader title="My Tasks" />

            <Table>
              <TableHeader>
                <TableHead>Task</TableHead>
                <TableHead>Project</TableHead>
                <TableHead>Status</TableHead>
                <TableHead align="right">Updated</TableHead>
              </TableHeader>

              <TableBody>
                <TableRow>
                  <TableCell>Implement refresh tokens</TableCell>
                  <TableCell muted>Auth Refactor</TableCell>
                  <TableCell>In Progress</TableCell>
                  <TableCell align="right" muted>
                    2h ago
                  </TableCell>
                </TableRow>

                <TableRow>
                  <TableCell>Fix webhook retries</TableCell>
                  <TableCell muted>Payments Revamp</TableCell>
                  <TableCell>Open</TableCell>
                  <TableCell align="right" muted>
                    1d ago
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </section>
        </div>
      </main>
    </>
  );
}

// inline helpers

function SectionHeader({ title }: { title: string }) {
  return <div className="mb-3 text-sm font-semibold">{title}</div>;
}

function ProjectRow({
  name,
  role,
  openTasks,
  totalTasks,
}: {
  name: string;
  role: "Owner" | "Member";
  openTasks: number;
  totalTasks: number;
}) {
  return (
    <Card>
      <div className="flex items-center justify-between px-2 py-1">
        <div>
          <div className="text-sm font-medium">{name}</div>
          <div className="text-xs text-(--text-muted)">{role}</div>
        </div>

        <div className="text-xs text-(--text-secondary)">
          {openTasks} / {totalTasks} open
        </div>
      </div>
    </Card>
  );
}
