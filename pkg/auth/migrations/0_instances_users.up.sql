CREATE TABLE IF NOT EXISTS instances(
    id bigint PRIMARY KEY,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name text
);

CREATE TABLE IF NOT EXISTS profiles(
    id bigint PRIMARY KEY,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    user_id bigint,
    roles text[],
    instance_id bigint
);

CREATE TABLE IF NOT EXISTS users(
    id bigint PRIMARY KEY,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name text,
    email text,
    username text,
    activated_at timestamp with time zone
);
