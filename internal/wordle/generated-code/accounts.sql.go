// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: accounts.sql

package wordle

import (
	"context"
)

const countAccountsByDiscordId = `-- name: CountAccountsByDiscordId :one
SELECT count(*)
FROM accounts
where discord_id = $1
`

func (q *Queries) CountAccountsByDiscordId(ctx context.Context, discordID string) (int64, error) {
	row := q.queryRow(ctx, q.countAccountsByDiscordIdStmt, countAccountsByDiscordId, discordID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createAccount = `-- name: CreateAccount :one
INSERT INTO accounts (discord_id)
VALUES ($1)
RETURNING discord_id, time_zone
`

func (q *Queries) CreateAccount(ctx context.Context, discordID string) (Account, error) {
	row := q.queryRow(ctx, q.createAccountStmt, createAccount, discordID)
	var i Account
	err := row.Scan(&i.DiscordID, &i.TimeZone)
	return i, err
}

const deleteAccount = `-- name: DeleteAccount :exec
DELETE
FROM accounts
WHERE discord_id = $1
`

func (q *Queries) DeleteAccount(ctx context.Context, discordID string) error {
	_, err := q.exec(ctx, q.deleteAccountStmt, deleteAccount, discordID)
	return err
}

const getAccount = `-- name: GetAccount :one
SELECT discord_id, time_zone
FROM accounts
WHERE discord_id = $1
LIMIT 1
`

func (q *Queries) GetAccount(ctx context.Context, discordID string) (Account, error) {
	row := q.queryRow(ctx, q.getAccountStmt, getAccount, discordID)
	var i Account
	err := row.Scan(&i.DiscordID, &i.TimeZone)
	return i, err
}

const listAccounts = `-- name: ListAccounts :many
SELECT discord_id, time_zone
FROM accounts
ORDER BY discord_id
`

func (q *Queries) ListAccounts(ctx context.Context) ([]Account, error) {
	rows, err := q.query(ctx, q.listAccountsStmt, listAccounts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Account
	for rows.Next() {
		var i Account
		if err := rows.Scan(&i.DiscordID, &i.TimeZone); err != nil {
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

const updateTimeZone = `-- name: UpdateTimeZone :one
UPDATE accounts
set time_zone = $2
where discord_id = $1
RETURNING discord_id, time_zone
`

type UpdateTimeZoneParams struct {
	DiscordID string `json:"discord_id"`
	TimeZone  string `json:"time_zone"`
}

func (q *Queries) UpdateTimeZone(ctx context.Context, arg UpdateTimeZoneParams) (Account, error) {
	row := q.queryRow(ctx, q.updateTimeZoneStmt, updateTimeZone, arg.DiscordID, arg.TimeZone)
	var i Account
	err := row.Scan(&i.DiscordID, &i.TimeZone)
	return i, err
}
