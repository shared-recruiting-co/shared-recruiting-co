-- name: GetUserByEmail :one
select
    id,
    email
from auth.users
where email = $1;

-- name: ListOAuthTokensByProvider :many
select
    user_id,
    provider,
    token,
    is_valid,
    created_at,
    updated_at
from public.user_oauth_token
where provider = $1;

-- name: ListValidOAuthTokensByProvider :many
select
    user_id,
    provider,
    token,
    is_valid,
    created_at,
    updated_at
from public.user_oauth_token
where provider = $1 and is_valid = true;

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
    examples_collected_at,
    created_at,
    updated_at
from public.user_email_sync_history
where user_id = $1;

-- name: UpsertUserEmailSyncHistory :exec
insert into public.user_email_sync_history(user_id, history_id, synced_at, examples_collected_at)
values ($1, $2, $3, $4)
on conflict (user_id) 
do update set 
    history_id = excluded.history_id,
    synced_at = excluded.synced_at,
    examples_collected_at = excluded.examples_collected_at;
