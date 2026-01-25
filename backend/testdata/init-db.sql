CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  -- External Identity Provider identifier
  idp_subject TEXT NOT NULL,

  -- The provider issuing the subject
  idp_provider TEXT NOT NULL DEFAULT 'keycloak',

  -- Your application's local profile fields
  username TEXT NOT NULL,
  display_name TEXT,
  email TEXT,
  avatar_url TEXT,

  -- User lifecycle
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE,
  last_login_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Every user must have a unique identity provider subject
CREATE UNIQUE INDEX ux_users_subject_provider
  ON users(idp_provider, idp_subject);

CREATE TABLE sessions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

  -- Rotating refresh token (hashed) for safety
  refresh_token_encrypted BYTEA NOT NULL,

  -- Device + context info
  user_agent TEXT,
  ip_address INET,
  device_name TEXT,

  -- Session lifecycle
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  last_active_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  revoked_at TIMESTAMP WITH TIME ZONE,
  expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Useful indexes
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_active ON sessions(user_id) WHERE revoked_at IS NULL;

CREATE TABLE projects (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	owner UUID NOT NULL REFERENCES users(id),
	name VARCHAR(255) NOT NULL,
	description TEXT,
	skills TEXT,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TYPE task_status AS ENUM('Unassigned', 'Ongoing', 'Completed', 'Abandoned');

CREATE TABLE tasks (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	project_id UUID NOT NULL REFERENCES projects(id),
	title VARCHAR(255) NOT NULL,
	description TEXT,
	status task_status NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TYPE user_role AS ENUM('Owner', 'Member'); 

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

CREATE TABLE comments (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	project_id UUID NOT NULL REFERENCES projects(id),
    task_id UUID NOT NULL REFERENCES tasks(id),
	user_id UUID NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TYPE request_status AS ENUM('Pending', 'Accepted', 'Rejected');

CREATE TABLE join_requests (
	project_id UUID NOT NULL REFERENCES projects(id),
	user_id UUID NOT NULL REFERENCES users(id),
	status request_status NOT NULL DEFAULT 'Pending',
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX ux_project_join_requests ON join_requests(project_id, user_id);

CREATE VIEW project_summary AS 
SELECT p.id, 
COUNT(t.id) FILTER (WHERE t.status='Unassigned') as unassigned_tasks,
COUNT(t.id) FILTER (WHERE t.status='Ongoing') as ongoing_tasks, 
COUNT(t.id) FILTER (WHERE t.status='Completed') as completed_tasks, 
COUNT(t.id) FILTER (WHERE t.status='Abandoned') as abandoned_tasks
FROM projects as p
LEFT JOIN tasks as t ON p.id=t.project_id
WHERE p.deleted_at IS NULL AND t.deleted_at IS NULL
GROUP BY p.id;