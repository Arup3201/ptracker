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
      <label className="text-sm font-medium text-text-primary">Status</label>

      <select
        value={status}
        onChange={(e) => onChange(e.target.value)}
        className={clsx(
          "h-9 w-full rounded-md bg-bg-elevated px-3 text-sm text-text-primary",
          "border border-border outline-none transition duration-fast",
          "focus:border-primary focus:shadow-focus-primary",
          "cursor-pointer",
        )}
      >
        {TASK_STATUS.map((status) => (
          <option
            key={status.value}
            value={status.value}
            className="bg-bg-elevated"
          >
            {status.label}
          </option>
        ))}
      </select>
    </>
  );
};

export { StatusSelector };
