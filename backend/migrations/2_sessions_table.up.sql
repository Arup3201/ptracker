CREATE TABLE sessions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

  -- Rotating refresh token (hashed) for safety
  refresh_token_hash TEXT UNIQUE NOT NULL,

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