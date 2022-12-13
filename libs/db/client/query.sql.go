// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: query.sql

package client

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/tabbed/pqtype"
)

const getUserByEmail = `-- name: GetUserByEmail :one
select
    id,
    email
from auth.users
where email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (AuthUser, error) {
	row := q.queryRow(ctx, q.getUserByEmailStmt, getUserByEmail, email)
	var i AuthUser
	err := row.Scan(&i.ID, &i.Email)
	return i, err
}

const getUserEmailSyncHistory = `-- name: GetUserEmailSyncHistory :one
select
    user_id,
    history_id,
    synced_at,
    examples_collected_at,
    created_at,
    updated_at
from public.user_email_sync_history
where user_id = $1
`

func (q *Queries) GetUserEmailSyncHistory(ctx context.Context, userID uuid.UUID) (UserEmailSyncHistory, error) {
	row := q.queryRow(ctx, q.getUserEmailSyncHistoryStmt, getUserEmailSyncHistory, userID)
	var i UserEmailSyncHistory
	err := row.Scan(
		&i.UserID,
		&i.HistoryID,
		&i.SyncedAt,
		&i.ExamplesCollectedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserOAuthToken = `-- name: GetUserOAuthToken :one
select
    user_id,
    provider,
    token,
    is_valid,
    created_at,
    updated_at
from public.user_oauth_token
where user_id = $1 and provider = $2
`

type GetUserOAuthTokenParams struct {
	UserID   uuid.UUID `json:"user_id"`
	Provider string    `json:"provider"`
}

func (q *Queries) GetUserOAuthToken(ctx context.Context, arg GetUserOAuthTokenParams) (UserOauthToken, error) {
	row := q.queryRow(ctx, q.getUserOAuthTokenStmt, getUserOAuthToken, arg.UserID, arg.Provider)
	var i UserOauthToken
	err := row.Scan(
		&i.UserID,
		&i.Provider,
		&i.Token,
		&i.IsValid,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listOAuthTokensByProvider = `-- name: ListOAuthTokensByProvider :many
select
    user_id,
    provider,
    token,
    is_valid,
    created_at,
    updated_at
from public.user_oauth_token
where provider = $1
`

func (q *Queries) ListOAuthTokensByProvider(ctx context.Context, provider string) ([]UserOauthToken, error) {
	rows, err := q.query(ctx, q.listOAuthTokensByProviderStmt, listOAuthTokensByProvider, provider)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UserOauthToken
	for rows.Next() {
		var i UserOauthToken
		if err := rows.Scan(
			&i.UserID,
			&i.Provider,
			&i.Token,
			&i.IsValid,
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

const listValidOAuthTokensByProvider = `-- name: ListValidOAuthTokensByProvider :many
select
    user_id,
    provider,
    token,
    is_valid,
    created_at,
    updated_at
from public.user_oauth_token
where provider = $1 and is_valid = true
`

func (q *Queries) ListValidOAuthTokensByProvider(ctx context.Context, provider string) ([]UserOauthToken, error) {
	rows, err := q.query(ctx, q.listValidOAuthTokensByProviderStmt, listValidOAuthTokensByProvider, provider)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UserOauthToken
	for rows.Next() {
		var i UserOauthToken
		if err := rows.Scan(
			&i.UserID,
			&i.Provider,
			&i.Token,
			&i.IsValid,
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

const upsertUserEmailSyncHistory = `-- name: UpsertUserEmailSyncHistory :exec
insert into public.user_email_sync_history(user_id, history_id, synced_at, examples_collected_at)
values ($1, $2, $3, $4)
on conflict (user_id) 
do update set 
    history_id = excluded.history_id,
    synced_at = excluded.synced_at,
    examples_collected_at = excluded.examples_collected_at
`

type UpsertUserEmailSyncHistoryParams struct {
	UserID              uuid.UUID    `json:"user_id"`
	HistoryID           int64        `json:"history_id"`
	SyncedAt            sql.NullTime `json:"synced_at"`
	ExamplesCollectedAt sql.NullTime `json:"examples_collected_at"`
}

func (q *Queries) UpsertUserEmailSyncHistory(ctx context.Context, arg UpsertUserEmailSyncHistoryParams) error {
	_, err := q.exec(ctx, q.upsertUserEmailSyncHistoryStmt, upsertUserEmailSyncHistory,
		arg.UserID,
		arg.HistoryID,
		arg.SyncedAt,
		arg.ExamplesCollectedAt,
	)
	return err
}

const upsertUserOAuthToken = `-- name: UpsertUserOAuthToken :exec
insert into public.user_oauth_token (user_id, provider, token, is_valid)
values ($1, $2, $3, $4)
on conflict (user_id, provider) 
do update set
    token = excluded.token,
    is_valid = excluded.is_valid
`

type UpsertUserOAuthTokenParams struct {
	UserID   uuid.UUID             `json:"user_id"`
	Provider string                `json:"provider"`
	Token    pqtype.NullRawMessage `json:"token"`
	IsValid  bool                  `json:"is_valid"`
}

func (q *Queries) UpsertUserOAuthToken(ctx context.Context, arg UpsertUserOAuthTokenParams) error {
	_, err := q.exec(ctx, q.upsertUserOAuthTokenStmt, upsertUserOAuthToken,
		arg.UserID,
		arg.Provider,
		arg.Token,
		arg.IsValid,
	)
	return err
}
