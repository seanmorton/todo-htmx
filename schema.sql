CREATE TABLE tasks (
  id            INTEGER     PRIMARY KEY,
  title         TEXT        NOT NULL,
  description   TEXT,
  assignee      TEXT,
  due_date      TIMESTAMP,
  completed_at  TIMESTAMP,
  recur_policy  BLOB,
  created_at    TIMESTAMP   DEFAULT current_timestamp
);
