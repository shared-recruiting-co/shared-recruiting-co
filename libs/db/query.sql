-- name: GetUserProfileByEmail :one
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
where email = $1;

-- name: ListUserOAuthTokens :many
select
    user_id,
    provider,
    token,
    is_valid,
    created_at,
    updated_at
from public.user_oauth_token
where provider = $1 and is_valid = $2;

-- name: GetUserOAuthToken :one
select
    user_id,
    provider,
    token,
    is_valid,
    created_at,
    updated_at
from public.user_oauth_token
where user_id = $1 and provider = $2;

-- name: UpsertUserOAuthToken :exec
insert into public.user_oauth_token (user_id, provider, token, is_valid)
values ($1, $2, $3, $4)
on conflict (user_id, provider)
do update set
    token = excluded.token,
    is_valid = excluded.is_valid;

-- name: GetUserEmailSyncHistory :one
select
    user_id,
    history_id,
    synced_at,
    created_at,
    updated_at
from public.user_email_sync_history
where user_id = $1;

-- name: UpsertUserEmailSyncHistory :exec
insert into public.user_email_sync_history(user_id, history_id, synced_at)
values ($1, $2, $3)
on conflict (user_id)
do update set
    history_id = excluded.history_id,
    synced_at = excluded.synced_at;

-- name: IncrementUserEmailStat :exec
insert into public.user_email_stat(user_id, email, stat_id, stat_value)
values ($1, $2, $3, $4)
on conflict (user_id, email, stat_id)
do update set
    stat_value = user_email_stat.stat_value + excluded.stat_value;

-- name: ListUserEmailJobs :many
select
    job_id,
    user_id,
    user_email,
    email_thread_id,
    company,
    job_title,
    data,
    created_at,
    updated_at
from public.user_email_job
where user_id = $1;

-- name: GetUserEmailJob :one
select
    job_id,
    user_id,
    user_email,
    email_thread_id,
    company,
    job_title,
    data,
    created_at,
    updated_at
from public.user_email_job
where job_id = $1;

-- name: InsertUserEmailJob :exec
insert into public.user_email_job(user_id, user_email, email_thread_id, company, job_title, data)
values ($1, $2, $3, $4, $5, $6);
