-- name: GetUserProfileByEmail :one
select
    user_profile.user_id,
    user_profile.email,
    user_profile.first_name,
    user_profile.last_name,
    user_profile.is_active,
    user_profile.auto_archive,
    user_profile.auto_contribute,
    user_profile.created_at,
    user_profile.updated_at
from public.user_profile
inner join public.user_oauth_token using (user_id)
where user_profile.email = $1 OR user_oauth_token.email = $1
limit 1;

-- name: ListUserOAuthTokens :many
select
    user_id,
    email,
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
    email,
    provider,
    token,
    is_valid,
    created_at,
    updated_at
from public.user_oauth_token
where user_id = $1 and email = $2 and provider = $3;

-- name: UpsertUserOAuthToken :exec
insert into public.user_oauth_token (user_id, email, provider, token, is_valid)
values ($1, $2, $3, $4, $5)
on conflict (user_id, email, provider)
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

-- name: DeleteUserEmailJobByEmailThreadID :exec
delete from public.user_email_job
where user_email = $1 and email_thread_id = $2;

-- name: CountUserEmailJobs :one
select count(*) as cnt
from public.user_email_job
where user_id = $1;

-- name: GetRecruiterByEmail :one
select 
    recruiter.user_id,
    recruiter.email,
    recruiter.first_name,
    recruiter.last_name,
    recruiter.email_settings,
    recruiter.company_id,
    recruiter.created_at,
    recruiter.updated_at
from recruiter
inner join public.user_oauth_token using (user_id)
where recruiter.email = $1 OR user_oauth_token.email = $1
limit 1;

-- name: GetRecruiterOutboundMessage :one
select
    recruiter_id,
    message_id,
    internal_message_id,
    from_email,
    to_email,
    sent_at,
    template_id,
    created_at,
    updated_at
from public.recruiter_outbound_message
where recruiter_id = $1 and message_id = $2;

-- name: GetRecruiterOutboundMessageByRecipient :one
select
    recruiter_id,
    message_id,
    internal_message_id,
    from_email,
    to_email,
    sent_at,
    template_id,
    created_at,
    updated_at
from public.recruiter_outbound_message
where to_email = $1 and internal_message_id = $2;

-- name: InsertRecruiterOutboundMessage :exec
insert into public.recruiter_outbound_message(recruiter_id, message_id, internal_message_id, from_email, to_email, sent_at, template_id)
values ($1, $2, $3, $4, $5, $6, $7);

-- name: InsertRecruiterOutboundTemplate :one
insert into public.recruiter_outbound_template(recruiter_id, job_id, subject, body, normalized_content, metadata)
values ($1, $2, $3, $4, $5, $6)
returning *;

-- name: ListSimilarRecruiterOutboundTemplates :many
select
    template_id,
    recruiter_id,
    job_id,
    subject,
    body,
    metadata,
    created_at,
    updated_at,
    similarity(normalized_content, @input::text) as "similarity"
from public.recruiter_outbound_template
where recruiter_id = @user_id::uuid
and normalized_content % @input::text
order by 9 desc;
