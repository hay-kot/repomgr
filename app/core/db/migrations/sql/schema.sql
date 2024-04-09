CREATE TABLE repository (
  id            INTEGER PRIMARY KEY,
  remote_id     TEXT NOT NULL UNIQUE,
  name          TEXT NOT NULL,
  username      TEXT NOT NULL,
  description   TEXT NOT NULL,
  clone_url     TEXT NOT NULL,
  clone_ssh_url TEXT NOT NULL,
  is_fork       BOOLEAN NOT NULL
);

CREATE TABLE repository_artifact (
  id            INTEGER PRIMARY KEY,
  type          TEXT    NOT NULL,
  data          BLOB    NOT NULL,
  repository_id INTEGER NOT NULL,
  FOREIGN KEY (repository_id) REFERENCES repository(id) ON DELETE CASCADE
); 
