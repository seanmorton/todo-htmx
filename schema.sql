CREATE TABLE users (
  id            INTEGER     PRIMARY KEY,
  name          TEXT        NOT NULL,
  created_at    TIMESTAMP   DEFAULT current_timestamp
);

CREATE TABLE projects (
  id            INTEGER     PRIMARY KEY,
  title         TEXT        NOT NULL,
  created_at    TIMESTAMP   DEFAULT current_timestamp
);

CREATE TABLE tasks (
  id            INTEGER     PRIMARY KEY,
  project_id    INTEGER,
  assignee_id   INTEGER,
  title         TEXT        NOT NULL,
  description   TEXT,
  due_date      TIMESTAMP,
  completed_at  TIMESTAMP,
  recur_policy  BLOB,
  created_at    TIMESTAMP   DEFAULT current_timestamp,
  FOREIGN KEY (project_id)  REFERENCES projects(id),
  FOREIGN KEY (assignee_id) REFERENCES users(id)
);
