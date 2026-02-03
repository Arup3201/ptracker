import clsx from "clsx";

const TASK_STATUS = [
  {
    value: "Unassigned",
    label: "Unassigned",
  },
  {
    value: "Ongoing",
    label: "Ongoing",
  },
  {
    value: "Completed",
    label: "Completed",
  },
  {
    value: "Abandoned",
    label: "Abandoned",
  },
];

interface StatusSelectorProps {
  status: string;
  onChange: (v: string) => void;
}

const StatusSelector: React.FC<StatusSelectorProps> = ({
  status,
  onChange,
}) => {
  return (
    <>
      <label className="text-[12px] font-medium text-(--text-primary)">
        Status
      </label>

      <select
        value={status}
        onChange={(e) => onChange(e.target.value)}
        className={clsx(
          "h-8 rounded-xs bg-(--bg-surface) px-2 text-sm text-(--text-primary)",
          "border border-(--border-default) outline-none",
          "focus:border-(--primary)",
        )}
      >
        {TASK_STATUS.map((status) => (
          <option value={status.value}>{status.label}</option>
        ))}
      </select>
    </>
  );
};

export { StatusSelector };
