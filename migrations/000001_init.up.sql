CREATE SCHEMA todoapp;

CREATE TABLE todoapp.users (
  id SERIAL PRIMARY KEY,
  version BIGINT NOT NULL DEFAULT 1,
  full_name VARCHAR(100) NOT NULL CHECK (full_name <> ''),
  phone_number VARCHAR(16) UNIQUE CHECK (
      phone_number ~ '^\+[1-9][0-9]{7,14}$'
  )
);

CREATE TABLE todoapp.tasks (
  id SERIAL PRIMARY KEY,
  version BIGINT NOT NULL DEFAULT 1,
  title VARCHAR(100) NOT NULL CHECK (title <> ''),
  description VARCHAR(1000) CHECK (title <> ''),
  completed BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL,
  completed_at TIMESTAMPTZ,
  author_id INTEGER NOT NULL REFERENCES todoapp.users(id)

  CHECK (
    (completed=FALSE AND completed_at IS NULL)
    OR
    (completed=TRUE AND completed_at IS NOT NULL AND completed_at >= created_at)
  )
);