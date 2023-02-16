// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: query.sql

package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	null "gopkg.in/guregu/null.v4"
)

const countUserEmailJobs = `-- name: CountUserEmailJobs :one
select count(*) as cnt
from public.user_email_job
where user_id = $1
`

func (q *Queries) CountUserEmailJobs(ctx context.Context, userID uuid.UUID) (int64, error) {
	row := q.queryRow(ctx, q.countUserEmailJobsStmt, countUserEmailJobs, userID)
	var cnt int64
	err := row.Scan(&cnt)
	return cnt, err
}

const getRecruiterByEmail = `-- name: GetRecruiterByEmail :one
select 
    user_id,
    email,
    first_name,
    last_name,
    company_id,
    created_at,
    updated_at
from recruiter
where email = $1
`

type GetRecruiterByEmailRow struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CompanyID uuid.UUID `json:"company_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) GetRecruiterByEmail(ctx context.Context, email string) (GetRecruiterByEmailRow, error) {
	row := q.queryRow(ctx, q.getRecruiterByEmailStmt, getRecruiterByEmail, email)
	var i GetRecruiterByEmailRow
	err := row.Scan(
		&i.UserID,
		&i.Email,
		&i.FirstName,
		&i.LastName,
		&i.CompanyID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserEmailJob = `-- name: GetUserEmailJob :one
select
    job_id,
    user_id,
    user_email,
    email_thread_id,
    emailed_at,
    company,
    job_title,
    data,
    created_at,
    updated_at
from public.user_email_job
where job_id = $1
`

func (q *Queries) GetUserEmailJob(ctx context.Context, jobID uuid.UUID) (UserEmailJob, error) {
	row := q.queryRow(ctx, q.getUserEmailJobStmt, getUserEmailJob, jobID)
	var i UserEmailJob
	err := row.Scan(
		&i.JobID,
		&i.UserID,
		&i.UserEmail,
		&i.EmailThreadID,
		&i.EmailedAt,
		&i.Company,
		&i.JobTitle,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserEmailSyncHistory = `-- name: GetUserEmailSyncHistory :one
select
    user_id,
    inbox_type,
    email,
    history_id,
    synced_at,
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
		&i.InboxType,
		&i.Email,
		&i.HistoryID,
		&i.SyncedAt,
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

const getUserProfileByEmail = `-- name: GetUserProfileByEmail :one
select
    user_id,
    email,
    first_name,
    last_name,
    is_active,
    auto_archive,
    auto_contribute,
    created_at,
    updated_at
from public.user_profile
where email = $1
`

