CREATE TABLE roles (
	project_id UUID NOT NULL REFERENCES projects(id),
	user_id UUID NOT NULL REFERENCES users(id),
	role user_role NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Every user must have one role in the project
CREATE UNIQUE INDEX ux_project_roles
	ON roles(project_id, user_id);