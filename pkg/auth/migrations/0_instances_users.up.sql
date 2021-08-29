CREATE TABLE IF NOT EXISTS instances(
  id bigserial PRIMARY KEY,
  name text NOT NULL,
  url text NOT NULL,
  created_at timestamp with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp with time zone NOT NULL DEFAULT NOW(),
  deleted_at timestamp with time zone
);
CREATE INDEX instance_name_idx ON instances(name);
CREATE INDEX instance_url_idx ON instances(url);
CREATE TABLE IF NOT EXISTS users(
  id bigserial PRIMARY KEY,
  name text,
  email text NOT NULL UNIQUE,
  password bytea NOT NULL,
  activated_at timestamp with time zone,
  created_at timestamp with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp with time zone NOT NULL DEFAULT NOW(),
  deleted_at timestamp with time zone
);
CREATE INDEX email_idx ON users(email);
CREATE TABLE IF NOT EXISTS profiles(
  id bigserial PRIMARY KEY,
  user_id bigint NOT NULL,
  instance_id bigint NOT NULL,
  role text,
  created_at timestamp with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp with time zone NOT NULL DEFAULT NOW(),
  deleted_at timestamp with time zone,
  CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id),
  CONSTRAINT fk_instance FOREIGN KEY(instance_id) REFERENCES instances(id)
);
CREATE TABLE IF NOT EXISTS tokens(
  id bigserial PRIMARY KEY,
  user_id bigint NOT NULL,
  profile_id bigint NOT NULL,
  token text NOT NULL,
  created_at timestamp with time zone NOT NULL DEFAULT NOW(),
  expires_at timestamp with time zone NOT NULL,
  CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id),
  CONSTRAINT fk_profile FOREIGN KEY(profile_id) REFERENCES profiles(id)
);
CREATE INDEX token_idx ON tokens(user_id, profile_id, token);
