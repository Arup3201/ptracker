CREATE TABLE join_requests (
	project_id UUID NOT NULL REFERENCES projects(id),
	user_id UUID NOT NULL REFERENCES users(id),
	status request_status NOT NULL DEFAULT 'Pending',
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX ux_project_join_requests ON join_requests(project_id, user_id);