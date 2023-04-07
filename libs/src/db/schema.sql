--------------------------------
-- Start: Supabase auth.users table
--------------------------------
-- Include a simplified Supabase auth.users table in the schema for SQLC compilation
create schema if not exists auth;
create table auth.users (
    id uuid primary key,
    email text not null
);
--------------------------------
-- End: Supabase auth.users table
--------------------------------

--------------------------------
-- Start: Postgres Extensions
--------------------------------

-- Add moddatetime extension
create extension if not exists moddatetime schema extensions;

-- Add pg_trgm extension
create extension if not exists pg_trgm schema extensions;

-- Enable the pgtap extension for testing
create extension pgtap with schema extensions;

--------------------------------
-- End: Postgres Extensions
--------------------------------

-- Enable Suapbase Realtime
begin;
  -- remove the supabase_realtime publication
  drop publication if exists supabase_realtime;

  -- re-create the supabase_realtime publication with no tables
  create publication supabase_realtime;
commit;

--------------------------------
-- Start: User OAuth Token Table
--------------------------------
create table public.user_oauth_token (
    user_id uuid references auth.users(id) on delete cascade not null,
    email text not null,
    provider text not null,
    token jsonb not null,
    is_valid boolean not null default true,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),

    -- Add a unique constraint to prevent duplicate tokens
    unique (email, provider),
    primary key (user_id, email, provider)
);

create trigger handle_updated_at_user_oauth_token before update on public.user_oauth_token
  for each row execute procedure moddatetime (updated_at);

alter table public.user_oauth_token enable row level security;

create policy "Users can view their own oauth tokens"
  on public.user_oauth_token for select
  using ( auth.uid() = user_id );

create policy "Users can insert their own oauth tokens."
  on public.user_oauth_token for insert
  with check ( auth.uid() = user_id );

create policy "Users can update own oauth tokens."
  on public.user_oauth_token for update
  using ( auth.uid() = user_id );

--------------------------------
-- End: User OAuth Token Table
--------------------------------

--------------------------------
-- Start: User Email Sync History Table
--------------------------------
create type inbox_type as enum ('candidate', 'recruiter');

-- User Email Sync History Table
create table public.user_email_sync_history (
    user_id uuid references auth.users(id) on delete cascade not null,
    inbox_type inbox_type not null,
    email text not null,
    history_id int8 not null,
    -- track successful sync attempts
    synced_at timestamp with time zone not null default now(),
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),

    primary key (user_id, inbox_type, email)
);

create trigger handle_updated_at_user_email_sync_history before update on public.user_email_sync_history
  for each row execute procedure moddatetime (updated_at);

-- lock down permissions on public.user_email_sync_history
-- only service role can insert. user cannot modify
alter table public.user_email_sync_history enable row level security;

create policy "Users can view their own email sync history"
  on public.user_email_sync_history for select
  using ( auth.uid() = user_id );

-- enable realtime
alter publication supabase_realtime add table user_email_sync_history;

--------------------------------
-- End: User Email Sync History Table
--------------------------------

--------------------------------
-- Start: Waitlist Table
--------------------------------

-- Waitlist table
create table public.waitlist (
    user_id uuid references auth.users(id) on delete cascade not null,
    -- duplicate email for convenience
    email text not null,
    first_name text not null,
    last_name text not null,
    linkedin_url text not null,
    responses jsonb not null default '{}'::jsonb,
    can_create_account boolean not null default false,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),

    primary key (user_id)
);

create trigger handle_updated_at_waitlist before update on public.waitlist
  for each row execute procedure moddatetime (updated_at);

-- enable RLS so users can't modify they waitlist entry
alter table public.waitlist enable row level security;

create policy "Users can view their own waitlist entry"
  on public.waitlist for select
  using ( auth.uid() = user_id );

--------------------------------
-- End: Waitlist Table
--------------------------------

--------------------------------
-- Start: User (Candidate) Profile Table
--------------------------------

