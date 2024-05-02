// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: nicknames.sql

package wordle

import (
	"context"
)

const countNicknameByDiscordIdAndServerId = `-- name: CountNicknameByDiscordIdAndServerId :one
SELECT count(*)
FROM nicknames
where discord_id = $1
  and server_id = $2
`

type CountNicknameByDiscordIdAndServerIdParams struct {
	DiscordID string `json:"discord_id"`
	ServerID  string `json:"server_id"`
}

func (q *Queries) CountNicknameByDiscordIdAndServerId(ctx context.Context, arg CountNicknameByDiscordIdAndServerIdParams) (int64, error) {
	row := q.queryRow(ctx, q.countNicknameByDiscordIdAndServerIdStmt, countNicknameByDiscordIdAndServerId, arg.DiscordID, arg.ServerID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createNickname = `-- name: CreateNickname :one
INSERT INTO nicknames (discord_id, server_id, nickname)
VALUES ($1, $2, $3)
RETURNING discord_id, server_id, nickname
`

type CreateNicknameParams struct {
	DiscordID string `json:"discord_id"`
	ServerID  string `json:"server_id"`
	Nickname  string `json:"nickname"`
}

func (q *Queries) CreateNickname(ctx context.Context, arg CreateNicknameParams) (Nickname, error) {
	row := q.queryRow(ctx, q.createNicknameStmt, createNickname, arg.DiscordID, arg.ServerID, arg.Nickname)
	var i Nickname
	err := row.Scan(&i.DiscordID, &i.ServerID, &i.Nickname)
	return i, err
}

const deleteNickname = `-- name: DeleteNickname :exec
DELETE
FROM nicknames
WHERE discord_id = $1
`

func (q *Queries) DeleteNickname(ctx context.Context, discordID string) error {
	_, err := q.exec(ctx, q.deleteNicknameStmt, deleteNickname, discordID)
	return err
}

const getNickname = `-- name: GetNickname :one
SELECT discord_id, server_id, nickname
FROM nicknames
WHERE discord_id = $1
and server_id = $2
LIMIT 1
`

type GetNicknameParams struct {
	DiscordID string `json:"discord_id"`
	ServerID  string `json:"server_id"`
}

func (q *Queries) GetNickname(ctx context.Context, arg GetNicknameParams) (Nickname, error) {
	row := q.queryRow(ctx, q.getNicknameStmt, getNickname, arg.DiscordID, arg.ServerID)
	var i Nickname
	err := row.Scan(&i.DiscordID, &i.ServerID, &i.Nickname)
	return i, err
}

const getNicknamesByDiscordId = `-- name: GetNicknamesByDiscordId :many
SELECT discord_id, server_id, nickname
FROM nicknames
where discord_id = $1
`

func (q *Queries) GetNicknamesByDiscordId(ctx context.Context, discordID string) ([]Nickname, error) {
	rows, err := q.query(ctx, q.getNicknamesByDiscordIdStmt, getNicknamesByDiscordId, discordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Nickname
	for rows.Next() {
		var i Nickname
		if err := rows.Scan(&i.DiscordID, &i.ServerID, &i.Nickname); err != nil {
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

const listNicknames = `-- name: ListNicknames :many
SELECT discord_id, server_id, nickname
FROM nicknames
ORDER BY discord_id
`

func (q *Queries) ListNicknames(ctx context.Context) ([]Nickname, error) {
	rows, err := q.query(ctx, q.listNicknamesStmt, listNicknames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Nickname
	for rows.Next() {
		var i Nickname
		if err := rows.Scan(&i.DiscordID, &i.ServerID, &i.Nickname); err != nil {
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

const updateNickname = `-- name: UpdateNickname :one
UPDATE nicknames
set nickname = $2
where discord_id = $1
  and server_id = $3
RETURNING discord_id, server_id, nickname
`

type UpdateNicknameParams struct {
	DiscordID string `json:"discord_id"`
	Nickname  string `json:"nickname"`
	ServerID  string `json:"server_id"`
}

func (q *Queries) UpdateNickname(ctx context.Context, arg UpdateNicknameParams) (Nickname, error) {
	row := q.queryRow(ctx, q.updateNicknameStmt, updateNickname, arg.DiscordID, arg.Nickname, arg.ServerID)
	var i Nickname
	err := row.Scan(&i.DiscordID, &i.ServerID, &i.Nickname)
	return i, err
}