func (q *Queries) GetUserProfileByEmail(ctx context.Context, email string) (UserProfile, error) {
	row := q.queryRow(ctx, q.getUserProfileByEmailStmt, getUserProfileByEmail, email)
	var i UserProfile
	err := row.Scan(
		&i.UserID,
		&i.Email,
		&i.FirstName,
		&i.LastName,
		&i.IsActive,
		&i.AutoArchive,
		&i.AutoContribute,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const incrementUserEmailStat = `-- name: IncrementUserEmailStat :exec
insert into public.user_email_stat(user_id, email, stat_id, stat_value)
values ($1, $2, $3, $4)
on conflict (user_id, email, stat_id)
do update set
    stat_value = user_email_stat.stat_value + excluded.stat_value
`

type IncrementUserEmailStatParams struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	StatID    string    `json:"stat_id"`
	StatValue int32     `json:"stat_value"`
}

func (q *Queries) IncrementUserEmailStat(ctx context.Context, arg IncrementUserEmailStatParams) error {
	_, err := q.exec(ctx, q.incrementUserEmailStatStmt, incrementUserEmailStat,
		arg.UserID,
		arg.Email,
		arg.StatID,
		arg.StatValue,
	)
	return err
}

const insertUserEmailJob = `-- name: InsertUserEmailJob :exec
insert into public.user_email_job(user_id, user_email, email_thread_id, emailed_at, company, job_title, data)
values ($1, $2, $3, $4, $5, $6, $7)
`

type InsertUserEmailJobParams struct {
	UserID        uuid.UUID       `json:"user_id"`
	UserEmail     string          `json:"user_email"`
	EmailThreadID string          `json:"email_thread_id"`
	EmailedAt     time.Time       `json:"emailed_at"`
	Company       string          `json:"company"`
	JobTitle      string          `json:"job_title"`
	Data          json.RawMessage `json:"data"`
}

func (q *Queries) InsertUserEmailJob(ctx context.Context, arg InsertUserEmailJobParams) error {
	_, err := q.exec(ctx, q.insertUserEmailJobStmt, insertUserEmailJob,
		arg.UserID,
		arg.UserEmail,
		arg.EmailThreadID,
		arg.EmailedAt,
		arg.Company,
		arg.JobTitle,
		arg.Data,
	)
	return err
}

const listCandidateOAuthTokens = `-- name: ListCandidateOAuthTokens :many
select
    user_id, provider, token, is_valid, created_at, updated_at
from public.candidate_oauth_token
where provider = $1 and is_valid = $2
limit $3
offset $4
`

type ListCandidateOAuthTokensParams struct {
	Provider string `json:"provider"`
	IsValid  bool   `json:"is_valid"`
	Limit    int32  `json:"limit"`
	Offset   int32  `json:"offset"`
}

func (q *Queries) ListCandidateOAuthTokens(ctx context.Context, arg ListCandidateOAuthTokensParams) ([]CandidateOauthToken, error) {
	rows, err := q.query(ctx, q.listCandidateOAuthTokensStmt, listCandidateOAuthTokens,
		arg.Provider,
		arg.IsValid,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CandidateOauthToken
	for rows.Next() {
		var i CandidateOauthToken
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

const listRecruiterOAuthTokens = `-- name: ListRecruiterOAuthTokens :many
select
    user_id, provider, token, is_valid, created_at, updated_at
from public.recruiter_oauth_token
where provider = $1 and is_valid = $2
limit $3
offset $4
`

type ListRecruiterOAuthTokensParams struct {
	Provider string `json:"provider"`
	IsValid  bool   `json:"is_valid"`
	Limit    int32  `json:"limit"`
	Offset   int32  `json:"offset"`
}

func (q *Queries) ListRecruiterOAuthTokens(ctx context.Context, arg ListRecruiterOAuthTokensParams) ([]RecruiterOauthToken, error) {
	rows, err := q.query(ctx, q.listRecruiterOAuthTokensStmt, listRecruiterOAuthTokens,
		arg.Provider,
		arg.IsValid,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RecruiterOauthToken
	for rows.Next() {
		var i RecruiterOauthToken
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

const listUserEmailJobs = `-- name: ListUserEmailJobs :many
select
    job_id,
    user_id,
    user_email,
    email_thread_id,
    emailed_at,
    company,
    job_title,
    data,
    created_at,
    updated_at
from public.user_email_job
where user_id = $1
order by emailed_at desc
limit $2
offset $3
`

type ListUserEmailJobsParams struct {
	UserID uuid.UUID `json:"user_id"`
	Limit  int32     `json:"limit"`
	Offset int32     `json:"offset"`
}

func (q *Queries) ListUserEmailJobs(ctx context.Context, arg ListUserEmailJobsParams) ([]UserEmailJob, error) {
	rows, err := q.query(ctx, q.listUserEmailJobsStmt, listUserEmailJobs, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UserEmailJob
	for rows.Next() {
		var i UserEmailJob
		if err := rows.Scan(
			&i.JobID,
			&i.UserID,
			&i.UserEmail,
			&i.EmailThreadID,
			&i.EmailedAt,
			&i.Company,
			&i.JobTitle,
			&i.Data,
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

const listUserOAuthTokens = `-- name: ListUserOAuthTokens :many
select
    user_id,
    provider,
    token,
    is_valid,
    created_at,
    updated_at
from public.user_oauth_token
where provider = $1 and is_valid = $2
order by created_at desc
limit $3
offset $4
`

type ListUserOAuthTokensParams struct {
	Provider string `json:"provider"`
	IsValid  bool   `json:"is_valid"`
	Limit    int32  `json:"limit"`
	Offset   int32  `json:"offset"`
}

func (q *Queries) ListUserOAuthTokens(ctx context.Context, arg ListUserOAuthTokensParams) ([]UserOauthToken, error) {
	rows, err := q.query(ctx, q.listUserOAuthTokensStmt, listUserOAuthTokens,
		arg.Provider,
		arg.IsValid,
		arg.Limit,
		arg.Offset,
	)
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
insert into public.user_email_sync_history(user_id, inbox_type, email, history_id, synced_at)
values ($1, $2, $3, $4, $5)
on conflict (user_id)
do update set
    history_id = excluded.history_id,
    inbox_type = excluded.inbox_type,
    email = excluded.email,
    synced_at = excluded.synced_at
`

type UpsertUserEmailSyncHistoryParams struct {
	UserID    uuid.UUID   `json:"user_id"`
	InboxType InboxType   `json:"inbox_type"`
	Email     null.String `json:"email"`
	HistoryID int64       `json:"history_id"`
	SyncedAt  time.Time   `json:"synced_at"`
}

func (q *Queries) UpsertUserEmailSyncHistory(ctx context.Context, arg UpsertUserEmailSyncHistoryParams) error {
	_, err := q.exec(ctx, q.upsertUserEmailSyncHistoryStmt, upsertUserEmailSyncHistory,
		arg.UserID,
		arg.InboxType,
		arg.Email,
		arg.HistoryID,
		arg.SyncedAt,
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
	UserID   uuid.UUID       `json:"user_id"`
	Provider string          `json:"provider"`
	Token    json.RawMessage `json:"token"`
	IsValid  bool            `json:"is_valid"`
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
