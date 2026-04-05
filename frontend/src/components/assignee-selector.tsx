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

  const options = members.map((member) => ({
    value: member.userId,
    label: member.username,
    email: member.email,
    avatar: member.avatarUrl,
  }));

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
    <div className="flex items-center gap-2">
      <div className="flex flex-col min-w-0">
        <span className="text-sm font-medium text-text-primary truncate">
          {option.label}
          <span className="text-xs text-text-muted ml-1">({option.email})</span>
        </span>
      </div>
    </div>
  );

  const customStyles: StylesConfig<Option, true> = {
    control: (base, state) => ({
      ...base,
      minHeight: "36px",
      backgroundColor: "var(--bg-elevated)",
      borderColor: state.isFocused ? "var(--primary)" : "var(--border-default)",
      borderRadius: "8px",
      boxShadow: state.isFocused ? "var(--tw-shadow-focus-primary)" : "none",
      "&:hover": {
        borderColor: state.isFocused
          ? "var(--primary)"
          : "var(--border-strong)",
      },
      transition: "border-color 150ms, box-shadow 150ms",
    }),
    valueContainer: (base) => ({
      ...base,
      padding: "2px 8px",
      gap: "4px",
    }),
    multiValue: (base) => ({
      ...base,
      backgroundColor: "var(--bg-overlay)",
      borderRadius: "6px",
      border: "1px solid var(--border-default)",
      margin: "0px",
      padding: "0px",
    }),
    multiValueLabel: (base) => ({
      ...base,
      color: "var(--text-secondary)",
      fontSize: "12px",
      fontWeight: "500",
      padding: "2px 6px",
      paddingLeft: "8px",
    }),
    multiValueRemove: (base) => ({
      ...base,
      color: "var(--text-muted)",
      cursor: "pointer",
      borderRadius: "0 6px 6px 0",
      paddingLeft: "2px",
      paddingRight: "6px",
      "&:hover": {
        backgroundColor: "var(--danger-muted)",
        color: "var(--danger)",
      },
    }),
    input: (base) => ({
      ...base,
      color: "var(--text-primary)",
      fontSize: "13px",
      margin: "0px",
      padding: "0px",
    }),
    placeholder: (base) => ({
      ...base,
      color: "var(--text-muted)",
      fontSize: "13px",
    }),
    menu: (base) => ({
      ...base,
      backgroundColor: "var(--bg-elevated)",
      borderRadius: "10px",
      border: "1px solid var(--border-default)",
      boxShadow:
        "0 8px 24px rgba(0,0,0,0.7), 0 0 0 0.5px rgba(255,255,255,0.04)",
      marginTop: "4px",
      overflow: "hidden",
    }),
    menuList: (base) => ({
      ...base,
      padding: "4px",
    }),
    option: (base, state) => ({
      ...base,
      backgroundColor: state.isSelected
        ? "var(--primary-muted)"
        : state.isFocused
          ? "var(--bg-overlay)"
          : "transparent",
      color: state.isSelected ? "var(--primary)" : "var(--text-primary)",
      cursor: "pointer",
      borderRadius: "6px",
      padding: "6px 8px",
      fontSize: "13px",
      transition: "background-color 150ms",
      "&:active": {
        backgroundColor: "var(--primary-muted)",
      },
    }),
    indicatorSeparator: () => ({
      display: "none",
    }),
    dropdownIndicator: (base, state) => ({
      ...base,
      color: "var(--text-muted)",
      padding: "0 8px",
      "&:hover": {
        color: "var(--text-secondary)",
      },
      transform: state.selectProps.menuIsOpen ? "rotate(180deg)" : undefined,
      transition: "transform 150ms, color 150ms",
    }),
    clearIndicator: (base) => ({
      ...base,
      color: "var(--text-muted)",
      padding: "0 4px",
      "&:hover": {
        color: "var(--text-secondary)",
      },
    }),
    noOptionsMessage: (base) => ({
      ...base,
      color: "var(--text-muted)",
      fontSize: "13px",
      padding: "8px",
    }),
  };

  return (
    <>
      <label className="text-sm font-medium text-text-primary">Assignees</label>

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
        <p className="text-xs text-text-muted">
          No assignees yet. Task will be unassigned.
        </p>
      )}
    </>
  );
};
export default AssigneeSelector;
