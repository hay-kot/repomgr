-- name: RepoCreate :one
INSERT INTO 
  repository (remote_id, name, username, description, clone_url, clone_ssh_url, is_fork)
VALUES 
  (?, ?, ?, ?, ?, ?, ?)
RETURNING 
  *;  

-- name: RepoUpsert :one  
INSERT INTO 
  repository (remote_id, name, username, description, clone_url, clone_ssh_url, is_fork)  
VALUES 
  (?, ?, ?, ?, ?, ?, ?) 
ON CONFLICT (remote_id) 
DO UPDATE SET 
  name = EXCLUDED.name, 
  username = EXCLUDED.username, 
  description = EXCLUDED.description, 
  clone_url = EXCLUDED.clone_url, 
  clone_ssh_url = EXCLUDED.clone_ssh_url, 
  is_fork = EXCLUDED.is_fork  
RETURNING *;

-- name: ReposByUsernameLike :many
SELECT 
  * 
FROM 
  repository 
WHERE 
  username LIKE ?;

-- name: ReposByNameLike :many  
SELECT 
  * 
FROM  
  repository 
WHERE 
  name LIKE ?;  

-- name: ReposGetAll :many 
SELECT 
  * 
FROM  
  repository; 

-- name: RepoArtifacts :many
SELECT
  * 
FROM  
  repository_artifact 
WHERE 
  repository_id = ?;

-- name: RepoArtifactByType :many
SELECT 
  * 
FROM  
  repository_artifact 
WHERE 
  repository_id = ? 
  AND data_type = ?;

-- name: RepoUpsertArtifact :one 
INSERT INTO 
  repository_artifact (repository_id, data_type, data)  
VALUES 
  (?, ?, ?)
ON CONFLICT (repository_id, data_type)
DO UPDATE SET 
  data = EXCLUDED.data  
RETURNING
  *;

-- name: RepoUpdateArtifact :exec
UPDATE
  repository_artifact
SET
  data = ?
WHERE 
  repository_id = ?
  AND data_type = ?