-- user_profile table
create table public.user_profile (
    user_id uuid references auth.users(id) on delete cascade not null,
    -- duplicate email for convenience
    email text not null,
    first_name text not null,
    last_name text not null,
    -- user email settings
    -- in future, this will be migrated to a separate table to support multi-email accounts
    is_active boolean not null default true,
    auto_archive boolean not null default false,
    auto_contribute boolean not null default false,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),

    primary key (user_id)
);

create trigger handle_updated_at_profile before update on public.user_profile
  for each row execute procedure moddatetime (updated_at);

alter table public.user_profile enable row level security;

create policy "Users can view their own profile"
  on public.user_profile for select
  using ( auth.uid() = user_id );

create policy "Users can update their own profile"
  on public.user_profile for update
  using ( auth.uid() = user_id );

--------------------------------
-- End: User (Candidate) Profile Table
--------------------------------

--------------------------------
-- Start: User (Candidate) Email Stat Table
--------------------------------

-- user_email_stat table
-- simple table to keep track of realtime user facing statistics
create table public.user_email_stat (
    user_id uuid references auth.users(id) on delete cascade not null,
    -- in future, user can have multiple emails
    email text not null,
    -- free form stat id
    stat_id text not null,
    -- integer for now
    stat_value integer not null default 0,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),

    primary key (user_id, email, stat_id)
);

create trigger handle_updated_at_user_email_stat before update on public.user_email_stat
  for each row execute procedure moddatetime (updated_at);

-- enable RLS
alter table public.user_email_stat enable row level security;
create policy "Users can view their own email stats"
  on public.user_email_stat for select
  using ( auth.uid() = user_id );

-- enable realtime
alter publication supabase_realtime add table user_email_stat;

create or replace function increment_user_email_stat (user_id uuid, email text, stat_id text, stat_value int)
returns void as
$$
  insert into public.user_email_stat(user_id, email, stat_id, stat_value)
  values (user_id, email, stat_id, stat_value)
  on conflict (user_id, email, stat_id)
  do update set
      stat_value = user_email_stat.stat_value + excluded.stat_value;
$$
language sql volatile;

--------------------------------
-- End: User (Candidate) Email Stat Table
--------------------------------

--------------------------------
-- Start: User (Candidate) Email Job Table
--------------------------------


-- user_email_job
create table public.user_email_job (
  job_id uuid not null default uuid_generate_v4(),
  user_id uuid references auth.users(id) on delete cascade not null,
  user_email text not null,
  -- thread id of the recruiting email
  email_thread_id text not null,
  -- time when the email was delivered to user
  emailed_at timestamp with time zone not null,
  company text not null,
  job_title text not null,
  -- use catch-all jsonb for everything else while we figure out the job schema
  data jsonb not null default '{}'::jsonb,
  created_at timestamp with time zone not null default now(),
  updated_at timestamp with time zone not null default now(),

  -- prevent duplicate email_thread_id processing
  unique (user_email, email_thread_id),
  -- prevent duplicate job posting
  -- for now consider, (user_id, company, job_title) to be a unique job posting
  unique (user_id, company, job_title),
  primary key (job_id)
);

create trigger handle_updated_at_user_email_job before update on public.user_email_job
  for each row execute procedure moddatetime (updated_at);

-- enable realtime
alter publication supabase_realtime add table user_email_job;

-- enable RLS
alter table public.user_email_job enable row level security;
create policy "Users can view their own jobs"
  on public.user_email_job for select
  using ( auth.uid() = user_id );

-- enable RLS to allow users to delete their own jobs
create policy "Users can delete their own jobs"
on "public"."user_email_job"
as permissive
for delete
to public
using ((auth.uid() = user_id));


--------------------------------
-- End: User (Candidate) Email Job Table
--------------------------------

--------------------------------
-- Start: Company Table
--------------------------------

-- company table
create table public.company (
    company_id uuid not null default uuid_generate_v4(),
    company_name text not null,
    website text not null,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),

    primary key (company_id)
);

create trigger handle_updated_at_company before update on public.company
  for each row execute procedure moddatetime (updated_at);

-- enable realtime
alter publication supabase_realtime add table public.company;
alter table public.company enable row level security;

