CREATE TABLE repository (
  id            INTEGER PRIMARY KEY,
  remote_id     TEXT NOT NULL UNIQUE,
  name          TEXT NOT NULL,
  username      TEXT NOT NULL,
  description   TEXT NOT NULL,
  html_url      TEXT NOT NULL, 
  clone_url     TEXT NOT NULL,
  clone_ssh_url TEXT NOT NULL,
  is_fork       BOOLEAN NOT NULL,
  fork_url      TEXT NOT NULL
);

CREATE TABLE repository_artifact (
  id            INTEGER PRIMARY KEY,
  data_type     TEXT    NOT NULL,
  data          BLOB    NOT NULL,
  repository_id INTEGER NOT NULL,
  FOREIGN KEY (repository_id) REFERENCES repository(id) ON DELETE CASCADE,
  UNIQUE(repository_id, data_type)
);
