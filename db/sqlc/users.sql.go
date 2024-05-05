// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: users.sql

package frog_blossom_db

import (
	"context"
	"database/sql"
)

const createUsers = `-- name: CreateUsers :one
INSERT INTO Users (
  username,
  email,
  password,
  role,
  first_name,
  last_name,
  avatar_url,
  bio,
  updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING id, username, email, password, role, first_name, last_name, avatar_url, bio, created_at, updated_at
`

type CreateUsersParams struct {
	Username  string         `json:"username"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	Role      sql.NullString `json:"role"`
	FirstName sql.NullString `json:"first_name"`
	LastName  sql.NullString `json:"last_name"`
	AvatarUrl sql.NullString `json:"avatar_url"`
	Bio       sql.NullString `json:"bio"`
	UpdatedAt sql.NullTime   `json:"updated_at"`
}

func (q *Queries) CreateUsers(ctx context.Context, arg CreateUsersParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUsers,
		arg.Username,
		arg.Email,
		arg.Password,
		arg.Role,
		arg.FirstName,
		arg.LastName,
		arg.AvatarUrl,
		arg.Bio,
		arg.UpdatedAt,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.Role,
		&i.FirstName,
		&i.LastName,
		&i.AvatarUrl,
		&i.Bio,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteUsers = `-- name: DeleteUsers :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUsers(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteUsers, id)
	return err
}

const getUsers = `-- name: GetUsers :one
SELECT id, username, email, password, role, first_name, last_name, avatar_url, bio, created_at, updated_at FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUsers(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUsers, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.Role,
		&i.FirstName,
		&i.LastName,
		&i.AvatarUrl,
		&i.Bio,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, username, email, password, role, first_name, last_name, avatar_url, bio, created_at, updated_at FROM users
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Email,
			&i.Password,
			&i.Role,
			&i.FirstName,
			&i.LastName,
			&i.AvatarUrl,
			&i.Bio,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUsers = `-- name: UpdateUsers :one
UPDATE users
  SET username = $2,
  email = $3,
  password = $4,
  role = $5,
  first_name = $6,
  last_name = $7,
  avatar_url = $8,
  bio = $9,
  updated_at = $10
WHERE id = $1
RETURNING id, username, email, password, role, first_name, last_name, avatar_url, bio, created_at, updated_at
`

type UpdateUsersParams struct {
	ID        int64          `json:"id"`
	Username  string         `json:"username"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	Role      sql.NullString `json:"role"`
	FirstName sql.NullString `json:"first_name"`
	LastName  sql.NullString `json:"last_name"`
	AvatarUrl sql.NullString `json:"avatar_url"`
	Bio       sql.NullString `json:"bio"`
	UpdatedAt sql.NullTime   `json:"updated_at"`
}

func (q *Queries) UpdateUsers(ctx context.Context, arg UpdateUsersParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUsers,
		arg.ID,
		arg.Username,
		arg.Email,
		arg.Password,
		arg.Role,
		arg.FirstName,
		arg.LastName,
		arg.AvatarUrl,
		arg.Bio,
		arg.UpdatedAt,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.Role,
		&i.FirstName,
		&i.LastName,
		&i.AvatarUrl,
		&i.Bio,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