--------------------------------
-- End: Company Table
--------------------------------

--------------------------------
-- Start: Recruiter Table
--------------------------------

-- recruiter table
-- similar to user_profile but for recruiters
create table public.recruiter (
    user_id uuid references auth.users(id) on delete cascade not null,
    -- duplicate email for convenience
    email text not null,
    first_name text not null,
    last_name text not null,
    -- email settings
    -- for now, use jsonb instead of a separate table
    -- json is keyed by email address
    email_settings jsonb not null default '{}'::jsonb,
    -- collect additional information when recruiters create an account
    responses jsonb not null default '{}'::jsonb,

    company_id uuid references public.company(company_id) not null,

    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),

    primary key (user_id)
);

create trigger handle_updated_at_recruiter before update on public.recruiter
  for each row execute procedure moddatetime (updated_at);

-- enable realtime
alter publication supabase_realtime add table public.recruiter;
alter table public.recruiter enable row level security;

create policy "Recruiters can view their own profile"
  on public.recruiter for select
  using ( auth.uid() = user_id );

create policy "Recruiters can update their own profile"
  on public.recruiter for update
  using ( auth.uid() = user_id );

create policy "Recruiters can view their own company"
  on public.company for select
  using ( auth.uid() in (
      select user_id
      from public.company
      inner join recruiter using (company_id)
      where user_id = auth.uid()
  ));

create policy "Recruiters can view their own company"
  on public.company for select
  using ( company_id in (
      select company_id
      from public.recruiter
      where user_id = auth.uid()
  ));

--------------------------------
-- End: Recruiter Table
--------------------------------

--------------------------------
-- Start: Jobs Table
--------------------------------

-- jobs
create table public.job (
    job_id uuid not null default uuid_generate_v4(),
    title text not null,
    description_url text not null,

    recruiter_id uuid references public.recruiter(user_id) on delete cascade not null,
    company_id uuid references public.company(company_id) on delete cascade not null,

    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),

    primary key (job_id)
);

create trigger handle_updated_at_job before update on public.job
  for each row execute procedure moddatetime (updated_at);

-- enable real-time
alter publication supabase_realtime add table public.job;
alter table public.job enable row level security;

create policy "Recruiters can view their jobs"
  on public.job for select
  using ( auth.uid() = recruiter_id );

create policy "Recruiters can insert their jobs"
  on public.job for insert
  with check ( auth.uid() = recruiter_id );

create policy "Recruiters can update their jobs"
  on public.job for update
  using ( auth.uid() = recruiter_id );

create policy "Recruiters can delete their jobs"
  on public.job for delete
  using ( auth.uid() = recruiter_id );

create policy "Companies can view their jobs"
  on public.job for select
  using ( company_id in (
      select company_id
      from public.recruiter
      where user_id = auth.uid()
  ));

--------------------------------
-- End: Jobs Table
--------------------------------

--------------------------------
-- Start: Views on User OAuth Tokens
--------------------------------

-- views on oauth tokens
create or replace view candidate_oauth_token with (security_invoker) as
select
  user_oauth_token.*
from user_oauth_token
inner join user_profile using (user_id);

create or replace view recruiter_oauth_token with (security_invoker) as
select
  user_oauth_token.*
from user_oauth_token
inner join recruiter using (user_id);

--------------------------------
-- End: Views on User OAuth Tokens
--------------------------------

--------------------------------
-- Start: Get Candidate/Recruiter Given Connected Email
--------------------------------

create or replace function get_user_profile_by_email (input text)
returns user_profile as
$$
select
  user_profile.*
from public.user_profile
left join public.user_oauth_token using (user_id)
where user_profile.email = input OR user_oauth_token.email = input
limit 1;
$$
language sql stable;

create or replace function get_recruiter_by_email (input text)
returns recruiter as
$$
select
  recruiter.*
from public.recruiter
left join public.user_oauth_token using (user_id)
where recruiter.email = input OR user_oauth_token.email = input
limit 1;
$$
language sql stable;

--------------------------------
-- End: Get Candidate/Recruiter Given Connected Email
--------------------------------

