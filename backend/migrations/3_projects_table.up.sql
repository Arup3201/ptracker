CREATE TABLE projects (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	code VARCHAR(255) NOT NULL,
	owner UUID NOT NULL REFERENCES users(id),
	name VARCHAR(255) NOT NULL,
	description TEXT,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX ux_project_code ON projects(code);