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
  AND type = ?;

-- name: RepoCreateArtifact :one 
INSERT INTO 
  repository_artifact (repository_id, type, data)  
VALUES 
  (?, ?, ?)
RETURNING
  *;