--------------------------------
-- Start: Recruiter Outbound Tables
--------------------------------


-- Recruiter outbound tables
-- These tables help us connect recruiter's existing outbound to SRC
create table public.recruiter_outbound_template (
  template_id uuid not null default uuid_generate_v4(),
  recruiter_id uuid references public.recruiter(user_id) on delete cascade not null,
  -- nullable job_id
  job_id uuid references public.job(job_id) on delete set null,

  subject text not null,
  body text not null,
  normalized_content text not null,
  metadata jsonb not null default '{}'::jsonb,

  created_at timestamp with time zone not null default now(),
  updated_at timestamp with time zone not null default now(),

  primary key (template_id)
);

create trigger handle_updated_at_recruiter_outbound_template before update on public.recruiter_outbound_template
  for each row execute procedure moddatetime (updated_at);

-- enable real-time
alter publication supabase_realtime add table public.recruiter_outbound_template;

-- enable RLS
alter table public.recruiter_outbound_template enable row level security;

create policy "Recruiters can view their outbound templates"
  on public.recruiter_outbound_template for select
  using ( auth.uid() = recruiter_id );

create policy "Recruiters can insert their outbound templates"
  on public.recruiter_outbound_template for insert
  with check ( auth.uid() = recruiter_id );

create policy "Recruiters can update their outbound templates"
  on public.recruiter_outbound_template for update
  using ( auth.uid() = recruiter_id );

create policy "Recruiters can delete their outbound templates"
  on public.recruiter_outbound_template for delete
  using ( auth.uid() = recruiter_id );

create or replace function list_similar_recruiter_outbound_templates (user_id uuid, input text)
returns table (
  template_id uuid,
  recruiter_id uuid,
  job_id uuid,
  subject text,
  body text,
  metadata jsonb,
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  similarity real
) as
$$
select
    template_id,
    recruiter_id,
    job_id,
    subject,
    body,
    metadata,
    created_at,
    updated_at,
    similarity(normalized_content, input) as "similarity"
from public.recruiter_outbound_template
where recruiter_id = user_id
and normalized_content % input
order by 9 desc;
$$
language sql stable;

-- add index for similarity search
create index idx_recruiter_outbound_template_normalized_content on public.recruiter_outbound_template using gin (normalized_content gin_trgm_ops);

-- default is 0.3
-- TODO: Adjust once we have more data
set pg_trgm.similarity_threshold = 0.5;

--------------------------------
--------------------------------

create table public.recruiter_outbound_message (
  recruiter_id uuid references public.recruiter(user_id) on delete cascade not null,
  -- message ID from the email provider
  message_id text not null,
  -- RFC2822 Message ID
  internal_message_id text not null,
  from_email text not null,
  to_email text not null,
  sent_at timestamp with time zone not null,

  template_id uuid references public.recruiter_outbound_template(template_id) on delete set null,

  created_at timestamp with time zone not null default now(),
  updated_at timestamp with time zone not null default now(),

  primary key (message_id)
);

create trigger handle_updated_at_recruiter_outbound_message before update on public.recruiter_outbound_message
  for each row execute procedure moddatetime (updated_at);

-- enable real-time
alter publication supabase_realtime add table public.recruiter_outbound_message;

-- enable RLS
alter table public.recruiter_outbound_message enable row level security;

create policy "Recruiters can view their outbound messages"
  on public.recruiter_outbound_message for select
  using ( auth.uid() = recruiter_id );

--------------------------------
-- End: Recruiter Outbound Tables
--------------------------------

--------------------------------
-- Start: Candidate Company Inbound Tables & Triggers
--------------------------------

-- Candidate Company Inbound Tables
-- These tables help us generate a job board for candidates from recruiter outbound
create table public.candidate_company_inbound (
  candidate_email text not null,
  -- nullable
  candidate_id uuid references public.user_profile(user_id) on delete set null,
  company_id uuid references public.company(company_id) on delete cascade not null,
  recruiter_id uuid references public.recruiter(user_id) on delete set null,
  template_id uuid references public.recruiter_outbound_template(template_id) on delete cascade not null,
  -- nullable
  job_id uuid references public.job(job_id) on delete set null,

  created_at timestamp with time zone not null default now(),
  updated_at timestamp with time zone not null default now(),

  primary key (candidate_email, company_id, template_id)
);

