CREATE TABLE assignees (
	project_id UUID NOT NULL REFERENCES projects(id),
    task_id UUID NOT NULL REFERENCES tasks(id),
	user_id UUID NOT NULL REFERENCES users(id),
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Every project task assignee should be unique
CREATE UNIQUE INDEX ux_project_task_assignee
  ON assignees(project_id, task_id, user_id);