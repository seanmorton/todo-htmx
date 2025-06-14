CREATE TABLE IF NOT EXISTS users (
  id            INTEGER     PRIMARY KEY,
  name          TEXT        NOT NULL UNIQUE,
  created_at    TIMESTAMP   DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS projects (
  id            INTEGER     PRIMARY KEY,
  name          TEXT        NOT NULL UNIQUE,
  created_at    TIMESTAMP   DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS tasks (
  id            INTEGER     PRIMARY KEY,
  project_id    INTEGER     REFERENCES projects(id) ON DELETE CASCADE NOT NULL,
  assignee_id   INTEGER     REFERENCES users(id) ON DELETE SET NULL,
  title         TEXT        NOT NULL,
  description   TEXT,
  due_date      TIMESTAMP,
  completed_at  TIMESTAMP,
  recur_policy  BLOB,
  created_at    TIMESTAMP   DEFAULT current_timestamp
);