create trigger handle_updated_at_candidate_company_inbound before update on public.candidate_company_inbound
  for each row execute procedure moddatetime (updated_at);

-- enable real-time
alter publication supabase_realtime add table public.candidate_company_inbound;

-- enable RLS
alter table public.candidate_company_inbound enable row level security;

-- TODO: Figure out company RLS policy
create policy "Recruiters can view their inbound candidates"
  on public.candidate_company_inbound for select
  using ( auth.uid() = recruiter_id );

create policy "Candidates can view their inbound"
  on public.candidate_company_inbound for select
  using ( auth.uid() = candidate_id );

create policy "Candidates can view inbound to their emails"
  on public.candidate_company_inbound for select
  using ( candidate_email in (
      select email
      from public.user_oauth_token
      where user_id = auth.uid()
  ));

-- create function to update job_id for a given template_id
create or replace function update_job_for_template_candidate_company_inbound_trigger()
returns trigger as
$$
begin
update public.candidate_company_inbound
set job_id = new.job_id
where template_id = new.template_id;
return null;
end;
$$
-- escalate to security definer to grant update permissions on candidate_company_inbound
security definer
language plpgsql volatile;

-- create function to update candidate_id for a given candidate_email
create or replace function update_candidate_for_email_candidate_company_inbound_trigger()
returns trigger as
$$
begin
update public.candidate_company_inbound
set candidate_id = new.user_id
where candidate_email = new.email;
return null;
end;
$$
-- escalate to security definer to grant update permissions on candidate_company_inbound
security definer
language plpgsql volatile;

-- create function to insert a row on every new recruiter_outbound_message
create or replace function insert_candidate_company_inbound_trigger()
returns trigger as
$$
begin
with input as (
  select
    *
  from (values (new.to_email, new.recruiter_id, new.template_id)) as t (candidate_email, recruiter_id, template_id)
),
c as (
  select
    user_id as candidate_id,
    coalesce(email, new.to_email) as candidate_email
  from get_user_profile_by_email(new.to_email)
),
t as (
  select
    template_id,
    job_id
  from public.recruiter_outbound_template
  where template_id = new.template_id
)
insert into public.candidate_company_inbound (
  candidate_email,
  candidate_id,
  template_id,
  job_id,
  recruiter_id,
  company_id
)
select
  input.candidate_email,
  c.candidate_id,
  input.template_id,
  t.job_id,
  input.recruiter_id,
  r.company_id
from input
inner join public.recruiter r on r.user_id = input.recruiter_id
inner join t using (template_id)
left join c using (candidate_email)
on conflict do nothing;
return null;
end;
$$
language plpgsql volatile;

-- create a trigger for inserts into to recruiter_outbound_message
create trigger insert_candidate_company_inbound_trigger_after_insert after insert on public.recruiter_outbound_message
  for each row execute function insert_candidate_company_inbound_trigger();

-- create a trigger for inserts or updates into to recruiter_outbound_template
-- TODO: Separate triggers for update and insert
create trigger candidate_company_inbound_trigger_recruiter_outbound_template after update on public.recruiter_outbound_template
  for each row
  when (old.job_id is distinct from new.job_id)
  execute function update_job_for_template_candidate_company_inbound_trigger();

-- create a trigger for inserts into to user_oauth_token
create trigger candidate_company_inbound_trigger_user_oauth_token after insert or update on public.user_oauth_token
  for each row execute function update_candidate_for_email_candidate_company_inbound_trigger();

-- create a trigger for inserts into to user_profile
create trigger candidate_company_inbound_trigger_user_profile after insert on public.user_profile
  for each row execute function update_candidate_for_email_candidate_company_inbound_trigger();

--------------------------------
-- End: Candidate Company Inbound Tables & Triggers
--------------------------------


