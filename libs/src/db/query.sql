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
where provider = $1 and is_valid = $2
order by created_at desc
limit $3
offset $4;


-- name: ListCandidateOAuthTokens :many
select
    *
from public.candidate_oauth_token
where provider = $1 and is_valid = $2
limit $3
offset $4;

-- name: ListRecruiterOAuthTokens :many
select
    *
from public.recruiter_oauth_token
where provider = $1 and is_valid = $2
limit $3
offset $4;

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
    inbox_type,
    email,
    history_id,
    synced_at,
    created_at,
    updated_at
from public.user_email_sync_history
where user_id = $1 and inbox_type = $2 and email = $3;

-- name: UpsertUserEmailSyncHistory :exec
insert into public.user_email_sync_history(user_id, inbox_type, email, history_id, synced_at)
values ($1, $2, $3, $4, $5)
on conflict (user_id, inbox_type, email)
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
offset $3;

-- name: GetUserEmailJob :one
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
where job_id = $1;

-- name: InsertUserEmailJob :exec
insert into public.user_email_job(user_id, user_email, email_thread_id, emailed_at, company, job_title, data)
values ($1, $2, $3, $4, $5, $6, $7);

-- name: CountUserEmailJobs :one
select count(*) as cnt
from public.user_email_job
where user_id = $1;

-- name: GetRecruiterByEmail :one
select 
    user_id,
    email,
    first_name,
    last_name,
    company_id,
    created_at,
    updated_at
from recruiter
where email = $1;
