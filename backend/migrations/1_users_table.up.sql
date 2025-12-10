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
  last_login_at TIMESTAMP WITH TIME ZONE
);

-- Every user must have a unique identity provider subject
CREATE UNIQUE INDEX ux_users_subject_provider
  ON users(idp_provider, idp_subject);