--------------------------------
-- Start: Views for Counts of Candidates (Recipients) per Template and for Jobs
--------------------------------
create or replace view outbound_template_recipient_count with (security_invoker) as
select
   template_id,
   count(distinct to_email) as num_recipients
from recruiter_outbound_message
group by template_id;

create or replace view job_candidate_count with (security_invoker) as
select
  job_id,
  count(distinct candidate_email) as num_candidates
from candidate_company_inbound
group by job_id;

create or replace view candidate_job_count with (security_invoker) as
select
  candidate_id,
  count(distinct job_id) as num_jobs
from candidate_company_inbound
where candidate_id is not null
group by candidate_id;

create or replace view candidate_job_count_unverified with (security_invoker) as
select
  user_id,
  count(distinct job_id) as num_jobs
from user_email_job
group by user_id;
--------------------------------
-- End: Views for Counts of Candidates (Recipients) per Template and for Jobs
--------------------------------

--------------------------------
-- Start: Candidate Job Interest
--------------------------------
create type job_interest as enum ('interested', 'not_interested', 'saved');

create table candidate_job_interest (
  candidate_id uuid references user_profile(user_id) on delete cascade,
  -- Note: no foreign key constraint because we want to share this table with user_email_job
  -- we need to create a single source of truth for jobs...
  job_id uuid not null,
  interest job_interest not null,
  created_at timestamp with time zone not null default now(),
  updated_at timestamp with time zone not null default now(),
  primary key (candidate_id, job_id)
);

create trigger handle_updated_at_candidate_job_interest before update on public.candidate_job_interest
  for each row execute procedure moddatetime (updated_at);

-- enable realtime
alter publication supabase_realtime add table candidate_job_interest;

alter table public.candidate_job_interest enable row level security;

create policy "Candidates can view their own job interest"
  on public.candidate_job_interest for select
  using ( auth.uid() = candidate_id );

create policy "Candidates can update their own job interest"
  on public.candidate_job_interest for update
  using ( auth.uid() = candidate_id );

create policy "Candidates can insert their own job interest"
  on public.candidate_job_interest for insert
  with check ( auth.uid() = candidate_id );

-- TODO: Allow recruiters to view job interest

--------------------------------
-- Start: Candidate Job Interest
--------------------------------

--------------------------------
-- Start: Job Board View
--------------------------------
-- create a view that combines verified and unofficial jobs into one table for candidates
create or replace view vw_job_board with (security_invoker ) as (
  with verified_jobs as (
    select
      cc.candidate_id as user_id,
      cc.candidate_email as user_email,
      job.job_id,
      job.title as job_title,
      job.description_url as job_description_url,
      candidate_job_interest.interest as job_interest,
      company.company_name,
      company.website as company_website,
      recruiter.first_name || ' ' || recruiter.last_name as recruiter_name,
      recruiter.email as recruiter_email,
      cc.created_at as emailed_at,
      true as is_verified
    from public.candidate_company_inbound cc
    inner join public.company on cc.company_id = company.company_id
    inner join public.recruiter on cc.recruiter_id = recruiter.user_id
    inner join public.job on cc.job_id = job.job_id
    left join public.candidate_job_interest on cc.candidate_id = candidate_job_interest.candidate_id and cc.job_id = candidate_job_interest.job_id
  ),
  unofficial_jobs as (
    select
      user_id,
      user_email,
      user_email_job.job_id,
      job_title,
      '' as job_description_url,
      candidate_job_interest.interest as job_interest,
      company as company_name,
      '' as company_website,
      data ->> 'recruiter' as recruiter_name,
      data ->> 'recruiter_email' as recruiter_email,
      emailed_at,
      false as is_verified
    from public.user_email_job
    left join public.candidate_job_interest
    on user_email_job.user_id = candidate_job_interest.candidate_id
    and user_email_job.job_id = candidate_job_interest.job_id
  )
  -- line up columns
  select * from verified_jobs
  union all
  select * from unofficial_jobs
);
--------------------------------
-- End: Job Board View
--------------------------------
