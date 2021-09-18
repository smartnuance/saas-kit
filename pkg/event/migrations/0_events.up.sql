--Events group workshops and assign them to an instance and an owner.
CREATE TABLE IF NOT EXISTS events(
  --use ObjectId as primary key
  id CHAR(20) PRIMARY KEY,
  info jsonb NOT NULL,
  --date range
  starts date NOT NULL,
  ends date,
  instance_id bigint NOT NULL,
  owner_id CHAR(20),
  created_at timestamp with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp with time zone NOT NULL DEFAULT NOW(),
  deleted_at timestamp with time zone
);
CREATE INDEX instance_idx ON events(instance_id);
CREATE INDEX owner_idx ON events(owner_id);
CREATE TABLE IF NOT EXISTS workshops(
  --use ObjectId as primary key
  id CHAR(20) PRIMARY KEY,
  info jsonb NOT NULL,
  --time range
  starts timestamp with time zone NOT NULL,
  ends timestamp with time zone,
  --workshops have to be assigned to an event
  event_id CHAR(20) NOT NULL,
  participants jsonb NOT NULL,
  created_at timestamp with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp with time zone NOT NULL DEFAULT NOW(),
  deleted_at timestamp with time zone,
  CONSTRAINT fk_event FOREIGN KEY(event_id) REFERENCES events(id)
);
