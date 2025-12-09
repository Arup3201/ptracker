CREATE TABLE roles (
	project_id UUID NOT NULL REFERENCES projects(id),
	user_id UUID NOT NULL REFERENCES users(user_id),
	role user_role NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);
