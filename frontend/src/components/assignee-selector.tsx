import { useMemo, useState } from "react";
import Select, { type MultiValue, type StylesConfig } from "react-select";
import type { Member } from "../types/project";

interface AssigneeSelectorProps {
  members: Member[];
  initialAssignees?: string[];
  onChange?: (_: string[]) => void;
  isDisabled?: boolean;
}

interface Option {
  value: string;
  label: string;
  email: string;
  avatar: string;
}
const AssigneeSelector = ({
  members,
  initialAssignees = [],
  onChange,
  isDisabled = false,
}: AssigneeSelectorProps) => {
  const [selectedAssignees, setSelectedAssignees] = useState(initialAssignees);

  // Convert members to react-select format
  const options = members.map((member) => ({
    value: member.id,
    label: member.username,
    email: member.email,
    avatar: member.avatarUrl,
  }));

  // Convert selected assignees to react-select format
  const selectedValues = useMemo<Option[]>(() => {
    return selectedAssignees
      .map((assigneeId) => options.find((opt) => opt.value === assigneeId))
      .filter((opt): opt is Option => opt !== undefined);
  }, [selectedAssignees, options]);

  const handleChange = (selected: MultiValue<Option>) => {
    const assigneeIds = selected ? selected.map((s) => s.value) : [];
    setSelectedAssignees(assigneeIds);
    onChange && onChange(assigneeIds);
  };

  const formatOptionLabel = (option: Option) => (
    <div className="flex items-center gap-2 p-1 rounded-xs bg-(--bg-tertiary) border border-(--border-muted)">
      <div className="flex flex-col min-w-0">
        <span className="text-sm font-medium text-(--text-primary) truncate">
          {option.label}
          <span className="text-[11px] text-(--text-secondary) truncate">
            ({option.email})
          </span>
        </span>
      </div>
    </div>
  );

  const customStyles: StylesConfig<Option, true> = {
    control: (base, state) => ({
      ...base,
      minHeight: "32px",
      backgroundColor: "var(--bg-surface)",
      borderColor: state.isFocused ? "var(--primary)" : "var(--border-default)",
      borderRadius: "2px", // rounded-xs
      boxShadow: "none",
      "&:hover": {
        borderColor: state.isFocused
          ? "var(--primary)"
          : "var(--border-default)",
      },
    }),
    valueContainer: (base) => ({
      ...base,
      padding: "2px 2px",
      gap: "4px",
    }),
    multiValue: (base) => ({
      ...base,
      backgroundColor: "var(--bg-tertiary)",
      borderRadius: "2px",
      padding: "0px",
      margin: "0px",
    }),
    multiValueLabel: (base) => ({
      ...base,
      color: "var(--text-primary)",
      fontSize: "0.875rem",
      fontWeight: "400",
      padding: "2px 6px",
      paddingLeft: "6px",
    }),
    multiValueRemove: (base) => ({
      ...base,
      color: "var(--text-secondary)",
      cursor: "pointer",
      paddingLeft: "2px",
      paddingRight: "4px",
      "&:hover": {
        backgroundColor: "var(--bg-hover)",
        color: "var(--text-primary)",
      },
    }),
    input: (base) => ({
      ...base,
      color: "var(--text-primary)",
      fontSize: "0.875rem",
      margin: "0px",
      padding: "0px",
    }),
    placeholder: (base) => ({
      ...base,
      color: "var(--text-tertiary)",
      fontSize: "0.875rem",
    }),
    menu: (base) => ({
      ...base,
      backgroundColor: "var(--bg-surface)",
      borderRadius: "4px",
      border: "1px solid var(--border-default)",
      boxShadow: "0 4px 6px -1px rgb(0 0 0 / 0.1)",
      marginTop: "4px",
    }),
    menuList: (base) => ({
      ...base,
      padding: "4px",
    }),
    option: (base, state) => ({
      ...base,
      backgroundColor: state.isSelected
        ? "var(--primary)"
        : state.isFocused
          ? "rgba(59, 130, 246, 0.08)" // Light blue on hover
          : "transparent",
      color: state.isSelected ? "white" : "var(--text-primary)",
      cursor: "pointer",
      borderRadius: "2px",
      padding: "4px",
      fontSize: "0.875rem",
      "&:active": {
        backgroundColor: state.isSelected
          ? "var(--primary)"
          : "rgba(59, 130, 246, 0.12)", // Slightly darker blue on click
      },
    }),
    indicatorSeparator: () => ({
      display: "none",
    }),
    dropdownIndicator: (base, state) => ({
      ...base,
      color: "var(--text-tertiary)",
      padding: "4px",
      "&:hover": {
        color: "var(--text-secondary)",
      },
      transform: state.selectProps.menuIsOpen ? "rotate(180deg)" : undefined,
      transition: "transform 0.2s",
    }),
    clearIndicator: (base) => ({
      ...base,
      color: "var(--text-tertiary)",
      padding: "4px",
      "&:hover": {
        color: "var(--text-secondary)",
      },
    }),
    noOptionsMessage: (base) => ({
      ...base,
      color: "var(--text-secondary)",
      fontSize: "0.875rem",
      padding: "8px",
    }),
  };

  return (
    <>
      <label className="text-[12px] font-medium text-(--text-primary)">
        Assignees
      </label>

      <Select
        isMulti
        options={options}
        value={selectedValues}
        onChange={handleChange}
        isDisabled={isDisabled}
        placeholder="Select assignees..."
        formatOptionLabel={formatOptionLabel}
        styles={customStyles}
        className="react-select-container"
        classNamePrefix="react-select"
        noOptionsMessage={() => "No members found"}
        closeMenuOnSelect={false}
      />

      {selectedAssignees.length === 0 && (
        <p className="text-[11px] text-(--text-tertiary)">
          No assignees yet. Task will be unassigned.
        </p>
      )}
    </>
  );
};
export default AssigneeSelector;
